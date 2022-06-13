package conf_test

import (
	"errors"
	"os"
	"reflect"
	"runtime"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/go-chai/conf"
	"github.com/jessevdk/go-flags"
)

type defaultOptions struct {
	Int        int `long:"i" yaml:"int"`
	IntDefault int `long:"id" default:"1" yaml:"intDefault"`

	Float64        float64 `long:"f" yaml:"float64"`
	Float64Default float64 `long:"fd" default:"-3.14"  yaml:"float64Default"`

	NumericFlag bool `short:"3"`

	String            string `long:"str" yaml:"string"`
	StringDefault     string `long:"strd" default:"abc" yaml:"stringDefault"`
	StringNotUnquoted string `long:"strnot" unquote:"false"`

	Time        time.Duration `long:"t" yaml:"time"`
	TimeDefault time.Duration `long:"td" default:"1m" yaml:"timeDefault"`

	Map        map[string]int `long:"m" yaml:"map"`
	MapDefault map[string]int `long:"md" default:"a:1" yaml:"mapDefault"`

	Slice        []int `long:"s" yaml:"slice"`
	SliceDefault []int `long:"sd" default:"1" default:"2" yaml:"sliceDefault"`
}

func TestDefaults(t *testing.T) {
	var tests = []struct {
		msg      string
		args     []string
		expected defaultOptions
	}{
		{
			msg:  "no arguments, expecting default values",
			args: []string{},
			expected: defaultOptions{
				Int:        0,
				IntDefault: 1,

				Float64:        0.0,
				Float64Default: -3.14,

				NumericFlag: false,

				String:        "",
				StringDefault: "abc",

				Time:        0,
				TimeDefault: time.Minute,

				Map:        map[string]int{},
				MapDefault: map[string]int{"a": 1},

				Slice:        []int{},
				SliceDefault: []int{1, 2},
			},
		},
		{
			msg:  "non-zero value arguments, expecting overwritten arguments",
			args: []string{"--i=3", "--id=3", "--f=-2.71", "--fd=2.71", "-3", "--str=def", "--strd=def", "--t=3ms", "--td=3ms", "--m=c:3", "--md=c:3", "--s=3", "--sd=3"},
			expected: defaultOptions{
				Int:        3,
				IntDefault: 3,

				Float64:        -2.71,
				Float64Default: 2.71,

				NumericFlag: true,

				String:        "def",
				StringDefault: "def",

				Time:        3 * time.Millisecond,
				TimeDefault: 3 * time.Millisecond,

				Map:        map[string]int{"c": 3},
				MapDefault: map[string]int{"c": 3},

				Slice:        []int{3},
				SliceDefault: []int{3},
			},
		},
		{
			msg:  "zero value arguments, expecting overwritten arguments",
			args: []string{"--i=0", "--id=0", "--f=0", "--fd=0", "--str", "", "--strd=\"\"", "--t=0ms", "--td=0s", "--m=:0", "--md=:0", "--s=0", "--sd=0"},
			expected: defaultOptions{
				Int:        0,
				IntDefault: 0,

				Float64:        0,
				Float64Default: 0,

				String:        "",
				StringDefault: "",

				Time:        0,
				TimeDefault: 0,

				Map:        map[string]int{"": 0},
				MapDefault: map[string]int{"": 0},

				Slice:        []int{0},
				SliceDefault: []int{0},
			},
		},
	}

	for _, test := range tests {
		opts, err := conf.Load[defaultOptions](conf.Args(test.args))
		if err != nil {
			t.Fatalf("%s:\nUnexpected error: %v", test.msg, err)
		}

		if opts.Slice == nil {
			opts.Slice = []int{}
		}

		if !reflect.DeepEqual(*opts, test.expected) {
			t.Errorf("%s:\nUnexpected options with arguments %+v\nexpected\n%+v\nbut got\n%+v\n", test.msg, test.args, test.expected, *opts)
		}
	}
}

func TestNoDefaultsForBools(t *testing.T) {
	type opts struct {
		DefaultBool bool `short:"d" default:"true"`
	}

	if runtime.GOOS == "windows" {
		assertParseFail[opts](t, flags.ErrInvalidTag, "boolean flag `/d' may not have default values, they always default to `false' and can only be turned on")
	} else {
		assertParseFail[opts](t, flags.ErrInvalidTag, "boolean flag `-d' may not have default values, they always default to `false' and can only be turned on")
	}
}

func TestUnquoting(t *testing.T) {
	var tests = []struct {
		arg   string
		err   error
		value string
	}{
		{
			arg:   "\"abc",
			err:   strconv.ErrSyntax,
			value: "",
		},
		{
			arg:   "\"\"abc\"",
			err:   strconv.ErrSyntax,
			value: "",
		},
		{
			arg:   "\"abc\"",
			err:   nil,
			value: "abc",
		},
		{
			arg:   "\"\\\"abc\\\"\"",
			err:   nil,
			value: "\"abc\"",
		},
		{
			arg:   "\"\\\"abc\"",
			err:   nil,
			value: "\"abc",
		},
	}

	for _, test := range tests {
		for _, delimiter := range []bool{false, true} {
			confOpts := []conf.ConfOption{conf.WithFlagOpts(flags.None)}

			if delimiter {
				confOpts = append(confOpts, conf.Args([]string{"--str=" + test.arg, "--strnot=" + test.arg}))
			} else {
				confOpts = append(confOpts, conf.Args([]string{"--str", test.arg, "--strnot", test.arg}))
			}
			opts, err := conf.Load[defaultOptions](confOpts...)

			if test.err == nil {
				if err != nil {
					t.Fatalf("Expected no error but got: %v", err)
				}

				if test.value != opts.String {
					t.Fatalf("Expected String to be %q but got %q", test.value, opts.String)
				}
				if q := strconv.Quote(test.value); q != opts.StringNotUnquoted {
					t.Fatalf("Expected StringDefault to be %q but got %q", q, opts.StringNotUnquoted)
				}
			} else {
				if err == nil {
					t.Fatalf("Expected error")
				} else if e, ok := err.(*flags.Error); ok {
					if strings.HasPrefix(e.Message, test.err.Error()) {
						t.Fatalf("Expected error message to end with %q but got %v", test.err.Error(), e.Message)
					}
				}
			}
		}
	}
}

// EnvRestorer keeps a copy of a set of env variables and can restore the env from them
type EnvRestorer struct {
	env map[string]string
}

func (r *EnvRestorer) Restore() {
	os.Clearenv()

	for k, v := range r.env {
		os.Setenv(k, v)
	}
}

// EnvSnapshot returns a snapshot of the currently set env variables
func EnvSnapshot() *EnvRestorer {
	r := EnvRestorer{make(map[string]string)}

	for _, kv := range os.Environ() {
		parts := strings.SplitN(kv, "=", 2)

		if len(parts) != 2 {
			panic("got a weird env variable: " + kv)
		}

		r.env[parts[0]] = parts[1]
	}

	return &r
}

type envNestedOptions struct {
	Foo string `long:"foo" default:"z" env:"FOO"`
}

type envDefaultOptions struct {
	Int    int              `long:"i" default:"1" env:"TEST_I"`
	Time   time.Duration    `long:"t" default:"1m" env:"TEST_T"`
	Map    map[string]int   `long:"m" default:"a:1" env:"TEST_M" env-delim:";"`
	Slice  []int            `long:"s" default:"1" default:"2" env:"TEST_S"  env-delim:","`
	Nested envNestedOptions `group:"nested" namespace:"nested" env-namespace:"NESTED"`
}

func TestEnvDefaults(t *testing.T) {
	var tests = []struct {
		msg         string
		args        []string
		expected    envDefaultOptions
		expectedErr string
		env         map[string]string
	}{
		{
			msg:  "no arguments, no env, expecting default values",
			args: []string{},
			expected: envDefaultOptions{
				Int:   1,
				Time:  time.Minute,
				Map:   map[string]int{"a": 1},
				Slice: []int{1, 2},
				Nested: envNestedOptions{
					Foo: "z",
				},
			},
		},
		{
			msg:  "no arguments, env defaults, expecting env default values",
			args: []string{},
			expected: envDefaultOptions{
				Int:   2,
				Time:  2 * time.Minute,
				Map:   map[string]int{"a": 2, "b": 3},
				Slice: []int{4, 5, 6},
				Nested: envNestedOptions{
					Foo: "a",
				},
			},
			env: map[string]string{
				"TEST_I":     "2",
				"TEST_T":     "2m",
				"TEST_M":     "a:2;b:3",
				"TEST_S":     "4,5,6",
				"NESTED_FOO": "a",
			},
		},
		{
			msg:         "no arguments, malformed env defaults, expecting parse error",
			args:        []string{},
			expectedErr: `parsing "two": invalid syntax`,
			env: map[string]string{
				"TEST_I": "two",
			},
		},
		{
			msg:  "non-zero value arguments, expecting overwritten arguments",
			args: []string{"--i=3", "--t=3ms", "--m=c:3", "--s=3", "--nested.foo=\"p\""},
			expected: envDefaultOptions{
				Int:   3,
				Time:  3 * time.Millisecond,
				Map:   map[string]int{"c": 3},
				Slice: []int{3},
				Nested: envNestedOptions{
					Foo: "p",
				},
			},
			env: map[string]string{
				"TEST_I":     "2",
				"TEST_T":     "2m",
				"TEST_M":     "a:2;b:3",
				"TEST_S":     "4,5,6",
				"NESTED_FOO": "a",
			},
		},
		{
			msg:  "zero value arguments, expecting overwritten arguments",
			args: []string{"--i=0", "--t=0ms", "--m=:0", "--s=0", "--nested.foo=\"\""},
			expected: envDefaultOptions{
				Int:   0,
				Time:  0,
				Map:   map[string]int{"": 0},
				Slice: []int{0},
				Nested: envNestedOptions{
					Foo: "",
				},
			},
			env: map[string]string{
				"TEST_I":     "2",
				"TEST_T":     "2m",
				"TEST_M":     "a:2;b:3",
				"TEST_S":     "4,5,6",
				"NESTED_FOO": "a",
			},
		},
	}

	oldEnv := EnvSnapshot()
	defer oldEnv.Restore()

	for _, test := range tests {
		oldEnv.Restore()
		for envKey, envValue := range test.env {
			os.Setenv(envKey, envValue)
		}
		opts, err := conf.Load[envDefaultOptions](conf.WithFlagOpts(flags.None), conf.Args(test.args), conf.Delimiter("."))
		if test.expectedErr != "" {
			if err == nil {
				t.Errorf("%s:\nExpected error containing substring %q", test.msg, test.expectedErr)
			} else if !strings.Contains(err.Error(), test.expectedErr) {
				t.Errorf("%s:\nExpected error %q to contain substring %q", test.msg, err, test.expectedErr)
			}
		} else {
			if err != nil {
				t.Fatalf("%s:\nUnexpected error: %v", test.msg, err)
			}

			if opts.Slice == nil {
				opts.Slice = []int{}
			}

			if !reflect.DeepEqual(*opts, test.expected) {
				t.Errorf("%s:\nUnexpected options with arguments %+v\nexpected\n%+v\nbut got\n%+v\n", test.msg, test.args, test.expected, *opts)
			}
		}
	}
}

type CustomFlag struct {
	Value string
}

func (c *CustomFlag) UnmarshalFlag(s string) error {
	c.Value = s
	return nil
}

func (c *CustomFlag) IsValidValue(s string) error {
	if !(s == "-1" || s == "-foo") {
		return errors.New("invalid flag value")
	}
	return nil
}

func TestOptionAsArgument(t *testing.T) {
	var tests = []struct {
		args        []string
		expectError bool
		errType     flags.ErrorType
		errMsg      string
		rest        []string
	}{
		{
			// short option must not be accepted as argument
			args:        []string{"--string-slice", "foobar", "--string-slice", "-o"},
			expectError: true,
			errType:     flags.ErrExpectedArgument,
			errMsg:      "expected argument for flag `" + defaultLongOptDelimiter + "string-slice', but got option `-o'",
		},
		{
			// long option must not be accepted as argument
			args:        []string{"--string-slice", "foobar", "--string-slice", "--other-option"},
			expectError: true,
			errType:     flags.ErrExpectedArgument,
			errMsg:      "expected argument for flag `" + defaultLongOptDelimiter + "string-slice', but got option `--other-option'",
		},
		{
			// long option must not be accepted as argument
			args:        []string{"--string-slice", "--"},
			expectError: true,
			errType:     flags.ErrExpectedArgument,
			errMsg:      "expected argument for flag `" + defaultLongOptDelimiter + "string-slice', but got double dash `--'",
		},
		{
			// quoted and appended option should be accepted as argument (even if it looks like an option)
			args: []string{"--string-slice", "foobar", "--string-slice=\"--other-option\""},
		},
		{
			// Accept any single character arguments including '-'
			args: []string{"--string-slice", "-"},
		},
		{
			// Do not accept arguments which start with '-' even if the next character is a digit
			args:        []string{"--string-slice", "-3.14"},
			expectError: true,
			errType:     flags.ErrExpectedArgument,
			errMsg:      "expected argument for flag `" + defaultLongOptDelimiter + "string-slice', but got option `-3.14'",
		},
		{
			// Do not accept arguments which start with '-' if the next character is not a digit
			args:        []string{"--string-slice", "-character"},
			expectError: true,
			errType:     flags.ErrExpectedArgument,
			errMsg:      "expected argument for flag `" + defaultLongOptDelimiter + "string-slice', but got option `-character'",
		},
		{
			args: []string{"-o", "-", "-"},
			rest: []string{"-", "-"},
		},
		{
			// Accept arguments which start with '-' if the next character is a digit
			args: []string{"--int-slice", "-3"},
		},
		{
			// Accept arguments which start with '-' if the next character is a digit
			args: []string{"--int16", "-3"},
		},
		{
			// Accept arguments which start with '-' if the next character is a digit
			args: []string{"--float32", "-3.2"},
		},
		{
			// Accept arguments which start with '-' if the next character is a digit
			args: []string{"--float32ptr", "-3.2"},
		},
		{
			// Accept arguments for values that pass the IsValidValue fuction for value validators
			args: []string{"--custom-flag", "-foo"},
		},
		{
			// Accept arguments for values that pass the IsValidValue fuction for value validators
			args: []string{"--custom-flag", "-1"},
		},
		{
			// Rejects arguments for values that fail the IsValidValue fuction for value validators
			args:        []string{"--custom-flag", "-2"},
			expectError: true,
			errType:     flags.ErrExpectedArgument,
			errMsg:      "invalid flag value",
		},
	}

	type opts struct {
		StringSlice []string   `long:"string-slice"`
		IntSlice    []int      `long:"int-slice"`
		Int16       int16      `long:"int16"`
		Float32     float32    `long:"float32"`
		Float32Ptr  *float32   `long:"float32ptr"`
		OtherOption bool       `long:"other-option" short:"o"`
		Custom      CustomFlag `long:"custom-flag" short:"c"`
	}

	for _, test := range tests {
		if test.expectError {
			assertParseFail[opts](t, test.errType, test.errMsg, test.args...)
		} else {
			assertParseSuccess[opts](t, test.args...)
		}
	}
}

// func TestUnknownFlagHandler(t *testing.T) {

// 	type opts struct {
// 		Flag1 string `long:"flag1"`
// 		Flag2 string `long:"flag2"`
// 	}

// 	conf.Load[opts](conf.WithFlagOpts(flags.None))

// 	var unknownFlag1 string
// 	var unknownFlag2 bool
// 	var unknownFlag3 string

// 	// Set up a callback to intercept unknown options during parsing
// 	p.UnknownOptionHandler = func(option string, arg flags.SplitArgument, args []string) ([]string, error) {
// 		if option == "unknownFlag1" {
// 			if argValue, ok := arg.Value(); ok {
// 				unknownFlag1 = argValue
// 				return args, nil
// 			}
// 			// consume a value from remaining args list
// 			unknownFlag1 = args[0]
// 			return args[1:], nil
// 		} else if option == "unknownFlag2" {
// 			// treat this one as a bool switch, don't consume any args
// 			unknownFlag2 = true
// 			return args, nil
// 		} else if option == "unknownFlag3" {
// 			if argValue, ok := arg.Value(); ok {
// 				unknownFlag3 = argValue
// 				return args, nil
// 			}
// 			// consume a value from remaining args list
// 			unknownFlag3 = args[0]
// 			return args[1:], nil
// 		}

// 		return args, fmt.Errorf("Unknown flag: %v", option)
// 	}

// 	// Parse args containing some unknown flags, verify that
// 	// our callback can handle all of them
// 	_, err := p.ParseArgs([]string{"--flag1=stuff", "--unknownFlag1", "blah", "--unknownFlag2", "--unknownFlag3=baz", "--flag2=foo"})

// 	if err != nil {
// 		assertErrorf(t, "Parser returned unexpected error %v", err)
// 	}

// 	assertString(t, opts.Flag1, "stuff")
// 	assertString(t, opts.Flag2, "foo")
// 	assertString(t, unknownFlag1, "blah")
// 	assertString(t, unknownFlag3, "baz")

// 	if !unknownFlag2 {
// 		assertErrorf(t, "Flag should have been set by unknown handler, but had value: %v", unknownFlag2)
// 	}

// 	// Parse args with unknown flags that callback doesn't handle, verify it returns error
// 	_, err = p.ParseArgs([]string{"--flag1=stuff", "--unknownFlagX", "blah", "--flag2=foo"})

// 	if err == nil {
// 		assertErrorf(t, "Parser should have returned error, but returned nil")
// 	}
// }

func TestChoices(t *testing.T) {
	type opts struct {
		Choice string `long:"choose" choice:"v1" choice:"v2"`
	}

	assertParseFail[opts](t, flags.ErrInvalidChoice, "Invalid value `invalid' for option `"+defaultLongOptDelimiter+"choose'. Allowed values are: v1 or v2", "--choose", "invalid")
	opt := assertParseSuccess[opts](t, "--choose", "v2")
	assertString(t, opt.Choice, "v2")
}

func TestEmbedded(t *testing.T) {
	type embedded struct {
		V bool `short:"v"`
	}
	type opts struct {
		embedded
	}

	opt := assertParseSuccess[opts](t, "-v")

	if !opt.V {
		t.Errorf("Expected V to be true")
	}
}

type command struct {
}

func (c *command) Execute(args []string) error {
	return nil
}

// func TestCommandHandlerNoCommand(t *testing.T) {
// 	var opts = struct {
// 		Value bool `short:"v"`
// 	}{}

// 	parser := NewParser(&opts, flags.Default&^flags.PrintErrors)

// 	var executedCommand flags.Commander
// 	var executedArgs []string

// 	executed := false

// 	parser.CommandHandler = func(command Commander, args []string) error {
// 		executed = true

// 		executedCommand = command
// 		executedArgs = args

// 		return nil
// 	}

// 	_, err := parser.ParseArgs([]string{"arg1", "arg2"})

// 	if err != nil {
// 		t.Fatalf("Unexpected parse error: %s", err)
// 	}

// 	if !executed {
// 		t.Errorf("Expected command handler to be executed")
// 	}

// 	if executedCommand != nil {
// 		t.Errorf("Did not exect an executed command")
// 	}

// 	assertStringArray(t, executedArgs, []string{"arg1", "arg2"})
// }

// func TestCommandHandler(t *testing.T) {
// 	var opts = struct {
// 		Value bool `short:"v"`

// 		Command command `command:"cmd"`
// 	}{}

// 	parser := NewParser(&opts, Default&^PrintErrors)

// 	var executedCommand Commander
// 	var executedArgs []string

// 	executed := false

// 	parser.CommandHandler = func(command Commander, args []string) error {
// 		executed = true

// 		executedCommand = command
// 		executedArgs = args

// 		return nil
// 	}

// 	_, err := parser.ParseArgs([]string{"cmd", "arg1", "arg2"})

// 	if err != nil {
// 		t.Fatalf("Unexpected parse error: %s", err)
// 	}

// 	if !executed {
// 		t.Errorf("Expected command handler to be executed")
// 	}

// 	if executedCommand == nil {
// 		t.Errorf("Expected command handler to be executed")
// 	}

// 	assertStringArray(t, executedArgs, []string{"arg1", "arg2"})
// }
