// Copyright 2021 by Red Hat, Inc. All rights reserved.
// Use of this source is goverend by the Apache License
// that can be found in the LICENSE file.

package weldr

import (
	"encoding/json"
	"sort"
)

// ListDistros returns a list of all of the available distributions
func (c Client) ListDistros() ([]string, *APIResponse, error) {
	j, resp, err := c.GetRaw("GET", "/distros/list")
	if err != nil {
		return nil, nil, err
	}
	if resp != nil {
		return nil, resp, nil
	}

	var distros struct {
		Distros []string
	}
	err = json.Unmarshal(j, &distros)
	if err != nil {
		return nil, nil, err
	}
	sort.Strings(distros.Distros)
	return distros.Distros, nil, nil
}
