// Copyright (c) 2023 Canonical Ltd
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

package fwstate

import (
	"github.com/canonical/pebble/internals/overlord/state"
	"gopkg.in/tomb.v2"
)

type FirmwareManager struct{}

// NewManager creates a new FirmwareManager.
func NewManager(runner *state.TaskRunner) *FirmwareManager {
	manager := &FirmwareManager{}
	runner.AddHandler("refresh", manager.refreshDoHandler, nil)
	return manager
}

func (m *FirmwareManager) refreshDoHandler(task *state.Task, tomb *tomb.Tomb) error {

	var args RefreshArgs
	err := task.Get("refresh-args", &args)
	if err != nil {
		return err
	}
	// todo
}
