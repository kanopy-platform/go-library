package cli

import (
	"github.com/kelseyhightower/envconfig"
	"github.com/spf13/pflag"
)

// EnvconfigProcessWithPflags can be used to run envconfig.Process without stomping on flags
// passed to the command line, since explicit flags should take precedence over the environment.
func EnvconfigProcessWithPflags(prefix string, flags *pflag.FlagSet, obj any) error {
	// discover all non-default flags before processing envconfig
	var changedFlags = map[string]string{}

	flags.VisitAll(func(f *pflag.Flag) {
		// only include non-default flags
		if f.Changed {
			changedFlags[f.Name] = f.Value.String()
		}
	})

	// apply envconfig values
	if err := envconfig.Process(prefix, obj); err != nil {
		return err
	}

	// re-apply changed flags to override environment
	for name, value := range changedFlags {
		if err := flags.Lookup(name).Value.Set(value); err != nil {
			return err
		}
	}

	return nil
}
