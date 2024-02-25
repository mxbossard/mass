package model

import (
	"fmt"
	"regexp"
	"time"

	"mby.fr/utils/ansi"
	"mby.fr/utils/cmdz"
)

type Action string

const (
	GlobalAction = Action("global")
	InitAction   = Action("init")
	TestAction   = Action("test")
	ReportAction = Action("report")
)

type Mapper[T any] func(expr, op string) (T, error)

type Validater[T any] func(rule Rule, value T) error

//type Configurer func(ctx Context) (Context, error)

type Asserter func(cmdz.Executer) (AssertionResult, error)

const (
	DefaultRulePrefix             = "@"
	ContextTokenEnvVarName        = "__CMDT_TOKEN"
	DefaultTestSuiteName          = "main"
	GlobalConfigTestSuiteName     = "__global"
	DefaultContainerImage         = "busybox"
	DefaultContainerDirtiesPolicy = "beforeSuite"
)

const (
	NamePattern = "[a-zA-Z][^/]*[a-zA-Z0-9]"
)

var (
	TempDirPrefix           = "cmdtest"
	ContextFilename         = "context.yaml"
	TestSequenceFilename    = "test-seq.txt"
	PassedSequenceFilename  = "passed-seq.txt"
	FailedSequenceFilename  = "failure-seq.txt"
	IgnoredSequenceFilename = "ignored-seq.txt"
	ErroredSequenceFilename = "error-seq.txt"
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
		ruleDef("keepOutputs", "", "="), ruleDef("silent", "", "="), ruleDef("timeout", "="),
		ruleDef("parallel", "="), ruleDef("runCount", "="), ruleDef("mock", "=", ":"),
		ruleDef("before", "="), ruleDef("after", "="), ruleDef("container", "", "="), ruleDef("dirtyContainer", "=")}
	// Config of test flow (init -> test -> report)
	FlowConfigs = []RuleDefinition{ruleDef("token", "=")}
	Assertions  = []RuleDefinition{ruleDef("success", ""), ruleDef("fail", ""), ruleDef("exit", "="),
		ruleDef("cmd", "="), ruleDef("stdout", "=", ":", "~", "!=", "!:", "!~"),
		ruleDef("stderr", "=", ":", "~", "!=", "!:", "!~"), ruleDef("exists", "=")}
)

type Rule struct {
	Prefix   string
	Name     string
	Op       string
	Expected string
}

type RuleKey struct {
	Name, Op string
}

func (r RuleKey) String() string {
	return fmt.Sprintf("%s%s", r.Name, r.Op)
}

type RuleDefinition struct {
	Name string
	Ops  []string
}

type Context0 struct {
	// TestSuite only
	Token        string        `yaml:""`
	Prefix       string        `yaml:""`
	TestSuite    string        `yaml:""`
	TestName     string        `yaml:""`
	Action       Action        `yaml:""`
	StartTime    time.Time     `yaml:""`
	LastTestTime time.Time     `yaml:""`
	SuiteTimeout time.Duration `yaml:""`
	ForkCount    int           `yaml:""`

	// Test or TestSuite
	PrintToken        bool
	ExportToken       bool
	ReportAll         bool
	Silent            *bool         `yaml:""`
	Ignore            *bool         `yaml:""`
	StopOnFailure     *bool         `yaml:""`
	KeepStdout        *bool         `yaml:""`
	KeepStderr        *bool         `yaml:""`
	Timeout           time.Duration `yaml:""`
	RunCount          int           `yaml:""`
	Parallel          int           `yaml:""`
	Mocks             []CmdMock     `yaml:""`
	Before            [][]string    `yaml:""`
	After             [][]string    `yaml:""`
	ContainerDisabled *bool         `yaml:""`
	ContainerImage    string        `yaml:""`
	ContainerDirties  string        `yaml:""`
	ContainerId       *string       `yaml:""`
	ContainerScope    *ConfigScope  `yaml:""`
}

func (c Context0) String() string {
	keepStdout := c.KeepStdout != nil && *c.KeepStdout
	keepStderr := c.KeepStderr != nil && *c.KeepStderr
	return fmt.Sprintf("[%s/%s] KeepStdout: %v, KeepStderr: %v", c.TestSuite, c.TestName, keepStdout, keepStderr)
}

type Config0 struct {
	Name  string
	Scope ConfigScope
	Value string
}

/*
func MergeContext(baseContext, overridingContext Context) Context {
	baseContext.Token = overridingContext.Token
	baseContext.Prefix = overridingContext.Prefix
	baseContext.TestName = overridingContext.TestName
	baseContext.Action = overridingContext.Action
	if !overridingContext.StartTime.IsZero() {
		baseContext.StartTime = overridingContext.StartTime
	}
	if !overridingContext.LastTestTime.IsZero() {
		baseContext.LastTestTime = overridingContext.LastTestTime
	}

	if baseContext.TestSuite == "" {
		baseContext.TestSuite = overridingContext.TestSuite
	}

	if overridingContext.Ignore != nil {
		baseContext.Ignore = overridingContext.Ignore
	}
	if baseContext.Ignore == nil {
		baseContext.Ignore = ptr.BoolPtr(false)
	}
	if overridingContext.StopOnFailure != nil {
		baseContext.StopOnFailure = overridingContext.StopOnFailure
	}
	if baseContext.StopOnFailure == nil {
		baseContext.StopOnFailure = ptr.BoolPtr(false)
	}
	if overridingContext.KeepStdout != nil {
		baseContext.KeepStdout = overridingContext.KeepStdout
	}
	if baseContext.KeepStdout == nil {
		baseContext.KeepStdout = ptr.BoolPtr(false)
	}
	if overridingContext.KeepStderr != nil {
		baseContext.KeepStderr = overridingContext.KeepStderr
	}
	if baseContext.KeepStderr == nil {
		baseContext.KeepStderr = ptr.BoolPtr(false)
	}
	if overridingContext.Timeout.Nanoseconds() > 0 {
		baseContext.Timeout = overridingContext.Timeout
	}
	if overridingContext.RunCount != 0 {
		baseContext.RunCount = overridingContext.RunCount
	}
	if overridingContext.Parallel != 0 {
		baseContext.Parallel = overridingContext.Parallel
	}
	if overridingContext.Silent != nil {
		baseContext.Silent = overridingContext.Silent
	}
	baseContext.Mocks = append(baseContext.Mocks, overridingContext.Mocks...)
	baseContext.Before = append(baseContext.Before, overridingContext.Before...)
	baseContext.After = append(baseContext.After, overridingContext.After...)

	if overridingContext.ContainerImage != "" {
		baseContext.ContainerImage = overridingContext.ContainerImage
	}
	if overridingContext.ContainerDirties != "" {
		baseContext.ContainerDirties = overridingContext.ContainerDirties
	}
	if overridingContext.ContainerId != nil {
		baseContext.ContainerId = overridingContext.ContainerId
	}
	if overridingContext.ContainerScope != nil {
		baseContext.ContainerScope = overridingContext.ContainerScope
	}
	if overridingContext.ContainerDisabled != nil {
		baseContext.ContainerDisabled = overridingContext.ContainerDisabled
	}

	return baseContext
}
*/

func ruleDef(name string, ops ...string) (r RuleDefinition) {
	r.Name = name
	r.Ops = ops
	return
}

func IsRuleOfKind(ruleDefs []RuleDefinition, r Rule) (ok bool, err error) {
	ok = false
	var expectedOperators []string
	for _, ruleDef := range ruleDefs {
		if r.Name == ruleDef.Name {
			expectedOperators = append(expectedOperators, ruleDef.Ops...)
			for _, op := range ruleDef.Ops {
				if r.Op == op {
					ok = true
					return
				}
			}
		}
	}

	if len(expectedOperators) > 0 {
		// name matched but not operator
		err = fmt.Errorf("rule %s expect one of operators: %s", r.Name, expectedOperators)
	}
	return
}

func NormalizeDurationInSec(d time.Duration) (duration string) {
	duration = fmt.Sprintf("%.3f s", float32(d.Milliseconds())/1000)
	return
}
