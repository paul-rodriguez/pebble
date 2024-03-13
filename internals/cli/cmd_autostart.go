// Copyright (c) 2014-2020 Canonical Ltd
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License version 3 as
// published by the Free Software Foundation.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

package cli

import (
	"fmt"

	"github.com/canonical/go-flags"

	"github.com/canonical/pebble/client"
)

const cmdAutoStartSummary = "Start services set to start by default"
const cmdAutoStartDescription = `
The autostart command starts the services that were configured
to start by default.
`

type cmdAutoStart struct {
	waitMixin
}

func init() {
	AddCommand(&CmdInfo{
		Name:        "autostart",
		Summary:     cmdAutoStartSummary,
		Description: cmdAutoStartDescription,
		ArgsHelp:    waitArgsHelp,
		New: func(opts *CmdOptions) flags.Commander {
			return &cmdAutoStart{}
		},
	})
}

func (cmd cmdAutoStart) Execute(args []string) error {
	if len(args) > 1 {
		return ErrExtraArgs
	}

	_, address := getEnvPaths()
	config := ConfigFromAddress(address)
	commandClient, err := client.New(config)
	if err != nil {
		return fmt.Errorf("cannot create client: %w", err)
	}

	servopts := client.ServiceOptions{}
	changeID, err := commandClient.AutoStart(&servopts)
	if err != nil {
		return err
	}

	if _, err := cmd.wait(commandClient, changeID); err != nil {
		if err == noWait {
			return nil
		}
		return err
	}
	maybePresentWarnings(commandClient.WarningsSummary())
	return nil
}
