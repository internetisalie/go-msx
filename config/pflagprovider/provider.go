// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package pflagprovider

import (
	"context"
	"cto-github.cisco.com/NFV-BU/go-msx/config"
	"cto-github.cisco.com/NFV-BU/go-msx/config/args"
	"cto-github.cisco.com/NFV-BU/go-msx/log"
	"fmt"
	"github.com/spf13/pflag"
	"os"
)

var logger = log.NewLogger("msx.config.pflagprovider")

type Provider struct {
	config.Named
	prefix  string
	extras  *args.CommandLineExtrasProvider
	flagSet *pflag.FlagSet
	config.SilentNotifier
}

func (p *Provider) flagExists(name string) bool {
	return p.flagSet.Lookup(name) != nil
}

func (p *Provider) Load(ctx context.Context) (entries config.ProviderEntries, err error) {
	extras, err := p.extras.Load(ctx)
	if err != nil {
		return
	}
	entries = extras.Clone()

	if !p.flagSet.Parsed() {
		if err := p.flagSet.Parse(os.Args[1:]); err != nil {
			return nil, err
		}
	}

	p.flagSet.VisitAll(func(flag *pflag.Flag) {
		key := config.PrefixWithName(p.prefix, flag.Name)
		entries = append(entries, config.NewEntry(p, key, flag.Value.String()))
	})

	return entries, nil
}

func NewProvider(name string, flagSet *pflag.FlagSet, prefix string) *Provider {
	p := &Provider{
		Named:   config.NewNamed(fmt.Sprintf("%s: pflag", name)),
		prefix:  prefix,
		flagSet: flagSet,
	}
	p.extras = args.NewCommandLineExtrasProvider(p, p.flagExists)
	return p
}
