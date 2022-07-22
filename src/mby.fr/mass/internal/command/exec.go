package command

import (
	"os/exec"

	"mby.fr/mass/internal/logger"
	"mby.fr/utils/inout"
)

func RunLogging(cmd *exec.Cmd, logger logger.ActionLogger) (err error) {
	errors := make(chan error, 10)

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		errors <- err
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		errors <- err
	}

	go inout.CopyChannelingErrors(stdout, logger.Out(), errors)
	go inout.CopyChannelingErrors(stderr, logger.Err(), errors)

	err = cmd.Start()
	if err != nil {
		//logger.Flush()
		errors <- err
	}
	err = cmd.Wait()
	if err == nil {
		// Use select to not block if no error in channel
		select {
		case err = <-errors:
		default:
		}
	}
	return
}
