// Copyright 2020-2021 by Red Hat, Inc. All rights reserved.
// Use of this source is goverend by the Apache License
// that can be found in the LICENSE file.

package weldr

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/BurntSushi/toml"
)

// ListComposes returns details about the composes on the server
func (c Client) ListComposes() ([]ComposeStatusV0, []APIErrorMsg, error) {
	var errors []APIErrorMsg
	j, resp, err := c.GetRaw("GET", "/compose/queue")
	if err != nil {
		return nil, nil, err
	}
	if resp != nil {
		errors = append(errors, resp.Errors...)
		return nil, errors, nil
	}

	var composes []ComposeStatusV0

	// queue returns new and run lists of ComposeStatusV0
	var queue struct {
		New []ComposeStatusV0
		Run []ComposeStatusV0
	}
	err = json.Unmarshal(j, &queue)
	if err != nil {
		errors = append(errors, APIErrorMsg{"JSONError", err.Error()})
	} else {
		composes = append(composes, queue.New...)
		composes = append(composes, queue.Run...)
	}

	j, resp, err = c.GetRaw("GET", "/compose/finished")
	if err != nil {
		return nil, nil, err
	}
	if resp != nil {
		errors = append(errors, resp.Errors...)
		return nil, errors, nil
	}

	// finished returns finished list
	var finished struct {
		Finished []ComposeStatusV0
	}
	err = json.Unmarshal(j, &finished)
	if err != nil {
		errors = append(errors, APIErrorMsg{"JSONError", err.Error()})
	} else {
		composes = append(composes, finished.Finished...)
	}

	j, resp, err = c.GetRaw("GET", "/compose/failed")
	if err != nil {
		return nil, nil, err
	}
	if resp != nil {
		errors = append(errors, resp.Errors...)
		return nil, errors, nil
	}

	// failed returns failed list
	var failed struct {
		Failed []ComposeStatusV0
	}
	err = json.Unmarshal(j, &failed)
	if err != nil {
		errors = append(errors, APIErrorMsg{"JSONError", err.Error()})
	} else {
		composes = append(composes, failed.Failed...)
	}
	if len(errors) > 0 {
		return nil, errors, nil
	}

	return composes, nil, nil
}

// GetComposeTypes returns a list of the compose types
func (c Client) GetComposeTypes(distro string) ([]string, *APIResponse, error) {
	var route string
	if len(distro) > 0 {
		route = fmt.Sprintf("/compose/types?distro=%s", distro)
	} else {
		route = "/compose/types"
	}

	j, resp, err := c.GetRaw("GET", route)
	if err != nil {
		return nil, nil, err
	}
	if resp != nil {
		return nil, resp, nil
	}

	var types struct {
		Types []ComposeTypesV0
	}
	err = json.Unmarshal(j, &types)
	if err != nil {
		return nil, nil, err
	}

	var enabled []string
	for i := range types.Types {
		if types.Types[i].Enabled {
			enabled = append(enabled, types.Types[i].Name)
		}
	}

	return enabled, nil, nil
}

// StartCompose will start a compose of a blueprint
// Returns the UUID of the build that was started
func (c Client) StartCompose(blueprint, composeType string, size uint) (string, *APIResponse, error) {
	return c.StartComposeTest(blueprint, composeType, size, 0)
}

// StartComposeTest will start a compose of a blueprint, optionally starting a test compose
// test = 1 creates a fake failed compose
// test = 2 creates a fake successful compose
func (c Client) StartComposeTest(blueprint, composeType string, size uint, test uint) (string, *APIResponse, error) {
	var settings struct {
		Name   string `json:"blueprint_name"`
		Type   string `json:"compose_type"`
		Branch string `json:"branch"`
		Size   uint   `json:"size"`
	}
	settings.Name = blueprint
	settings.Type = composeType
	settings.Branch = "master"
	settings.Size = size
	return c.startComposeTest(settings, test)
}

// StartComposeUpload will start a compose of a blueprint and upload it to a provider
// Returns the UUID of the build that was started
func (c Client) StartComposeUpload(blueprint, composeType, imageName, profileFile string, size uint) (string, *APIResponse, error) {
	return c.StartComposeTestUpload(blueprint, composeType, imageName, profileFile, size, 0)
}

// StartComposeTestUpload will start a compose of a blueprint, optionally starting a test compose
// it will also upload the image to a provider.
// test = 1 creates a fake failed compose
// test = 2 creates a fake successful compose
func (c Client) StartComposeTestUpload(blueprint, composeType, imageName, profileFile string, size uint, test uint) (string, *APIResponse, error) {
	var settings struct {
		Name   string `json:"blueprint_name"`
		Type   string `json:"compose_type"`
		Branch string `json:"branch"`
		Size   uint   `json:"size"`
		Upload struct {
			Provider  string      `json:"provider" toml:"provider"`
			ImageName string      `json:"image_name" toml:"image_name"`
			Settings  interface{} `json:"settings" toml:"settings"`
		} `json:"upload"`
	}
	settings.Name = blueprint
	settings.Type = composeType
	settings.Branch = "master"
	settings.Size = size

	// Read the profile toml file into settings.Upload
	_, err := toml.DecodeFile(profileFile, &settings.Upload)
	if err != nil {
		return "", nil, err
	}
	settings.Upload.ImageName = imageName

	return c.startComposeTest(settings, test)
}

// StartOSTreeCompose will start a compose of a blueprint
// Returns the UUID of the build that was started
func (c Client) StartOSTreeCompose(blueprint, composeType, ref, parent, url string, size uint) (string, *APIResponse, error) {
	return c.StartOSTreeComposeTest(blueprint, composeType, ref, parent, url, size, 0)
}

// StartOSTreeComposeTest will start a compose of a blueprint, optionally starting a test compose
// test = 1 creates a fake failed compose
// test = 2 creates a fake successful compose
func (c Client) StartOSTreeComposeTest(blueprint, composeType, ref, parent, url string, size uint, test uint) (string, *APIResponse, error) {
	var settings struct {
		Name   string `json:"blueprint_name"`
		Type   string `json:"compose_type"`
		Branch string `json:"branch"`
		Size   uint   `json:"size"`
		OSTree struct {
			Ref    string `json:"ref"`
			Parent string `json:"parent"`
			URL    string `json:"url"`
		} `json:"ostree"`
	}
	settings.Name = blueprint
	settings.Type = composeType
	settings.Branch = "master"
	settings.Size = size
	settings.OSTree.Ref = ref
	settings.OSTree.Parent = parent
	settings.OSTree.URL = url

	return c.startComposeTest(settings, test)
}

// StartOSTreeComposeUpload will start a compose of a blueprint and upload it to a provider
// Returns the UUID of the build that was started
func (c Client) StartOSTreeComposeUpload(blueprint, composeType, imageName, profileFile, ref, parent, url string, size uint) (string, *APIResponse, error) {
	return c.StartOSTreeComposeTestUpload(blueprint, composeType, imageName, profileFile, ref, parent, url, size, 0)
}

// StartOSTreeComposeTestUpload will start a compose of a blueprint, optionally starting a test compose
// test = 1 creates a fake failed compose
// test = 2 creates a fake successful compose
func (c Client) StartOSTreeComposeTestUpload(blueprint, composeType, imageName, profileFile, ref, parent, url string, size uint, test uint) (string, *APIResponse, error) {
	var settings struct {
		Name   string `json:"blueprint_name"`
		Type   string `json:"compose_type"`
		Branch string `json:"branch"`
		Size   uint   `json:"size"`
		OSTree struct {
			Ref    string `json:"ref"`
			Parent string `json:"parent"`
			URL    string `json:"url"`
		} `json:"ostree"`
		Upload struct {
			Provider  string      `json:"provider" toml:"provider"`
			ImageName string      `json:"image_name" toml:"image_name"`
			Settings  interface{} `json:"settings" toml:"settings"`
		} `json:"upload"`
	}
	settings.Name = blueprint
	settings.Type = composeType
	settings.Branch = "master"
	settings.Size = size
	settings.OSTree.Ref = ref
	settings.OSTree.Parent = parent
	settings.OSTree.URL = url

	// Read the profile toml file into settings.Upload
	_, err := toml.DecodeFile(profileFile, &settings.Upload)
	if err != nil {
		return "", nil, err
	}
	settings.Upload.ImageName = imageName

	return c.startComposeTest(settings, test)
}

// startComposeTest is the common function for starting composes
// It passes through the request to the server, only setting the test flag if it is > 0
// And returns the server response to the caller
func (c Client) startComposeTest(request interface{}, test uint) (string, *APIResponse, error) {
	data, err := json.Marshal(request)
	if err != nil {
		return "", nil, err
	}

	var route string
	if test > 0 {
		route = fmt.Sprintf("/compose?test=%d", test)
	} else {
		route = "/compose"
	}
	body, resp, err := c.PostJSON(route, string(data))
	if resp != nil || err != nil {
		return "", resp, err
	}
	var build ComposeStartV0
	err = json.Unmarshal(body, &build)
	if err != nil {
		return "", nil, err
	}

	return build.ID, resp, err
}

// DeleteComposes removes a list of composes from the server
func (c Client) DeleteComposes(ids []string) ([]ComposeDeleteV0, []APIErrorMsg, error) {
	var errors []APIErrorMsg
	route := fmt.Sprintf("/compose/delete/%s", strings.Join(ids, ","))
	j, resp, err := c.DeleteRaw(route)
	if err != nil {
		return nil, nil, err
	}
	if resp != nil {
		errors = append(errors, resp.Errors...)
		return nil, errors, nil
	}

	// delete returns the status of each build id it was asked to delete
	var r struct {
		UUIDs  []ComposeDeleteV0
		Errors []APIErrorMsg
	}
	err = json.Unmarshal(j, &r)
	if err != nil {
		errors = append(errors, APIErrorMsg{"JSONError", err.Error()})
	}
	if len(errors) > 0 {
		return nil, errors, nil
	}
	if len(r.Errors) > 0 {
		errors = append(errors, r.Errors...)
	}
	return r.UUIDs, errors, nil
}

// CancelCompose cancels a compose that is waiting or running on the server
func (c Client) CancelCompose(id string) (ComposeCancelV0, []APIErrorMsg, error) {
	var r ComposeCancelV0
	var errors []APIErrorMsg
	route := fmt.Sprintf("/compose/cancel/%s", id)
	j, resp, err := c.DeleteRaw(route)
	if err != nil {
		return r, nil, err
	}
	if resp != nil {
		errors = append(errors, resp.Errors...)
		return r, errors, nil
	}

	// cancel returns the status of the single build id it was asked to cancel
	err = json.Unmarshal(j, &r)
	if err != nil {
		errors = append(errors, APIErrorMsg{"JSONError", err.Error()})
	}
	return r, errors, nil
}

// ComposeLog returns the last 1k of logs from a running compose
func (c Client) ComposeLog(id string, size int) (string, *APIResponse, error) {
	route := fmt.Sprintf("/compose/log/%s?size=%d", id, size)
	body, resp, err := c.GetRaw("GET", route)
	if err != nil {
		return "", resp, err
	}
	if resp != nil {
		return "", resp, nil
	}
	return string(body), nil, nil
}

// ComposeLogs returns the tar file of logs from the selected compose
// It returns a temporary filename, the name of the file from the server, and the type.
// The caller must clean up the temporary file when finished
func (c Client) ComposeLogs(id string) (tempFile, fileName, cType string, apiResponse *APIResponse, err error) {
	route := fmt.Sprintf("/compose/logs/%s", id)
	tempFile, cDisposition, cType, apiResponse, err := c.GetFile(route)

	if err != nil || apiResponse != nil {
		return
	}
	fileName, err = GetContentFilename(cDisposition)
	return
}

// ComposeMetadata returns the tar file of compose's metadata
// It returns a temporary filename, the name of the file from the server, and the type.
// The caller must clean up the temporary file when finished
func (c Client) ComposeMetadata(id string) (tempFile, fileName, cType string, apiResponse *APIResponse, err error) {
	route := fmt.Sprintf("/compose/metadata/%s", id)
	tempFile, cDisposition, cType, apiResponse, err := c.GetFile(route)

	if err != nil || apiResponse != nil {
		return
	}
	fileName, err = GetContentFilename(cDisposition)
	return
}

// ComposeResults returns the tar file of compose's results
// It returns a temporary filename, the name of the file from the server, and the type.
// The caller must clean up the temporary file when finished
func (c Client) ComposeResults(id string) (tempFile, fileName, cType string, apiResponse *APIResponse, err error) {
	route := fmt.Sprintf("/compose/results/%s", id)
	tempFile, cDisposition, cType, apiResponse, err := c.GetFile(route)

	if err != nil || apiResponse != nil {
		return
	}
	fileName, err = GetContentFilename(cDisposition)
	return
}

// ComposeImage returns the tar file of compose's image
// It returns a temporary filename, the name of the file from the server, and the type.
// The caller must clean up the temporary file when finished
func (c Client) ComposeImage(id string) (tempFile, fileName, cType string, apiResponse *APIResponse, err error) {
	route := fmt.Sprintf("/compose/image/%s", id)
	tempFile, cDisposition, cType, apiResponse, err := c.GetFile(route)

	if err != nil || apiResponse != nil {
		return
	}
	fileName, err = GetContentFilename(cDisposition)
	return
}

// ComposeInfo returns details about a specific compose
func (c Client) ComposeInfo(id string) (info ComposeInfoV0, resp *APIResponse, err error) {

	route := fmt.Sprintf("/compose/info/%s", id)
	j, resp, err := c.GetRaw("GET", route)
	if err != nil {
		return info, nil, err
	}
	if resp != nil {
		return info, resp, nil
	}

	err = json.Unmarshal(j, &info)
	if err != nil {
		resp = &APIResponse{false, []APIErrorMsg{{"JSONError", err.Error()}}}
	}
	return info, resp, nil
}
