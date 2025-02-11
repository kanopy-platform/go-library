package cli_test

import (
	"os"
	"testing"

	"github.com/kanopy-platform/go-library/cli"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

type testCLI struct {
	String           string
	Int              int
	SubString        string `split_words:"true"`
	OverrideFlagRoot string `split_words:"true"`
	OverrideFlagSub  string `split_words:"true"`
}

func (c *testCLI) RootCmd() *cobra.Command {
	emptyRun := func(_ *cobra.Command, _ []string) {}

	root := &cobra.Command{
		Use: "root",
		Run: emptyRun,
		PersistentPreRunE: func(cmd *cobra.Command, _ []string) error {
			return cli.EnvconfigProcessWithPflags("APP", cmd.Flags(), c)
		},
	}

	root.PersistentFlags().StringVar(&c.String, "string", "default", "")
	root.PersistentFlags().IntVar(&c.Int, "int", 1, "")
	root.PersistentFlags().StringVar(&c.OverrideFlagRoot, "override", "root", "")

	sub := &cobra.Command{Use: "sub", Run: emptyRun}
	sub.Flags().StringVar(&c.SubString, "substring", "default", "")
	sub.PersistentFlags().StringVar(&c.OverrideFlagSub, "override", "sub", "")

	root.AddCommand(sub)

	return root
}

// checks that the correct order of precedence is adhered to: flag > env > default
func TestPrecedence(t *testing.T) {
	tests := []struct {
		desc string
		env  map[string]string
		args []string
		want testCLI
	}{
		{
			desc: "test root defaults",
			want: testCLI{
				String:           "default",
				Int:              1,
				SubString:        "default",
				OverrideFlagRoot: "root",
				OverrideFlagSub:  "sub",
			},
		},
		{
			desc: "test root override",
			args: []string{"--override=flag"},
			want: testCLI{
				String:           "default",
				Int:              1,
				SubString:        "default",
				OverrideFlagRoot: "flag",
				OverrideFlagSub:  "sub",
			},
		},
		{
			desc: "test sub defaults",
			args: []string{"sub"},
			want: testCLI{
				String:           "default",
				Int:              1,
				SubString:        "default",
				OverrideFlagRoot: "root",
				OverrideFlagSub:  "sub",
			},
		},
		{
			desc: "test sub env",
			args: []string{"sub"},
			env: map[string]string{
				"APP_STRING":            "env",
				"APP_INT":               "2",
				"APP_SUB_STRING":        "env",
				"APP_OVERRIDE_FLAG_SUB": "env",
			},
			want: testCLI{
				String:           "env",
				Int:              2,
				SubString:        "env",
				OverrideFlagRoot: "root",
				OverrideFlagSub:  "env",
			},
		},
		{
			desc: "test sub flags",
			args: []string{
				"sub",
				"--string=flag",
				"--int=3",
				"--substring=flag",
				"--override=flag",
			},
			want: testCLI{
				String:           "flag",
				Int:              3,
				SubString:        "flag",
				OverrideFlagRoot: "root",
				OverrideFlagSub:  "flag",
			},
		},
		{
			desc: "test sub flags > env",
			env: map[string]string{
				"APP_STRING":             "env",
				"APP_INT":                "2",
				"APP_SUB_STRING":         "env",
				"APP_OVERRIDE_FLAG_ROOT": "env",
				"APP_OVERRIDE_FLAG_SUB":  "env",
			},
			args: []string{
				"sub",
				"--string=flag",
				"--int=3",
				"--substring=flag",
				"--override=flag",
			},
			want: testCLI{
				String:           "flag",
				Int:              3,
				SubString:        "flag",
				OverrideFlagRoot: "env",
				OverrideFlagSub:  "flag",
			},
		},
	}

	for _, test := range tests {
		for k, v := range test.env {
			assert.NoError(t, os.Setenv(k, v))
		}

		cli := &testCLI{}
		cmd := cli.RootCmd()
		cmd.SetArgs(test.args)
		assert.NoError(t, cmd.Execute())

		assert.Equal(t, test.want, *cli, test.desc)

		for k := range test.env {
			assert.NoError(t, os.Unsetenv(k))
		}
	}
}
