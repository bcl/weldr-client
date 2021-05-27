// Copyright 2021 by Red Hat, Inc. All rights reserved.
// Use of this source is goverend by the Apache License
// that can be found in the LICENSE file.

package weldr

import (
	"encoding/json"
	"fmt"
	"strings"
)

// ListModules returns a list of all of the modules available
// NOTE: These are just packages, the server does not support modules directly
func (c Client) ListModules(distro string) ([]ModuleV0, *APIResponse, error) {
	var route string
	if len(distro) > 0 {
		route = fmt.Sprintf("/modules/list?distro=%s", distro)
	} else {
		route = "/modules/list"
	}

	body, resp, err := c.GetJSONAll(route)
	if resp != nil || err != nil {
		return nil, resp, err
	}
	var list ModulesListV0
	err = json.Unmarshal(body, &list)
	if err != nil {
		return nil, nil, err
	}
	return list.Modules, nil, nil
}

// ModulesInfo returns a list of detailed info about the modules, including deps
func (c Client) ModulesInfo(names []string, distro string) ([]ProjectV0, *APIResponse, error) {
	route := fmt.Sprintf("/modules/info/%s", strings.Join(names, ","))

	if len(distro) > 0 {
		route = fmt.Sprintf("%s?distro=%s", route, distro)
	}

	j, resp, err := c.GetRaw("GET", route)
	if err != nil {
		return nil, nil, err
	}
	if resp != nil {
		return nil, resp, nil
	}

	var r struct {
		Modules []ProjectV0
	}
	err = json.Unmarshal(j, &r)
	if err != nil {
		resp = &APIResponse{Status: false, Errors: []APIErrorMsg{{"JSONError", err.Error()}}}
	}
	return r.Modules, resp, nil
}
