package dao

import (
	"database/sql"
	"fmt"
	"time"

	"mby.fr/cmdtest/model"
	"mby.fr/utils/zql"
)

func NewTest(db *zql.SynchronizedDB) (d Test, err error) {
	d.db = db
	d.init()
	return
}

type Test struct {
	db *zql.SynchronizedDB
}

func (d Test) init() (err error) {
	_, err = d.db.Exec(`
		CREATE TABLE IF NOT EXISTS tested (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			title TEXT NOT NULL,
			seq INTEGER NOT NULL,
			outcome TEXT NOT NULL,
			errorMsg TEXT NULL,
			duration INTEGER NOT NULL,
			passed INTEGER NOT NULL,
			failed INTEGER NOT NULL,
			errored INTEGER NOT NULL,
			ignored INTEGER NOT NULL,
			exitCode INTEGER NOT NULL,
			stdout TEXT NULL,
			stderr TEXT NULL,
			report TEXT NULL,
			FOREIGN KEY(suite) REFERENCES suite(name)
		);

		CREATE TABLE IF NOT EXISTS assertion_result (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			testId INTEGER NOT NULL,
			prefix TEXT NOT NULL,
			name TEXT NOT NULL,
			op TEXT NOT NULL,
			expected TEXT NOT NULL,
			value TEXT NOT NULL,
			errorMsg TEXT NULL,
			success INTEGER NOT NULL,

			FOREIGN KEY(testId) REFERENCES test(name)
		);		
	`)
	return
}

func (d Test) GetSuiteOutcome(suite string) (outcome model.SuiteOutcome, err error) {
	var passedCount, failedCount, erroredCount, ignoredCount uint32
	var title, oc, prefix, name, op, expected, value, errorMsg string
	var success bool
	var startTime, endTime int64
	var failedAssertionsMessages []string

	tx, err := d.db.Begin()
	if err != nil {
		return
	}
	defer tx.Rollback()

	row := tx.QueryRow(`
		SELECT sum(passed), sum(failed), sum(errored), sum(ignored)
		FROM tested
		WHERE suite = @suite
	`, sql.Named("suite", suite))
	err = row.Scan(&passedCount, &failedCount, &erroredCount, &ignoredCount)
	if err != nil {
		return
	}

	rows, err := tx.Query(`
		SELECT t.title, t.outcome, a.prefix, a.name, a.op, a.expected, a.value, a.errorMsg, a.success
		FROM assertion_result a INNER JOIN tested t ON a.testId = t.id
		WHERE t.name = @suite 
	`, sql.Named("suite", suite))
	if err != nil {
		return
	}
	defer rows.Close()

	row = tx.QueryRow(`
		SELECT startTime, endTime
		FROM suite
		WHERE name = @suite 
	`, sql.Named("suite", suite))
	err = row.Scan(&startTime, &endTime)
	if err != nil {
		return
	}

	err = tx.Commit()
	if err != nil {
		return
	}

	for rows.Next() {
		if err = rows.Scan(&title, &oc, &prefix, &name, &op, &expected, &value, &errorMsg, &success); err != nil {
			return
		}
		ocm := model.Outcome(oc)
		var message string
		switch ocm {
		case model.PASSED:
			// Nothing to do
		case model.FAILED:
			message = title + "  => " + prefix + name + op + expected
		case model.IGNORED:
			// Nothing to do
		case model.ERRORED:
			message = title + "  => not executed"
		case model.TIMEOUT:
			message = title + "  => timed out"
		default:
			err = fmt.Errorf("outcome %s not supported", ocm)
		}

		failedAssertionsMessages = append(failedAssertionsMessages, message)
	}

	testCount := passedCount + failedCount + erroredCount + ignoredCount
	duration := time.Duration((endTime - startTime) * 1000)
	var ocm model.Outcome
	if testCount == passedCount {
		ocm = model.PASSED
	} else if testCount == ignoredCount {
		ocm = model.IGNORED
	} else if erroredCount > 0 {
		ocm = model.ERRORED
	} else {
		ocm = model.FAILED
	}

	outcome.TestSuite = suite
	outcome.Duration = duration
	outcome.TestCount = testCount
	outcome.PassedCount = passedCount
	outcome.FailedCount = failedCount
	outcome.ErroredCount = erroredCount
	outcome.IgnoredCount = ignoredCount
	outcome.Outcome = ocm
	outcome.FailureReports = failedAssertionsMessages

	return
}

func (d Test) SaveTestOutcome(suite, title string, outcome model.TestOutcome) (err error) {
	seq := outcome.Seq
	errorMsg := outcome.Err
	micros := outcome.Duration.Microseconds()
	passed := outcome.Outcome == model.PASSED
	ignored := outcome.Outcome == model.IGNORED
	failed := outcome.Outcome == model.FAILED
	errored := outcome.Outcome == model.ERRORED
	exitCode := outcome.ExitCode
	stdout := outcome.Stdout
	stderr := outcome.Stderr
	report := "NO REPORT"
	_, err = d.db.Exec(`
		INSERT INTO tested(suite, seq, title, errorMsg, duration, passed, ignored, failed, errored, exitCode, stdout, stderr, report) 
		VALUES (@suite, @seq, @errorMsg, @micros, @passed, @ignored, @failed, @errored, @exitCode, @stdout, @stderr, @report);
	`, sql.Named("suite", suite), sql.Named("seq", seq), sql.Named("title", title),
		sql.Named("errorMsg", errorMsg), sql.Named("micros", micros),
		sql.Named("passed", passed), sql.Named("ignored", ignored),
		sql.Named("failed", failed), sql.Named("errored", errored),
		sql.Named("exitCode", exitCode), sql.Named("stdout", stdout),
		sql.Named("stderr", stderr), sql.Named("report", report))
	return
}

func (d Test) ClearSuite(suite string) (err error) {
	_, err = d.db.Exec(`
		DELETE FROM assertion_result
		WHERE testId IN (select id from tested where suite = @suite);
		DELETE FROM tested
		WHERE suite = @suite;
	`, sql.Named("suite", suite))
	return
}

/*
func (d Test) SaveTested(suite string, seq int, cfg model.Config) (err error) {
	// TODO
	return
}
*/
