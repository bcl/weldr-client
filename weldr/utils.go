// Copyright 2020 by Red Hat, Inc. All rights reserved.
// Use of this source is goverend by the Apache License
// that can be found in the LICENSE file.

// nolint: deadcode,unused // These functions are only used by the *_test.go code

package weldr

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"path"
	"strconv"
)

// TestState holds the state of the mocked testing client and information about the
// environment for the tests to use.
type TestState struct {
	client     *Client
	apiVersion int
	repoDir    string
	unitTest   bool
	distros    []string
}

func setUpTestState(socketPath string, unitTest bool) (*TestState, error) {
	client := InitClientUnixSocket(context.Background(), 1, socketPath)
	state := TestState{unitTest: unitTest, client: &client}

	// Make sure the server is running
	status, resp, err := state.client.ServerStatus()
	if err != nil {
		return nil, fmt.Errorf("status request failed with client error: %s", err)
	}
	if resp != nil {
		return nil, fmt.Errorf("status request failed: %v", resp)
	}
	apiVersion, e := strconv.Atoi(status.API)
	if e != nil {
		state.apiVersion = 0
	} else {
		state.apiVersion = apiVersion
	}
	fmt.Printf("Running tests against %s %s server using V%d API\n\n", status.Backend, status.Build, state.apiVersion)
	return &state, nil
}

// SetUpTemporaryRepository creates a temporary repository
func SetUpTemporaryRepository() (string, error) {
	dir, err := ioutil.TempDir("/tmp", "osbuild-composer-test-")
	if err != nil {
		return "", err
	}
	cmd := exec.Command("createrepo_c", path.Join(dir))
	err = cmd.Start()
	if err != nil {
		return "", err
	}
	err = cmd.Wait()
	if err != nil {
		return "", err
	}
	return dir, nil
}

// TearDownTemporaryRepository removes the temporary repository
func TearDownTemporaryRepository(dir string) error {
	return os.RemoveAll(dir)
}

// MockClient implements the HTTPClient interface for testing client requests
// Set DoFunc to a function that returns whatever response is required
type MockClient struct {
	DoFunc func(req *http.Request) (*http.Response, error)
	Req    http.Request
}

// Do saves the request in m.Req and runs the function set in m.DoFunc
// instead of making an actual network query
func (m *MockClient) Do(req *http.Request) (*http.Response, error) {
	m.Req = *req
	return m.DoFunc(req)
}
