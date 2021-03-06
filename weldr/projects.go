// Copyright 2021 by Red Hat, Inc. All rights reserved.
// Use of this source is goverend by the Apache License
// that can be found in the LICENSE file.

package weldr

import (
	"encoding/json"
	"fmt"
	"strings"
)

// ListProjects returns a list of all of the projects available
func (c Client) ListProjects() ([]ProjectV0, *APIResponse, error) {
	body, resp, err := c.GetJSONAll("/projects/list")
	if resp != nil || err != nil {
		return nil, resp, err
	}
	var list ProjectsListV0
	err = json.Unmarshal(body, &list)
	if err != nil {
		return nil, nil, err
	}
	return list.Projects, nil, nil
}

// ProjectsInfo returns a list of detailed info about the projects
func (c Client) ProjectsInfo(projs []string) ([]ProjectV0, *APIResponse, error) {
	route := fmt.Sprintf("/projects/info/%s", strings.Join(projs, ","))
	j, resp, err := c.GetRaw("GET", route)
	if err != nil {
		return nil, nil, err
	}
	if resp != nil {
		return nil, resp, nil
	}

	var r struct {
		Projects []ProjectV0
	}
	err = json.Unmarshal(j, &r)
	if err != nil {
		resp = &APIResponse{Status: false, Errors: []APIErrorMsg{{"JSONError", err.Error()}}}
	}
	return r.Projects, resp, nil
}
