package cobraprovider

import (
	"context"
	"cto-github.cisco.com/NFV-BU/go-msx/config/args"
	"cto-github.cisco.com/NFV-BU/go-msx/config"
	"cto-github.cisco.com/NFV-BU/go-msx/log"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"os"
	"sync"
)

var logger = log.NewLogger("msx.cli.cobra")

type ConfigProvider struct {
	prefix  string
	appName string
	extras  map[string]string
	flagset *pflag.FlagSet
	once    sync.Once
}

func (f *ConfigProvider) Load(ctx context.Context) (settings map[string]string, err error) {
	logger.Info("Loading command line config")

	f.once.Do(func() {
		f.extras = args.Extras(func(name string) bool {
			return f.flagset.Lookup(name) != nil
		})
	})

	if !f.flagset.Parsed() {
		if err := f.flagset.Parse(os.Args[1:]); err != nil {
			return nil, err
		}
	}

	settings = make(map[string]string)
	f.flagset.VisitAll(func(flag *pflag.Flag) {
		key := config.NormalizeKey(f.prefix + flag.Name)
		settings[key] = flag.Value.String()
	})

	// Apply extras
	for k, v := range f.extras {
		settings[k] = v
	}

	// Apply application name
	settings["spring.application.name"] = f.appName

	return settings, nil
}

func ExtractFlagSet(command *cobra.Command) *pflag.FlagSet {
	flagSet := pflag.NewFlagSet(command.Name(), pflag.ContinueOnError)
	flagSet.ParseErrorsWhitelist.UnknownFlags = true

	command.InheritedFlags().VisitAll(func(flag *pflag.Flag) {
		flagSet.AddFlag(flag)
	})

	command.LocalFlags().VisitAll(func(flag *pflag.Flag) {
		flagSet.AddFlag(flag)
	})

	return flagSet
}

func NewCobraSource(command *cobra.Command, prefix string) *ConfigProvider {
	flagSet := ExtractFlagSet(command)
	return &ConfigProvider{
		prefix:  prefix,
		flagset: flagSet,
		appName: command.Root().Name(),
	}
}
