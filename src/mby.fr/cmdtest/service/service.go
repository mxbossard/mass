package service

import (
	"fmt"
	"log"
	"log/slog"
	"os"
	"path/filepath"

	"mby.fr/cmdtest/display"
	"mby.fr/cmdtest/facade"
	"mby.fr/cmdtest/model"
	"mby.fr/cmdtest/repo"
	"mby.fr/cmdtest/utils"
	"mby.fr/utils/cmdz"
	"mby.fr/utils/printz"
	"mby.fr/utils/utilz"
)

var dpl = display.New()

var logger = slog.New(slog.NewTextHandler(os.Stderr, model.DefaultLoggerOpts))

func usage() {
	usagePrinter := printz.NewStandard()
	cmd := filepath.Base(os.Args[0])
	usagePrinter.Errf("cmdtest tool is usefull to test various scripts cli and command behaviors.\n")
	usagePrinter.Errf("You must initialize a test suite (%[1]s @init) before running tests and then report the test (%[1]s @report).\n", cmd)
	usagePrinter.Errf("usage: \t%s @init[=TEST_SUITE_NAME] [@CONFIG_1] ... [@CONFIG_N] \n", cmd)
	usagePrinter.Errf("usage: \t%s <COMMAND> [ARG_1] ... [ARG_N] [@CONFIG_1] ... [@CONFIG_N] [@ASSERTION_1] ... [@ASSERTION_N]\n", cmd)
	usagePrinter.Errf("usage: \t%s @report[=TEST_SUITE_NAME] \n", cmd)
	usagePrinter.Errf("\tCONFIG available: @ignore @stopOnFailure @keepStdout @keepStderr @keepOutputs @timeout=Duration @fork=N\n")
	usagePrinter.Errf("\tCOMMAND and ARGs: the command on which to run tests\n")
	usagePrinter.Errf("\tASSERTIONs available: @fail @success @exit=N @stdout= @stdout~ @stderr= @stderr~ @cmd= @exists=\n")
	usagePrinter.Errf("In complex cases assertions must be correlated by a token. You can generate a token with @init @printToken or @init @exportToken and supply it with @token=\n")
	usagePrinter.Flush()
}

func GlobalConfig(ctx facade.GlobalContext) (exitCode int, err error) {
	// Init or Update Global config
	exitCode = 0
	err = ctx.Save()
	return
}

func InitTestSuite(ctx facade.SuiteContext) (exitCode int, err error) {
	// Clear and Init new test suite
	exitCode = 0
	cfg := ctx.Config

	var token string
	if cfg.PrintToken.Is(true) {
		token, err = utils.ForgeUuid()
		if err != nil {
			return
		}
		fmt.Printf("%s\n", token)
		cfg.Token = utilz.OptionalOf(token)
	} else if cfg.ExportToken.Is(true) {
		token, err = utils.ForgeUuid()
		if err != nil {
			return
		}
		fmt.Printf("export %s=%s\n", model.ContextTokenEnvVarName, token)
		cfg.Token = utilz.OptionalOf(token)
	}

	err = ctx.InitSuite()
	if err != nil {
		ctx.NoErrorOrFatal(err)
	}

	dpl.Suite(ctx)
	return
}

func ReportAllTestSuites(ctx facade.GlobalContext) (exitCode int, err error) {
	exitCode = 1
	token := ctx.Token

	var testSuites []string
	testSuites, err = ctx.Repo.ListTestSuites()
	if err != nil {
		return
	}

	logger.Info("Reporting all suites", "token", token, "suites", testSuites)
	if len(testSuites) == 0 {
		err = fmt.Errorf("you must perform some test prior to report")
		return
	}

	exitCode = 0
	for _, testSuite := range testSuites {
		suiteCtx := facade.NewSuiteContext(token, testSuite, false, model.ReportAction, ctx.Config)
		if suiteCtx.Repo.TestCount(testSuite) > 0 {
			code := 0
			code, err = ReportTestSuite(suiteCtx)
			if code != 0 {
				exitCode = code
			}
			if err != nil {
				// FIXME: aggregate errors
				return
			}
		}
	}
	dpl.ReportAllFooter(ctx)

	return
}

func ReportTestSuite(ctx facade.SuiteContext) (exitCode int, err error) {
	exitCode = 1
	cfg := ctx.Config
	testSuite := cfg.TestSuite.Get()
	if ctx.Repo.TestCount(testSuite) == 0 {
		err = fmt.Errorf("you must perform some test prior to report")
		return
	}

	logger.Info("Reporting suite", "ctx", ctx)

	var suiteOutcome model.SuiteOutcome
	suiteOutcome, err = ctx.Repo.LoadSuiteOutcome(testSuite)
	if err != nil {
		return
	}
	dpl.ReportSuite(ctx, suiteOutcome)

	if suiteOutcome.FailedCount == 0 && suiteOutcome.ErroredCount == 0 {
		exitCode = 0
	}

	err = ctx.Repo.ClearTestSuite(suiteOutcome.TestSuite)

	return
}

func PerformTest(testDef model.TestDefinition) (exitCode int, err error) {
	exitCode = 1
	cfg := testDef.Config
	ctx := facade.NewTestContext2(testDef)
	seq := testDef.Seq

	dpl.TestTitle(ctx, seq)

	if cfg.Ignore.Is(true) {
		ctx.IncrementIgnoredCount()
		dpl.TestOutcome(ctx, ctx.IgnoredTestOutcome(seq))
		exitCode = 0
		return
	}

	err = ctx.ConfigMocking()
	ctx.NoErrorOrFatal(err)

	for _, before := range cfg.Before {
		cmdBefore := cmdz.Cmd(before...)
		beforeExit, beforeErr := cmdBefore.BlockRun()
		// FIXME: what to do of before exit code or beforeErr ?
		_ = beforeExit
		if beforeErr != nil {
			err = fmt.Errorf("error running before cmd: [%s]: %w", cmdBefore.String(), beforeErr)
			return
		}
	}

	outcome := ctx.AssertCmdExecBlocking(seq, testDef.Assertions)

	dpl.TestOutcome(ctx, outcome)

	for _, after := range cfg.After {
		cmdAfter := cmdz.Cmd(after...)
		afterExit, afterErr := cmdAfter.BlockRun()
		// FIXME: what to do of after exit code or afterErr ?
		_ = afterExit
		if afterErr != nil {
			err = fmt.Errorf("error running after cmd: [%s]: %w", cmdAfter.String(), afterErr)
		}
	}

	if cfg.StopOnFailure.Is(true) && outcome.Outcome != model.PASSED {
		exitCode = max(outcome.ExitCode, 1)
		//ReportTestSuite(ctx)
		// FIXME do we need to call ReportTestSuite ?
	} else {
		exitCode = 0
	}
	return
}

func ProcessTestDef(testDef model.TestDefinition) (exitCode int) {
	token := testDef.Token
	testSuite := testDef.TestSuite
	testCfg := testDef.Config
	testCtx := facade.NewTestContext2(testDef)

	tooMuchFailures := testCtx.ProcessTooMuchFailures()

	if tooMuchFailures == 1 {
		// First time detecting TOO MUCH FAILURES
		dpl.TooMuchFailures(testSuite)
	}
	if tooMuchFailures > 0 {
		exitCode = 0
		return
	}

	var err error
	if !utils.IsWithinContainer() && (testCfg.ContainerDisabled.Is(true) || testCfg.ContainerImage.IsEmpty()) && len(testCfg.RootMocks) > 0 {
		err = fmt.Errorf("cannot mock absolute path outside a container")
		testCtx.NoErrorOrFatal(err)
	}

	logger.Info("Processing test action", "token", token)

	dpl.Quiet(testCfg.Quiet)
	if testCfg.ContainerDisabled.Is(true) || testCfg.ContainerImage.IsEmpty() {
		logger.Debug("Performing test outside container", "image", testCfg.ContainerImage, "containerDisabled", testCfg.ContainerDisabled, "testConfig", testCfg)
		exitCode, err = PerformTest(testDef)
		testCtx.NoErrorOrFatal(err)
	} else {
		logger.Info("Performing test inside container", "image", testCfg.ContainerImage, "containeriId", testCfg.ContainerId, "testConfig", testCfg)
		var ctId string
		ctId, exitCode, err = PerformTestInContainer(testCtx)
		testCtx.NoErrorOrFatal(err)
		if testCfg.ContainerScope.Is(model.GLOBAL_SCOPE) {
			globalCfg, err2 := testCtx.Repo.LoadGlobalConfig()
			testCtx.NoErrorOrFatal(err2)
			globalCfg.ContainerId.Set(ctId)
			err2 = testCtx.Repo.SaveGlobalConfig(globalCfg)
			testCtx.NoErrorOrFatal(err2)
		} else if testCfg.ContainerScope.Is(model.SUITE_SCOPE) {
			suiteCfg, err2 := testCtx.Repo.LoadSuiteConfig(testSuite, true)
			testCtx.NoErrorOrFatal(err2)
			suiteCfg.ContainerId.Set(ctId)
			err2 = testCtx.Repo.SaveSuiteConfig(suiteCfg)
			testCtx.NoErrorOrFatal(err2)
		}
		//testCtx.NoErrorOrFatal(err)
	}
	return
}

func ProcessArgs(allArgs []string) (exitCode int) {
	exitCode = 1

	if len(allArgs) == 1 {
		usage()
		return
	}

	// FIXME: if token supplied by ENV should be retrieved FIRST to get token and load Global config
	defaultCfg := model.NewGlobalDefaultConfig()
	envToken := utils.ReadEnvToken()
	if envToken != "" {
		envCtx := facade.NewGlobalContext(envToken, defaultCfg)
		defaultCfg = envCtx.Config
	}
	rulePrefix := defaultCfg.Prefix.Get()

	args := allArgs[1:]
	logger.Debug("Processing cmdt args", "args", args)
	inputConfig, assertions, agg := ParseArgs(rulePrefix, args)
	//loadedConfig := inputConfig
	inputConfig.Token.Default(envToken)
	logger.Debug("Parsed args", "inputConfig", inputConfig, "assertions", assertions, "error", agg)
	if inputConfig.Debug.IsPresent() {
		model.LoggerLevel.Set(slog.Level(8 - inputConfig.Debug.Get()*4))
	}

	token := inputConfig.Token.GetOr("")
	action := inputConfig.Action.Get()

	var err error
	switch action {
	case model.GlobalAction:
		if agg.GotError() {
			log.Fatal(agg)
		}
		globalCtx := facade.NewGlobalContext(token, inputConfig)
		logger.Debug("Forged context", "ctx", globalCtx)
		logger.Info("Processing global action", "token", token)
		dpl.Quiet(globalCtx.Config.Quiet)
		exitCode, err = GlobalConfig(globalCtx)
	case model.InitAction:
		testSuite := inputConfig.TestSuite.Get()
		suiteCtx := facade.NewSuiteContext(token, testSuite, false, action, inputConfig)
		logger.Debug("Forged context", "ctx", suiteCtx)
		suiteCtx.NoErrorOrFatal(agg.Return())
		logger.Info("Processing init action", "token", token)
		dpl.Quiet(suiteCtx.Config.Quiet)
		exitCode, err = InitTestSuite(suiteCtx)
	case model.ReportAction:
		if inputConfig.ReportAll.Is(true) {
			if agg.GotError() {
				log.Fatal(agg)
			}
			globalCtx := facade.NewGlobalContext(token, inputConfig)
			logger.Debug("Forged context", "ctx", globalCtx)
			logger.Info("Processing report all action", "token", token)
			dpl.Quiet(globalCtx.Config.Quiet)
			exitCode, err = ReportAllTestSuites(globalCtx)
		} else {
			testSuite := inputConfig.TestSuite.Get()
			suiteCtx := facade.NewSuiteContext(token, testSuite, false, action, inputConfig)
			logger.Debug("Forged context", "ctx", suiteCtx)
			suiteCtx.NoErrorOrFatal(agg.Return())
			logger.Info("Processing report suite action", "token", token)
			dpl.Quiet(suiteCtx.Config.Quiet)
			exitCode, err = ReportTestSuite(suiteCtx)
		}
	case model.TestAction:
		testSuite := inputConfig.TestSuite.Get()
		testCtx := facade.NewTestContext(token, testSuite, inputConfig)
		logger.Debug("Forged context", "ctx", testCtx)
		testCfg := testCtx.Config
		var seq int
		//if testCfg.ContainerDisabled.Is(true) || testCfg.ContainerImage.IsEmpty() {
		seq = testCtx.IncrementTestCount()
		//}

		testCtx.NoErrorOrFatal(agg.Return())

		testDef := model.TestDefinition{Token: token, TestSuite: testSuite, Seq: seq, Config: testCfg}
		testOp := repo.TestOperation{TestSuite: testSuite, Def: testDef, Blocking: !testCfg.Async.Get()}
		testCtx.Repo.QueueOperation(&testOp)

		//exitCode = processTestDef(testDef)
		exitCode = 0
	default:
		err = fmt.Errorf("action: [%v] not known", inputConfig.Action)
	}

	logger.Info("exiting", "exitCode", exitCode)

	if err != nil {
		log.Fatal(inputConfig.TestSuite, inputConfig.Token, err)
	}
	return
}
