// Copyright (c) 2021 Canonical Ltd
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

package daemon

import (
	"encoding/json"
	"fmt"
	"io"
	"mime"
	"mime/multipart"
	"net/http"
	pathpkg "path"

	"github.com/canonical/pebble/internals/osutil"
	"github.com/canonical/pebble/internals/osutil/sys"
)

func v1DeviceInstall(_ *Command, req *http.Request, _ *UserState) Response {
	contentType := req.Header.Get("Content-Type")
	mediaType, params, err := mime.ParseMediaType(contentType)
	if err != nil {
		return statusBadRequest("invalid Content-Type %q", contentType)
	}

	switch mediaType {
	case "multipart/form-data":
		boundary := params["boundary"]
		if len(boundary) < minBoundaryLength {
			return statusBadRequest("invalid boundary %q", boundary)
		}
		return installFiles(req.Body, boundary)
	default:
		return statusBadRequest("invalid media type %q", mediaType)
	}
}

// Writing files

type installFilesItem struct {
	Path        string `json:"path"`
	MakeDirs    bool   `json:"make-dirs"`
	Permissions string `json:"permissions"`
	UserID      *int   `json:"user-id"`
	User        string `json:"user"`
	GroupID     *int   `json:"group-id"`
	Group       string `json:"group"`
}

func installFiles(body io.Reader, boundary string) Response {
	// Read metadata part (field name "request").
	mr := multipart.NewReader(body, boundary)
	part, err := mr.NextPart()
	if err != nil {
		return statusBadRequest("cannot read request metadata: %v", err)
	}
	if part.FormName() != "request" {
		return statusBadRequest(`metadata field name must be "request", got %q`, part.FormName())
	}

	// Decode metadata about files to write.
	var payload struct {
		Action string             `json:"action"`
		Files  []installFilesItem `json:"files"`
	}
	decoder := json.NewDecoder(part)
	if err := decoder.Decode(&payload); err != nil {
		return statusBadRequest("cannot decode request metadata: %v", err)
	}
	if payload.Action != "write" {
		return statusBadRequest(`multipart action must be "write", got %q`, payload.Action)
	}
	if len(payload.Files) == 0 {
		return statusBadRequest("must specify one or more files")
	}
	infos := make(map[string]installFilesItem)
	for _, file := range payload.Files {
		infos[file.Path] = file
	}

	errors := make(map[string]error)
	for i := 0; ; i++ {
		part, err = mr.NextPart()
		if err == io.EOF {
			break
		}
		if err != nil {
			return statusBadRequest("cannot read file part %d: %v", i, err)
		}
		if part.FormName() != "files" {
			return statusBadRequest(`field name must be "files", got %q`, part.FormName())
		}
		path := multipartFilename(part)
		info, ok := infos[path]
		if !ok {
			return statusBadRequest("no metadata for path %q", path)
		}
		errors[path] = installFile(info, part)
		part.Close()
	}

	// Build list of results with any errors.
	result := make([]fileResult, len(payload.Files))
	for i, file := range payload.Files {
		err, ok := errors[file.Path]
		if !ok {
			// Ensure we wrote all the files in the metadata.
			err = fmt.Errorf("no file content for path %q", file.Path)
		}
		result[i] = fileResult{
			Path:  file.Path,
			Error: fileErrorToResult(err),
		}
	}
	return SyncResponse(result)
}

func installFile(item installFilesItem, source io.Reader) error {
	if !pathpkg.IsAbs(item.Path) {
		return nonAbsolutePathError(item.Path)
	}

	uid, gid, err := normalizeUidGid(item.UserID, item.GroupID, item.User, item.Group)
	if err != nil {
		return fmt.Errorf("cannot look up user and group: %w", err)
	}

	// Create parent directory if needed.
	if item.MakeDirs {
		err := mkdirAllUserGroup(pathpkg.Dir(item.Path), 0o755, uid, gid)
		if err != nil {
			return fmt.Errorf("cannot create directory: %w", err)
		}
	}

	// Atomically write file content to destination.
	perm, err := parsePermissions(item.Permissions, 0o644)
	if err != nil {
		return err
	}
	sysUid, sysGid := sys.UserID(osutil.NoChown), sys.GroupID(osutil.NoChown)
	if uid != nil && gid != nil {
		sysUid, sysGid = sys.UserID(*uid), sys.GroupID(*gid)
	}
	return atomicWriteChown(item.Path, source, perm, osutil.AtomicWriteChmod, sysUid, sysGid)
}
