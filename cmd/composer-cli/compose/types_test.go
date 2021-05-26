// Copyright 2020 by Red Hat, Inc. All rights reserved.
// Use of this source is goverend by the Apache License
// that can be found in the LICENSE file.

package compose

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/osbuild/weldr-client/cmd/composer-cli/root"
)

func TestCmdComposeTypes(t *testing.T) {
	// Test the "compose types" command
	mc := root.SetupCmdTest(func(request *http.Request) (*http.Response, error) {
		json := `{
    "types": [
        {
            "name": "ami",
            "enabled": true
        },
        {
            "name": "fedora-iot-commit",
            "enabled": true
        },
        {
            "name": "openstack",
            "enabled": true
        },
        {
            "name": "qcow2",
            "enabled": true
        },
        {
            "name": "vhd",
            "enabled": true
        },
        {
            "name": "vmdk",
            "enabled": true
        }
    ]
}`

		return &http.Response{
			StatusCode: 200,
			Body:       ioutil.NopCloser(bytes.NewReader([]byte(json))),
		}, nil
	})

	// Get the compose types
	cmd, out, err := root.ExecuteTest("compose", "types")
	defer out.Close()
	require.Nil(t, err)
	require.NotNil(t, out.Stdout)
	require.NotNil(t, out.Stderr)
	require.NotNil(t, cmd)
	assert.Equal(t, cmd, typesCmd)
	stdout, err := ioutil.ReadAll(out.Stdout)
	assert.Nil(t, err)
	assert.Contains(t, string(stdout), "openstack")
	assert.Contains(t, string(stdout), "qcow2")
	stderr, err := ioutil.ReadAll(out.Stderr)
	assert.Nil(t, err)
	assert.Equal(t, []byte(""), stderr)
	assert.Equal(t, "GET", mc.Req.Method)
}

func TestCmdComposeTypesDistro(t *testing.T) {
	// Test the "compose types --distro=test-distro" command
	mc := root.SetupCmdTest(func(request *http.Request) (*http.Response, error) {
		json := `{
    "types": [
        {
            "name": "tar",
            "enabled": true
        },
        {
            "name": "qcow2",
            "enabled": true
        }
    ]
}`

		return &http.Response{
			StatusCode: 200,
			Body:       ioutil.NopCloser(bytes.NewReader([]byte(json))),
		}, nil
	})

	// Get the compose types
	cmd, out, err := root.ExecuteTest("compose", "types", "--distro=test-distro")
	defer out.Close()
	require.Nil(t, err)
	require.NotNil(t, out.Stdout)
	require.NotNil(t, out.Stderr)
	require.NotNil(t, cmd)
	assert.Equal(t, cmd, typesCmd)
	stdout, err := ioutil.ReadAll(out.Stdout)
	assert.Nil(t, err)
	assert.Contains(t, string(stdout), "tar")
	assert.Contains(t, string(stdout), "qcow2")
	stderr, err := ioutil.ReadAll(out.Stderr)
	assert.Nil(t, err)
	assert.Equal(t, []byte(""), stderr)
	assert.Equal(t, "GET", mc.Req.Method)
}

func TestCmdComposeTypesBadDistro(t *testing.T) {
	// Test the "compose types --distro=unknown" command
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
	cmd, out, err := root.ExecuteTest("compose", "types", "--distro=homer")
	defer out.Close()
	require.NotNil(t, err)
	require.NotNil(t, out.Stdout)
	require.NotNil(t, out.Stderr)
	require.NotNil(t, cmd)
	assert.Equal(t, cmd, typesCmd)
	stdout, err := ioutil.ReadAll(out.Stdout)
	assert.Nil(t, err)
	assert.Equal(t, []byte(""), stdout)
	stderr, err := ioutil.ReadAll(out.Stderr)
	assert.Nil(t, err)
	assert.Contains(t, string(stderr), "DistroError: Invalid distro: homer")
	assert.Equal(t, "GET", mc.Req.Method)
}
