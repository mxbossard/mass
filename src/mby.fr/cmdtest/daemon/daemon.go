package daemon

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/gofrs/flock"

	"mby.fr/cmdtest/model"
	"mby.fr/cmdtest/repo"
	"mby.fr/cmdtest/service"
	"mby.fr/utils/filez"
)

const (
	DaemonLockFilename = "daemon.lock"
	DaemonPidFilename  = "daemon.pid"
	LockWatingSecs     = 5
	ExtraRunningSecs   = 10
)

var logger = slog.New(slog.NewTextHandler(os.Stderr, model.DefaultLoggerOpts))

/*
Keys:
- Only one daemon running at once
- Do not miss queued tests

Ideas:
- Init process
  - acquire lockfile (daemon is starting) or wait
  - if PID file is present => return
  - Write PID file
  - Start daemon process
  - Release lock file
  - Start daemon

- Running daemon
  - Loop Unqueuing tests
  - If no test for x seconds => acquire lock file (daemon is stopping)
  - Last unqueue (May take a lot of time testing)
  - Stop the daemon

- Stop daemon
  - Rm PID file
  - Release lock file
  - Exit 0

*/

type daemon struct {
	repo repo.FileRepo
}

func (d *daemon) run() {
	var lastUnqueue time.Time
	for {
		testOp, err := d.repo.UnqueueOperation()
		if err != nil {
			panic(err)
		}
		if testOp == nil {
			// nothing to unqueue wait 1ms
			if time.Since(lastUnqueue) > ExtraRunningSecs*time.Second {
				// More than ExtraRunningSecs since last unqueue
				break
			}
			time.Sleep(time.Millisecond)
			continue
		}
		lastUnqueue = time.Now()
		d.performTest(testOp.Def)
		d.repo.State.ReportOperationDone(testOp)
	}
}

func (d daemon) performTest(testDef model.TestDefinition) {
	//logger.Warn("daemon performing test", "testDef", testDef)
	_ = service.ProcessTestDef(testDef)
}

func (d daemon) ReadPid() string {
	pidFilepath := filepath.Join(d.repo.BackingFilepath(), DaemonPidFilename)
	pid, err := filez.ReadString(pidFilepath)
	if os.IsNotExist(err) {
		return ""
	} else if err != nil {
		panic(err)
	}
	return pid
}

func (d daemon) WritePid() {
	pidFilepath := filepath.Join(d.repo.BackingFilepath(), DaemonPidFilename)
	err := filez.WriteString(pidFilepath, fmt.Sprintf("%d", os.Getpid()), 0600)
	if err != nil {
		panic(err)
	}
}

func (d daemon) ClearPid() {
	pidFilepath := filepath.Join(d.repo.BackingFilepath(), DaemonPidFilename)
	err := os.Remove(pidFilepath)
	if err != nil && !os.IsNotExist(err) {
		panic(err)
	}
}

func TakeOver() {
	if os.Args[1] == "@_daemon" {
		//fmt.Printf("@_daemon args: %s", os.Args)
		if len(os.Args) != 3 {
			panic("bad usage of @_daemon")
		}
	} else {
		return
	}

	token := os.Args[2]
	repo := repo.New(token)
	d := daemon{repo: repo}
	lockFilepath := filepath.Join(repo.BackingFilepath(), DaemonLockFilename)
	fileLock := flock.New(lockFilepath)

	// Wait to acquire file lock
	lockCtx, cancel := context.WithTimeout(context.Background(), LockWatingSecs*time.Second)
	defer cancel()
	locked, err := fileLock.TryLockContext(lockCtx, time.Millisecond)
	if err != nil {
		panic(err)
	}
	if !locked {
		os.Exit(2)
	}

	// If PID file already exists exit => already running
	if d.ReadPid() != "" {
		fileLock.Unlock()
		os.Exit(3)
	}

	// Write PID file
	d.WritePid()

	// Release file lock
	err = fileLock.Unlock()
	if err != nil {
		panic(err)
	}

	// Run daemon
	d.run()

	// Lock prior last unqueue
	locked, err = fileLock.TryLockContext(lockCtx, time.Millisecond)
	if err != nil {
		panic(err)
	}
	if !locked {
		os.Exit(2)
	}

	// Last unqueue
	testOp, err := d.repo.UnqueueOperation()
	if err != nil {
		panic(err)
	}
	if testOp != nil {
		d.performTest(testOp.Def)
	}

	// Clear PID file
	d.ClearPid()

	// Release file lock
	fileLock.Unlock()

	os.Exit(0)
}

func LanchProcessIfNeeded(token string) error {
	if token == "" {
		// No token => no daemon to launch
		return nil
	}
	// FIXME: add retries ?
	cwd, err := os.Getwd()
	if err != nil {
		return err
	}
	cmd := exec.Command(os.Args[0], "@_daemon", token)
	cmd.Dir = cwd
	cmd.Env = os.Environ()
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Start()
	if err != nil {
		return err
	}
	err = cmd.Process.Release()
	if err != nil {
		return err
	}
	logger.Debug("daemon process released")
	return nil
}
