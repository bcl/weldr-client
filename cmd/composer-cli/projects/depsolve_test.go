// Copyright 2021 by Red Hat, Inc. All rights reserved.
// Use of this source is goverend by the Apache License
// that can be found in the LICENSE file.

package projects

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/osbuild/weldr-client/cmd/composer-cli/root"
)

func TestCmdProjectsDepsolve(t *testing.T) {
	// Test the "projects depsolve" command
	json := `{
		"projects": [
            {
                "arch": "x86_64",
                "check_gpg": true,
                "checksum": "sha256:92c1615d385b32088f78a6574a2bf89a6bb29d9858abdd71471ef5113ef0831f",
                "epoch": 0,
                "name": "acl",
                "release": "9.fc33",
                "remote_location": "https://mirrors.rit.edu/fedora/fedora/linux/releases/33/Everything/x86_64/os/Packages/a/acl-2.2.53-9.fc33.x86_64.rpm",
                "version": "2.2.53"
            },
            {
                "arch": "noarch",
                "check_gpg": true,
                "checksum": "sha256:f4efaa5bc8382246d8230ece8bacebd3c29eb9fd52b509b1e6575e643953851b",
                "epoch": 0,
                "name": "basesystem",
                "release": "10.fc33",
                "remote_location": "https://mirrors.rit.edu/fedora/fedora/linux/releases/33/Everything/x86_64/os/Packages/b/basesystem-11-10.fc33.noarch.rpm",
                "version": "11"
            }
	]}`
	mc := root.SetupCmdTest(func(request *http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: 200,
			Body:       ioutil.NopCloser(bytes.NewReader([]byte(json))),
		}, nil
	})

	cmd, out, err := root.ExecuteTest("projects", "depsolve", "bash")
	defer out.Close()
	require.Nil(t, err)
	require.NotNil(t, out.Stdout)
	require.NotNil(t, out.Stderr)
	require.NotNil(t, cmd)
	assert.Equal(t, cmd, depsolveCmd)
	stdout, err := ioutil.ReadAll(out.Stdout)
	assert.Nil(t, err)
	assert.Contains(t, string(stdout), "acl-2.2.53-9.fc33.x86_64")
	assert.Contains(t, string(stdout), "basesystem-11-10.fc33.noarch")
	stderr, err := ioutil.ReadAll(out.Stderr)
	assert.Nil(t, err)
	assert.Equal(t, []byte(""), stderr)
	assert.Equal(t, "GET", mc.Req.Method)
}

func TestCmdProjectsDepsolveDistro(t *testing.T) {
	// Test the "projects depsolve --distro=test-distro" command
	json := `{
		"projects": [
            {
                "arch": "x86_64",
                "check_gpg": true,
                "checksum": "sha256:92c1615d385b32088f78a6574a2bf89a6bb29d9858abdd71471ef5113ef0831f",
                "epoch": 0,
                "name": "acl",
                "release": "9.fc33",
                "remote_location": "https://mirrors.rit.edu/fedora/fedora/linux/releases/33/Everything/x86_64/os/Packages/a/acl-2.2.53-9.fc33.x86_64.rpm",
                "version": "2.2.53"
            },
            {
                "arch": "noarch",
                "check_gpg": true,
                "checksum": "sha256:f4efaa5bc8382246d8230ece8bacebd3c29eb9fd52b509b1e6575e643953851b",
                "epoch": 0,
                "name": "basesystem",
                "release": "10.fc33",
                "remote_location": "https://mirrors.rit.edu/fedora/fedora/linux/releases/33/Everything/x86_64/os/Packages/b/basesystem-11-10.fc33.noarch.rpm",
                "version": "11"
            }
	]}`
	mc := root.SetupCmdTest(func(request *http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: 200,
			Body:       ioutil.NopCloser(bytes.NewReader([]byte(json))),
		}, nil
	})

	cmd, out, err := root.ExecuteTest("projects", "depsolve", "--distro=test-distro", "bash")
	defer out.Close()
	require.Nil(t, err)
	require.NotNil(t, out.Stdout)
	require.NotNil(t, out.Stderr)
	require.NotNil(t, cmd)
	assert.Equal(t, cmd, depsolveCmd)
	stdout, err := ioutil.ReadAll(out.Stdout)
	assert.Nil(t, err)
	assert.Contains(t, string(stdout), "acl-2.2.53-9.fc33.x86_64")
	assert.Contains(t, string(stdout), "basesystem-11-10.fc33.noarch")
	stderr, err := ioutil.ReadAll(out.Stderr)
	assert.Nil(t, err)
	assert.Equal(t, []byte(""), stderr)
	assert.Equal(t, "GET", mc.Req.Method)
}

func TestCmdProjectsDepsolveBadDistro(t *testing.T) {
	// Test the "projects depsolve --distro=homer" command
	mc := root.SetupCmdTest(func(request *http.Request) (*http.Response, error) {
		json := `{
        "errors": [
            {
                "id": "DistroError",
                "msg": "Invalid distro: homer"
            }
        ],
        "status": false
}`

		return &http.Response{
			Request:    request,
			StatusCode: 400,
			Body:       ioutil.NopCloser(bytes.NewReader([]byte(json))),
		}, nil
	})

	// Get the compose types
	cmd, out, err := root.ExecuteTest("projects", "depsolve", "--distro=homer", "bash")
	defer out.Close()
	require.NotNil(t, err)
	require.NotNil(t, out.Stdout)
	require.NotNil(t, out.Stderr)
	require.NotNil(t, cmd)
	assert.Equal(t, cmd, depsolveCmd)
	stdout, err := ioutil.ReadAll(out.Stdout)
	assert.Nil(t, err)
	assert.Equal(t, []byte(""), stdout)
	stderr, err := ioutil.ReadAll(out.Stderr)
	assert.Nil(t, err)
	assert.Contains(t, string(stderr), "DistroError: Invalid distro: homer")
	assert.Equal(t, "GET", mc.Req.Method)
}

func TestCmdProjectsDepsolveUnknown(t *testing.T) {
	// Test the "projects depsolve" command with an unknown package
	json := `{
        "errors": [
            {
                "id": "ProjectsError",
                "msg": "BadRequest: DNF error occured: MarkingErrors: Error occurred when marking packages for installation: Problems in request:\nmissing packages: homer"
            }
        ],
        "status": false
	}`
	mc := root.SetupCmdTest(func(request *http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: 200,
			Body:       ioutil.NopCloser(bytes.NewReader([]byte(json))),
		}, nil
	})

	cmd, out, err := root.ExecuteTest("projects", "depsolve", "bash,homer")
	defer out.Close()
	require.NotNil(t, err)
	assert.Equal(t, root.ExecutionError(cmd, ""), err)
	require.NotNil(t, out.Stdout)
	require.NotNil(t, out.Stderr)
	require.NotNil(t, cmd)
	assert.Equal(t, cmd, depsolveCmd)
	stdout, err := ioutil.ReadAll(out.Stdout)
	assert.Nil(t, err)
	assert.Equal(t, []byte(""), stdout)
	stderr, err := ioutil.ReadAll(out.Stderr)
	assert.Nil(t, err)
	assert.Contains(t, string(stderr), "missing packages: homer")
	assert.Equal(t, "GET", mc.Req.Method)
}
