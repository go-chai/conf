package conf_test

import (
	"testing"
	"time"

	"github.com/go-chai/conf"
	"github.com/stretchr/testify/require"
)

func Test_Load_ConfigFiles(t *testing.T) {
	var tcs = map[string]struct {
		args     []string
		confFlag string
		confFile []string
		expected any
	}{
		"no args, no default overrides": {
			args:     []string{},
			confFlag: "--conf",
			confFile: []string{"testdata/config-default.yaml"},
			expected: &defaultOptions{
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
			},
		},
		"no args, with default overrides": {
			args:     []string{},
			confFlag: "--conf",
			confFile: []string{"testdata/config.yaml"},
			expected: &defaultOptions{
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
					"a":     1,
				},
				Slice:        []int{1, 2, 3},
				SliceDefault: []int{4, 5, 6},
			},
		},
		"with args, with default overrides": {
			args:     []string{"--i=4", "--id=5", "--t=17s", "--td=19s"},
			confFlag: "--conf",
			confFile: []string{"testdata/config.yaml"},
			expected: &defaultOptions{
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
					"a":     1,
				},
				Slice:        []int{1, 2, 3},
				SliceDefault: []int{4, 5, 6},
			},
		},
	}

	for _, tc := range tcs {
		cfg, err := conf.Load[defaultOptions](conf.Args(tc.args), conf.ConfigFlag(tc.confFlag, tc.confFile...))
		require.NoError(t, err)

		require.Equal(t, tc.expected, cfg)
	}
}
