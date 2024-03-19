package model

import (
	"fmt"
	"log/slog"
	"regexp"
	"strings"
	"time"

	"mby.fr/utils/ansi"
	"mby.fr/utils/utilz"
)

const (
	DefaultRulePrefix           = "@"
	DefaultVerboseLevel         = SHOW_PASSED
	DefaultInitedVerboseLevel   = SHOW_FAILED_OUTS
	DefaultInitlessVerboseLevel = DefaultVerboseLevel
	DefaultInitedAsync          = true
	DefaultInitlessAsync        = false
	StartDebugLevel             = WARN
	DefaultDebugLevel           = INFO
	DefaultTooMuchFailures      = 3
	TooMuchFailuresNoLimit      = -1

	ContextTokenEnvVarName    = "__CMDT_TOKEN"
	GlobalConfigTestSuiteName = "__global"
	DefaultTestSuiteName      = "main"

	DefaultContainerDirtiesPolicy = "beforeSuite"
	DefaultContainerImage         = "busybox"
	EnvContainerScopeKey          = "__CMDT_CONTAINER_SCOPE"
	EnvContainerImageKey          = "__CMDT_CONTAINER_IMAGE"
	EnvContainerIdKey             = "__CMDT_CONTAINER_ID"

	NamePattern = "[a-zA-Z][^/]*[a-zA-Z0-9]"

	TempDirPrefix           = "cmdtest"
	ContextFilename         = "context.yaml"
	TestSequenceFilename    = "test-seq.txt"
	PassedSequenceFilename  = "passed-seq.txt"
	FailedSequenceFilename  = "failed-seq.txt"
	IgnoredSequenceFilename = "ignored-seq.txt"
	ErroredSequenceFilename = "errored-seq.txt"
	TooMuchSequenceFilename = "tooMuch-seq.txt"
	StdoutFilename          = "stdout.log"
	StderrFilename          = "stderr.log"
	ReportFilename          = "report.log"

	MessageColor = ansi.HiPurple
	TestColor    = ansi.HiCyan
	SuccessColor = ansi.BoldGreen
	FailureColor = ansi.BoldRed
	ReportColor  = ansi.Yellow
	WarningColor = ansi.BoldHiYellow
	ErrorColor   = ansi.Red
)

const (
	GlobalAction = Action("global")
	InitAction   = Action("init")
	TestAction   = Action("test")
	ReportAction = Action("report")
)

var (
	LoggerLevel       slog.LevelVar
	DefaultLoggerOpts = &slog.HandlerOptions{
		Level: &LoggerLevel,
	}
)

var (
	DefaultTestTimeout, _ = time.ParseDuration("24h")
	AbsNamePattern        = fmt.Sprintf("(%s/)?(%s)?", NamePattern, NamePattern)
	NameRegexp            = regexp.MustCompile("^" + NamePattern + "$")
	AbsNameRegexp         = regexp.MustCompile("^" + AbsNamePattern + "$")
)

var (
	Actions = []RuleDefinition{ruleDef("global", ""), ruleDef("init", "", "="), ruleDef("test", "", "="),
		ruleDef("report", "", "=")}
	// Global config available to global
	GlobalConfigs = []RuleDefinition{ruleDef("fork", "="), ruleDef("suiteTimeout", "="), ruleDef("prefix", "=")}
	// Suite config available to suite
	SuiteConfigs = append(GlobalConfigs, []RuleDefinition{
		ruleDef("exportToken", ""), ruleDef("printToken", "")}...)
	// Test config available at all levels (global, suite and test)
	TestConfigs = []RuleDefinition{ruleDef("before", "="), ruleDef("after", "="), ruleDef("ignore", "", "="),
		ruleDef("stopOnFailure", "", "="), ruleDef("keepStdout", "", "="), ruleDef("keepStderr", "", "="),
		ruleDef("keepOutputs", "", "="), ruleDef("quiet", "", "="), ruleDef("timeout", "="),
		ruleDef("parallel", "="), ruleDef("runCount", "="), ruleDef("mock", "=", ":"),
		ruleDef("container", "", "="), ruleDef("dirtyContainer", "=")}
	// Config of test flow (init -> test -> report)
	FlowConfigs = []RuleDefinition{ruleDef("token", "="), ruleDef("verbose", "", "="),
		ruleDef("debug", "", "="), ruleDef("failuresLimit", "="), ruleDef("async", "", "=")}
	Assertions = []RuleDefinition{ruleDef("success", ""), ruleDef("fail", ""), ruleDef("exit", "="),
		ruleDef("cmd", "="), ruleDef("exists", "="),
		ruleDef("stdout", "=", ":", "~", "!=", "!:", "!~", "@=", "@:"),
		ruleDef("stderr", "=", ":", "~", "!=", "!:", "!~", "@=", "@:")}
	Concatenables = []RuleDefinition{
		ruleDef("init", "="), ruleDef("test", "="), ruleDef("report", "="),
		ruleDef("before", "="), ruleDef("after", "="),
		ruleDef("cmd", "="), ruleDef("exists", "="),
		ruleDef("stdout", "=", ":", "~", "!=", "!:", "!~", "@=", "@:"),
		ruleDef("stderr", "=", ":", "~", "!=", "!:", "!~", "@=", "@:"),
	}
)

func NewGlobalDefaultConfig() Config {
	return Config{
		Prefix: utilz.OptionalOf(DefaultRulePrefix),
		//Verbose: utilz.OptionalOf(DefaultInitlessVerboseLevel),
		//Verbose:           utilz.OptionalOf(DefaultInitedVerboseLevel),
		GlobalStartTime:   utilz.OptionalOf(time.Now()),
		ForkCount:         utilz.OptionalOf(1),
		Ignore:            utilz.OptionalOf(false),
		StopOnFailure:     utilz.OptionalOf(false),
		KeepStdout:        utilz.OptionalOf(false),
		KeepStderr:        utilz.OptionalOf(false),
		Timeout:           utilz.OptionalOf(10 * time.Second),
		RunCount:          utilz.OptionalOf(1),
		Parallel:          utilz.OptionalOf(1),
		ContainerDisabled: utilz.OptionalOf(true),
		ContainerImage:    utilz.OptionalOf(DefaultContainerImage),
		ContainerDirties:  utilz.OptionalOf(DirtyBeforeSuite),
		ContainerScope:    utilz.OptionalOf(SUITE_SCOPE),
	}
}

func NewSuiteDefaultConfig() Config {
	return Config{
		Async:           utilz.OptionalOf(DefaultInitedAsync),
		TooMuchFailures: utilz.OptionalOf(DefaultTooMuchFailures),
		SuiteStartTime:  utilz.OptionalOf(time.Now()),
		SuiteTimeout:    utilz.OptionalOf(120 * time.Second),
		Verbose:         utilz.OptionalOf(DefaultInitedVerboseLevel),
	}
}

func NewInitlessSuiteDefaultConfig() Config {
	return Config{
		Async:           utilz.OptionalOf(DefaultInitlessAsync),
		TooMuchFailures: utilz.OptionalOf(TooMuchFailuresNoLimit),
		SuiteStartTime:  utilz.OptionalOf(time.Now()),
		SuiteTimeout:    utilz.OptionalOf(3600 * time.Second),
		Verbose:         utilz.OptionalOf(DefaultInitlessVerboseLevel),
	}
}

type ConfigScope int

const (
	GLOBAL_SCOPE ConfigScope = iota // How to use this ?
	SUITE_SCOPE                     // can be placed on suite init only
	TEST_SCOPE                      // can be placed on test or on suite to configure all tests
	RUN_SCOPE
)

type DirtyScope string

const (
	DirtyBeforeSuite = DirtyScope("beforeSuite")
	DirtyAfterSuite  = DirtyScope("afterSuite")
	DirtyBeforeTest  = DirtyScope("beforeTest")
	DirtyAfterTest   = DirtyScope("afterTest")
	DirtyBeforeRun   = DirtyScope("beforeRun")
	DirtyAfterRun    = DirtyScope("afterRun")
)

type VerboseLevel int

const (
	SHOW_FAILED_ONLY VerboseLevel = iota
	SHOW_FAILED_OUTS
	SHOW_PASSED
	SHOW_PASSED_OUTS
	SHOW_ALL
)

type DebugLevel int

const (
	ERROR DebugLevel = iota
	WARN
	INFO
	DEBUG
	TRACE
)

type CmdMock struct {
	Op               string
	Cmd              string
	Args             []string
	StdinOp          string
	Stdin            *string
	Stdout           string
	Stderr           string
	ExitCode         int
	Delegate         bool
	OnCallCmdAndArgs []string
}

type Config struct {
	// TestSuite only
	Token     utilz.Optional[string]
	Action    utilz.Optional[Action] `yaml:""`
	TestSuite utilz.Optional[string] `yaml:""`
	TestName  utilz.Optional[string] `yaml:""`
	Async     utilz.Optional[bool]   `yaml:""`

	Prefix          utilz.Optional[string]        `yaml:""`
	CmdAndArgs      []string                      `yaml:""`
	GlobalStartTime utilz.Optional[time.Time]     `yaml:""`
	SuiteStartTime  utilz.Optional[time.Time]     `yaml:""`
	TooMuchFailures utilz.Optional[int]           `yaml:""`
	LastTestTime    utilz.Optional[time.Time]     `yaml:""`
	SuiteTimeout    utilz.Optional[time.Duration] `yaml:""`
	ForkCount       utilz.Optional[int]           `yaml:""`
	BeforeSuite     [][]string                    `yaml:""`
	AfterSuite      [][]string                    `yaml:""`

	// Test or TestSuite
	PrintToken    utilz.Optional[bool]          `yaml:""`
	ExportToken   utilz.Optional[bool]          `yaml:""`
	ReportAll     utilz.Optional[bool]          `yaml:""`
	Verbose       utilz.Optional[VerboseLevel]  `yaml:""`
	Debug         utilz.Optional[DebugLevel]    `yaml:""`
	Quiet         utilz.Optional[bool]          `yaml:""`
	Ignore        utilz.Optional[bool]          `yaml:""`
	StopOnFailure utilz.Optional[bool]          `yaml:""`
	KeepStdout    utilz.Optional[bool]          `yaml:""`
	KeepStderr    utilz.Optional[bool]          `yaml:""`
	Timeout       utilz.Optional[time.Duration] `yaml:""`
	RunCount      utilz.Optional[int]           `yaml:""`
	Parallel      utilz.Optional[int]           `yaml:""`

	Mocks     []CmdMock  `yaml:""`
	RootMocks []CmdMock  `yaml:""`
	Before    [][]string `yaml:""`
	After     [][]string `yaml:""`

	ContainerDisabled utilz.Optional[bool]        `yaml:""`
	ContainerImage    utilz.Optional[string]      `yaml:""`
	ContainerDirties  utilz.Optional[DirtyScope]  `yaml:""`
	ContainerId       utilz.Optional[string]      `yaml:""`
	ContainerScope    utilz.Optional[ConfigScope] `yaml:""`
}

func (c Config) IsRule(s string) bool {
	prefix := c.Prefix.Get()
	return strings.HasPrefix(s, prefix)
}

func (c Config) SplitRuleExpr(ruleExpr string) (ok bool, r Rule) {
	ok = false
	prefix := c.Prefix.Get()
	assertionRulePattern := regexp.MustCompile("^" + prefix + "([a-zA-Z]+)([=~:!@]{1,2})?(.+)?$")
	submatch := assertionRulePattern.FindStringSubmatch(ruleExpr)
	if submatch != nil {
		ok = true
		r.Prefix = prefix
		r.Name = submatch[1]
		r.Op = submatch[2]
		r.Expected = submatch[3]
	}
	return
}

func (c *Config) Merge(right Config) {
	c.Action.Merge(right.Action)
	c.TestSuite.Merge(right.TestSuite)
	c.TestSuite.Merge(right.TestSuite)
	c.TestName.Merge(right.TestName)
	c.Async.Merge(right.Async)

	c.Prefix.Merge(right.Prefix)
	c.TooMuchFailures.Merge(right.TooMuchFailures)
	c.GlobalStartTime.Merge(right.GlobalStartTime)
	c.SuiteStartTime.Merge(right.SuiteStartTime)
	c.LastTestTime.Merge(right.LastTestTime)
	c.SuiteTimeout.Merge(right.SuiteTimeout)
	c.ForkCount.Merge(right.ForkCount)
	if len(right.CmdAndArgs) > 0 {
		c.CmdAndArgs = right.CmdAndArgs
	}
	if len(right.BeforeSuite) > 0 {
		c.BeforeSuite = append(c.BeforeSuite, right.BeforeSuite...)
	}
	if len(right.AfterSuite) > 0 {
		c.AfterSuite = append(c.AfterSuite, right.AfterSuite...)
	}

	c.PrintToken.Merge(right.PrintToken)
	c.ExportToken.Merge(right.ExportToken)
	c.ReportAll.Merge(right.ReportAll)
	c.Quiet.Merge(right.Quiet)
	c.Ignore.Merge(right.Ignore)
	c.StopOnFailure.Merge(right.StopOnFailure)
	c.KeepStdout.Merge(right.KeepStdout)
	c.KeepStderr.Merge(right.KeepStderr)
	c.Timeout.Merge(right.Timeout)
	c.RunCount.Merge(right.RunCount)
	c.Parallel.Merge(right.Parallel)
	c.Verbose.Merge(right.Verbose)
	c.Debug.Merge(right.Debug)

	if len(right.Mocks) > 0 {
		c.Mocks = append(c.Mocks, right.Mocks...)
	}
	if len(right.RootMocks) > 0 {
		c.RootMocks = append(c.RootMocks, right.RootMocks...)
	}
	if len(right.Before) > 0 {
		c.Before = append(c.Before, right.Before...)
	}
	if len(right.After) > 0 {
		c.After = append(c.After, right.After...)
	}

	c.ContainerDisabled.Merge(right.ContainerDisabled)
	c.ContainerImage.Merge(right.ContainerImage)
	c.ContainerDirties.Merge(right.ContainerDirties)
	c.ContainerId.Merge(right.ContainerId)
	c.ContainerScope.Merge(right.ContainerScope)
}
