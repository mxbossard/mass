package repo

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"mby.fr/cmdtest/dao"
	"mby.fr/cmdtest/model"
	"mby.fr/utils/errorz"
	"mby.fr/utils/zql"
)

type dbRepo struct {
	dirpath   string
	token     string
	isolation string
	db        *zql.SynchronizedDB
	suiteDao  dao.Suite
	queueDao  dao.Queue
	testDao   dao.Test
	//lastUpdate time.Time
}

func (r *dbRepo) Init() (err error) {
	if r.db != nil {
		err = r.db.Close()
		if err != nil {
			return
		}
	}
	db, err := dao.DbOpen(r.dirpath)
	if err != nil {
		return
	}
	if r.db == nil {
		r.db = db
	} else {
		*r.db = *db
	}
	return
}

func (r dbRepo) BackingFilepath() string {
	path, err := forgeWorkDirectoryPath(r.token, r.isolation)
	if err != nil {
		errorz.Fatal(err)
	}
	return path
}

func (r dbRepo) MockDirectoryPath(testSuite string, testId uint16) (mockDir string, err error) {
	var path string
	path, err = testSuiteDirectoryPath(testSuite, r.token, r.isolation)
	if err != nil {
		return
	}
	mockDir = filepath.Join(path, fmt.Sprintf("__mock_%d", testId))
	// create a mock dir
	err = os.MkdirAll(mockDir, 0755)
	if err != nil {
		return
	}
	return
}

func (r dbRepo) SaveGlobalConfig(cfg model.Config) (err error) {
	err = r.suiteDao.SaveGlobalConfig(cfg)
	return
}

func (r dbRepo) GetGlobalConfig() (cfg model.Config, err error) {
	found, err := r.suiteDao.FindGlobalConfig()
	if err != nil {
		return
	}
	if found != nil {
		cfg = *found
	} else {
		// global config does not exists yet
		// create a new default one
		cfg = model.NewGlobalDefaultConfig()
		cfg.Token.Set(r.token)
		cfg.GlobalStartTime.Set(time.Now())
		err = r.SaveGlobalConfig(cfg)
	}
	return
}

func (r dbRepo) InitSuite(cfg model.Config) (err error) {
	err = r.ClearTestSuite(cfg.TestSuite.Get())
	if err != nil {
		return
	}
	//err = persistSuiteConfig(r.token, cfg)
	err = r.SaveSuiteConfig(cfg)
	if err != nil {
		err = fmt.Errorf("unable to init suite: %w", err)
	}
	return
}

func (r dbRepo) SaveSuiteConfig(cfg model.Config) (err error) {
	err = r.suiteDao.SaveSuiteConfig(cfg.TestSuite.Get(), cfg)
	return
}

func (r dbRepo) GetSuiteConfig(testSuite string, initless bool) (cfg model.Config, err error) {
	// TODO
	defer func() {
		if err != nil {
			err = fmt.Errorf("cannot load suite config: %w", err)
		}
	}()

	found, err := r.suiteDao.FindSuiteConfig(testSuite)
	if err != nil {
		return
	}
	if found != nil {
		cfg = *found
	} else {
		// suite config does not exists yet
		// create a new default one
		cfg, err = r.GetGlobalConfig()
		if err != nil {
			return
		}
		var suiteCfg model.Config
		if initless {
			//logger.Warn("Saving new initless config", "testSuite", testSuite)
			suiteCfg = model.NewInitlessSuiteDefaultConfig()
		} else {
			//logger.Warn("Saving new inited config", "testSuite", testSuite)
			suiteCfg = model.NewSuiteDefaultConfig()
		}
		suiteCfg.TestSuite.Set(testSuite)
		suiteCfg.SuiteStartTime.Set(time.Now())
		cfg.Merge(suiteCfg)
		err = r.SaveSuiteConfig(cfg)
	}
	return
}

func (r dbRepo) ClearTestSuite(testSuite string) (err error) {
	err = r.testDao.DeleteTest(testSuite)
	if err != nil {
		return
	}
	err = r.suiteDao.DeleteSuite(testSuite)
	return
}

func (r dbRepo) ListTestSuites() (suites []string, err error) {
	suites, err = r.suiteDao.ListPassedFailedErrored()
	return
}

func (r dbRepo) SaveTestOutcome(outcome model.TestOutcome) (err error) {
	err = r.testDao.SaveTestOutcome(outcome)
	if err != nil {
		return
	}
	err = r.suiteDao.UpdateSuiteOutcome(outcome.TestSuite, outcome.Outcome)
	return
}

func (r dbRepo) UpdateLastTestTime(testSuite string) {
	err := r.suiteDao.UpdateSuiteEndTime(testSuite, time.Now())
	if err != nil {
		errorz.Fatal(err)
	}
}

func (r dbRepo) LoadSuiteOutcome(testSuite string) (outcome model.SuiteOutcome, err error) {
	outcome, err = r.testDao.GetSuiteOutcome(testSuite)
	return
}

func (r dbRepo) IncrementSuiteSeq(testSuite, name string) (n uint16) {
	// FIXME should this be used ?

	var err error
	if name == model.TestSequenceFilename {
		n, err = r.suiteDao.NextSeq(testSuite)
	} else if name == model.TooMuchSequenceFilename {
		n, err = r.suiteDao.IncrementTooMuchCount(testSuite)
	} else {
		// Other seq increment are supported by counter in DB
		return 9999
	}
	if err != nil {
		errorz.Fatal(err)
	}
	return
}

func (r dbRepo) TestCount(testSuite string) (n uint16) {
	n, err := r.suiteDao.TestCount(testSuite)
	if err != nil {
		errorz.Fatal(err)
	}
	return
}

func (r dbRepo) PassedCount(testSuite string) (n uint16) {
	n, err := r.testDao.PassedCount(testSuite)
	if err != nil {
		errorz.Fatal(err)
	}
	return
}

func (r dbRepo) IgnoredCount(testSuite string) (n uint16) {
	n, err := r.testDao.IgnoredCount(testSuite)
	if err != nil {
		errorz.Fatal(err)
	}
	return
}

func (r dbRepo) FailedCount(testSuite string) (n uint16) {
	n, err := r.testDao.FailedCount(testSuite)
	if err != nil {
		errorz.Fatal(err)
	}
	return
}

func (r dbRepo) ErroredCount(testSuite string) (n uint16) {
	n, err := r.testDao.ErroredCount(testSuite)
	if err != nil {
		errorz.Fatal(err)
	}
	return
}

func (r dbRepo) TooMuchCount(testSuite string) (n uint16) {
	n, err := r.suiteDao.TooMuchCount(testSuite)
	if err != nil {
		errorz.Fatal(err)
	}
	return
}

func (r dbRepo) QueueOperation(op model.Operater) (err error) {
	err = r.queueDao.QueueOperater(op)
	if err == nil {
		logger.Info("Queued operation", "testSuite", op.Suite(), "kind", op.Kind(), "seq", op.Seq())
	}
	return
}

func (r dbRepo) UnqueueOperation() (op model.Operater, err error) {
	op, err = r.queueDao.UnqueueOperater()
	if op != nil {
		logger.Info("Unqueued operation", "testSuite", op.Suite(), "kind", op.Kind(), "seq", op.Seq(), "err", err)
	}
	return
}

func (r dbRepo) Done(op model.Operater) (err error) {
	if op == nil {
		return
	}

	err = r.queueDao.Done(op)
	//logger.Warn("Unblock() unblocked", "opId", op.Id())
	return
}

func (r dbRepo) WaitOperationDone(op model.Operater, timeout time.Duration) (exitCode int16, err error) {
	exitCode = -1
	start := time.Now()
	for time.Since(start) < timeout {
		var done bool
		done, exitCode, err = r.queueDao.IsOperationsDone(op)
		if done || err != nil {
			// Operater done
			return
		}
		time.Sleep(1 * time.Millisecond)
		logger.Trace("waiting ...", "op", op)
	}
	err = errors.New("WaitOperationDone() timed out")
	return
}

func (r dbRepo) WaitEmptyQueue(testSuite string, timeout time.Duration) (err error) {
	start := time.Now()
	for time.Since(start) < timeout {
		var count int
		count, err = r.queueDao.QueuedOperationsCountBySuite(testSuite, nil)
		if err != nil {
			return
		}
		logger.Trace("WaitEmptyQueue()", "testSuite", testSuite, "count", count)
		if count == 0 {
			// Queue is empty
			return
		}
		time.Sleep(1 * time.Millisecond)
	}
	err = errors.New("WaitEmptyQueue() timed out")
	return
}

func (r dbRepo) WaitAllEmpty(timeout time.Duration) (err error) {
	start := time.Now()
	for time.Since(start) < timeout {
		var count int
		count, err = r.queueDao.QueuedOperationsCount()
		if err != nil {
			return
		}
		if count == 0 {
			// No operation queued
			return
		}
		time.Sleep(1 * time.Millisecond)
	}
	err = errors.New("WaitAllEmpty() timed out")
	return
}

func (r dbRepo) unqueue() (ok bool, op model.Operater, err error) {
	queuedOperationsCount, err := r.queueDao.QueuedOperationsCount()
	if err != nil {
		return
	}
	//logger.Warn("unqueue()", "globalOperationsCount", globalOperationsCount)
	if queuedOperationsCount == 0 {
		return
	}

	op, err = r.queueDao.UnqueueOperater()
	if err != nil || op == nil {
		return
	}

	ok = true
	return
}

func (r dbRepo) Unqueue0() (ok bool, op model.Operater, err error) {
	ok, op, err = r.unqueue()
	//	logger.Warn("Unqueue()", "kind", op.Kind(), "opId", op.Id())
	return
}

func newDbRepo(dirpath, isolation, token string) (r dbRepo, err error) {
	r.dirpath = dirpath
	r.token = token
	r.isolation = isolation

	err = r.Init()
	if err != nil {
		return
	}
	db := r.db
	r.queueDao, err = dao.NewQueue(db)
	if err != nil {
		return
	}
	r.suiteDao, err = dao.NewSuite(db)
	if err != nil {
		return
	}
	r.testDao, err = dao.NewTest(db)
	return
}
