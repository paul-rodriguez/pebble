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
	"fmt"
	"io"

	"github.com/canonical/pebble/internals/overlord/state"
)

type RefreshArgs struct {
	source io.Reader
	path   string
	size   int
}

func NewRefreshTask(st *state.State, args *RefreshArgs) (*state.Task, error) {
	task := st.NewTask("refresh", fmt.Sprintf("Refresh file %s", args.path))
	task.Set("refresh-args", args)
	return task, nil
}
