// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package build

import (
	"cto-github.cisco.com/NFV-BU/go-msx/cli"
	"github.com/spf13/cobra"
)

type Target struct {
	Name        string
	Description string
	Fn          cli.CommandFunc
}

func AddTarget(name, description string, fn cli.CommandFunc) *cobra.Command {
	wrapper := func(args []string) error {
		logger.Infof("Executing target '%s': %s", name, description)
		err := fn(args)
		if err != nil {
			logger.Infof("Target failed: '%s': %s", name, err.Error())
		} else {
			logger.Infof("Target succeeded: '%s'", name)
		}
		return err
	}

	if cmd, err := cli.AddCommand(name, description, wrapper); err != nil {
		panic(err.Error())
	} else {
		cmd.PreRunE = loadConfig
		cmd.FParseErrWhitelist = cli.RootCmd().FParseErrWhitelist
		return cmd
	}
}
