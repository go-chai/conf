package conf_test

import (
	"testing"
	"time"

	"github.com/go-chai/conf"
	"github.com/stretchr/testify/require"
)

var onlyDefaults = &defaultOptions{
	IntDefault:     1,
	Float64Default: -3.14,
	StringDefault:  "abc",
	TimeDefault:    1 * time.Minute,
	MapDefault: map[string]int{
		"a": 1,
	},
	Map:          map[string]int{},
	SliceDefault: []int{1, 2},
}

var nonDefaultOverrides = &defaultOptions{
	Int:            3,
	IntDefault:     1,
	Float64:        2.712,
	Float64Default: -3.14,
	NumericFlag:    false,
	String:         "asdf",
	StringDefault:  "abc",
	Time:           13 * time.Second,
	TimeDefault:    1 * time.Minute,
	Map: map[string]int{
		"val1": 3,
		"val2": 4,
	},
	MapDefault: map[string]int{
		"a": 1,
	},
	Slice:        []int{1, 2, 3},
	SliceDefault: []int{1, 2},
}

var defaultOverrides = &defaultOptions{
	Int:            3,
	IntDefault:     13,
	Float64:        2.712,
	Float64Default: 1.1234,
	NumericFlag:    false,
	String:         "asdf",
	StringDefault:  "defg",
	Time:           13 * time.Second,
	TimeDefault:    11 * time.Minute,
	Map: map[string]int{
		"val1": 3,
		"val2": 4,
	},
	MapDefault: map[string]int{
		"val21": 21,
		"val22": 22,
	},
	Slice:        []int{1, 2, 3},
	SliceDefault: []int{4, 5, 6},
}

var flagOverrides = &defaultOptions{
	Int:            4,
	IntDefault:     5,
	Float64:        2.712,
	Float64Default: 1.1234,
	NumericFlag:    false,
	String:         "asdf",
	StringDefault:  "defg",
	Time:           17 * time.Second,
	TimeDefault:    19 * time.Second,
	Map: map[string]int{
		"val1": 3,
		"val2": 4,
	},
	MapDefault: map[string]int{
		"val21": 21,
		"val22": 22,
	},
	Slice:        []int{1, 2, 3},
	SliceDefault: []int{4, 5, 6},
}

func Test_Load_ConfigFiles(t *testing.T) {
	var tcs = map[string]struct {
		opts     []conf.ConfOption
		expected any
	}{
		"no args > no config": {
			opts:     []conf.ConfOption{conf.Args([]string{})},
			expected: onlyDefaults,
		},
		"no args > paths > empty YAML": {
			opts:     []conf.ConfOption{conf.Paths("testdata/config-empty.yaml"), conf.Args([]string{})},
			expected: onlyDefaults,
		},
		"conf args > empty YAML": {
			opts:     []conf.ConfOption{conf.ConfigFlag("conf", ""), conf.Args([]string{"--conf=testdata/config-empty.yaml"})},
			expected: onlyDefaults,
		},
		"no args > default flag config path > empty YAML": {
			opts:     []conf.ConfOption{conf.ConfigFlag("conf", "testdata/config-empty.yaml"), conf.Args([]string{})},
			expected: onlyDefaults,
		},
		"no args > paths > empty YML": {
			opts:     []conf.ConfOption{conf.Paths("testdata/config-empty.yml"), conf.Args([]string{})},
			expected: onlyDefaults,
		},
		"conf args > empty YML": {
			opts:     []conf.ConfOption{conf.ConfigFlag("conf", ""), conf.Args([]string{"--conf=testdata/config-empty.yml"})},
			expected: onlyDefaults,
		},
		"no args > default flag config path > empty YML": {
			opts:     []conf.ConfOption{conf.ConfigFlag("conf", "testdata/config-empty.yml"), conf.Args([]string{})},
			expected: onlyDefaults,
		},
		"no args > paths > empty TOML": {
			opts:     []conf.ConfOption{conf.Paths("testdata/config-empty.toml"), conf.Args([]string{})},
			expected: onlyDefaults,
		},
		"conf args > empty TOML": {
			opts:     []conf.ConfOption{conf.ConfigFlag("conf", ""), conf.Args([]string{"--conf=testdata/config-empty.toml"})},
			expected: onlyDefaults,
		},
		"no args > default flag config path > empty TOML": {
			opts:     []conf.ConfOption{conf.ConfigFlag("conf", "testdata/config-empty.toml"), conf.Args([]string{})},
			expected: onlyDefaults,
		},
		"no args > paths > empty JSON": {
			opts:     []conf.ConfOption{conf.Paths("testdata/config-empty.json"), conf.Args([]string{})},
			expected: onlyDefaults,
		},
		"conf args > empty JSON": {
			opts:     []conf.ConfOption{conf.ConfigFlag("conf", ""), conf.Args([]string{"--conf=testdata/config-empty.json"})},
			expected: onlyDefaults,
		},
		"no args > default flag config path > empty JSON": {
			opts:     []conf.ConfOption{conf.ConfigFlag("conf", "testdata/config-empty.json"), conf.Args([]string{})},
			expected: onlyDefaults,
		},
		"no args > default flag config path > YAML with non-default overrides": {
			opts:     []conf.ConfOption{conf.ConfigFlag("conf", "testdata/config-default.yaml"), conf.Args([]string{})},
			expected: nonDefaultOverrides,
		},
		"no args > default flag config path > YAML with default overrides": {
			opts:     []conf.ConfOption{conf.ConfigFlag("conf", "testdata/config.yaml"), conf.Args([]string{})},
			expected: defaultOverrides,
		},
		"value args > default flag config path > YAML with default overrides": {
			opts:     []conf.ConfOption{conf.ConfigFlag("conf", "testdata/config.yaml"), conf.Args([]string{"--i=4", "--id=5", "--t=17s", "--td=19s"})},
			expected: flagOverrides,
		},
		"no args > default flag config path > TOML with non-default overrides": {
			opts:     []conf.ConfOption{conf.ConfigFlag("conf", "testdata/config-default.toml"), conf.Args([]string{})},
			expected: nonDefaultOverrides,
		},
		"no args > default flag config path > TOML with default overrides": {
			opts:     []conf.ConfOption{conf.ConfigFlag("conf", "testdata/config.toml"), conf.Args([]string{})},
			expected: defaultOverrides,
		},
		"value args > default flag config path > TOML with default overrides": {
			opts:     []conf.ConfOption{conf.ConfigFlag("conf", "testdata/config.toml"), conf.Args([]string{"--i=4", "--id=5", "--t=17s", "--td=19s"})},
			expected: flagOverrides,
		},
		"no args > paths > TOML with non-default overrides": {
			opts:     []conf.ConfOption{conf.Paths("testdata/config-default.toml"), conf.Args([]string{})},
			expected: nonDefaultOverrides,
		},
		"no args > paths > TOML with default overrides": {
			opts:     []conf.ConfOption{conf.Paths("testdata/config.toml"), conf.Args([]string{})},
			expected: defaultOverrides,
		},
		"value args > paths > TOML with default overrides": {
			opts:     []conf.ConfOption{conf.Paths("testdata/config.toml"), conf.Args([]string{"--i=4", "--id=5", "--t=17s", "--td=19s"})},
			expected: flagOverrides,
		},
		"no args > paths > JSON with non-default overrides": {
			opts:     []conf.ConfOption{conf.Paths("testdata/config-default.json"), conf.Args([]string{})},
			expected: nonDefaultOverrides,
		},
		"no args > paths > JSON with default overrides": {
			opts:     []conf.ConfOption{conf.Paths("testdata/config.json"), conf.Args([]string{})},
			expected: defaultOverrides,
		},
		"value args > paths > JSON with default overrides": {
			opts:     []conf.ConfOption{conf.Paths("testdata/config.json"), conf.Args([]string{"--i=4", "--id=5", "--t=17s", "--td=19s"})},
			expected: flagOverrides,
		},
		"value args > default flag config path > YAML with non-default overrides > TOML with default overrides": {
			opts:     []conf.ConfOption{conf.ConfigFlag("conf", "testdata/config.yaml", "testdata/config.toml"), conf.Args([]string{"--i=4", "--id=5", "--t=17s", "--td=19s"})},
			expected: flagOverrides,
		},
	}

	for name, tc := range tcs {
		t.Run(name, func(t *testing.T) {
			cfg, err := conf.Load[defaultOptions](conf.ConfOptions(tc.opts))
			require.NoError(t, err)

			require.Equal(t, tc.expected, cfg)
		})
	}
}
