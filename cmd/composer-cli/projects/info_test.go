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

func TestCmdProjectsInfo(t *testing.T) {
	// Test the "modules list" command
	mc := root.SetupCmdTest(func(request *http.Request) (*http.Response, error) {
		json := `{
    "projects": [
        {
            "builds": [
                {
                    "arch": "x86_64",
                    "build_config_ref": "BUILD_CONFIG_REF",
                    "build_env_ref": "BUILD_ENV_REF",
                    "build_time": "2020-07-27T13:17:35",
                    "changelog": "CHANGELOG_NEEDED",
                    "epoch": 0,
                    "metadata": {},
                    "release": "2.fc33",
                    "source": {
                        "license": "GPLv3+",
                        "metadata": {},
                        "source_ref": "SOURCE_REF",
                        "version": "5.0.17"
                    }
                }
            ],
            "description": "The GNU Bourne Again shell (Bash) is a shell or command language\ninterpreter that is compatible with the Bourne shell (sh). Bash\nincorporates useful features from the Korn shell (ksh) and the C shell\n(csh). Most sh scripts can be run by bash without modification.",
            "homepage": "https://www.gnu.org/software/bash",
            "name": "bash",
            "summary": "The GNU Bourne Again shell",
            "upstream_vcs": "UPSTREAM_VCS"
        }
    ]}`

		return &http.Response{
			StatusCode: 200,
			Body:       ioutil.NopCloser(bytes.NewReader([]byte(json))),
		}, nil
	})

	cmd, out, err := root.ExecuteTest("projects", "info", "bash")
	defer out.Close()
	require.Nil(t, err)
	require.NotNil(t, out.Stdout)
	require.NotNil(t, out.Stderr)
	require.NotNil(t, cmd)
	assert.Equal(t, cmd, infoCmd)
	stdout, err := ioutil.ReadAll(out.Stdout)
	assert.Nil(t, err)
	assert.Contains(t, string(stdout), "Summary: The GNU Bourne Again shell")
	assert.Contains(t, string(stdout), "             shell (sh). Bash")
	assert.Contains(t, string(stdout), "     5.0.17-2.fc33.x86_64 at")
	stderr, err := ioutil.ReadAll(out.Stderr)
	assert.Nil(t, err)
	assert.Equal(t, []byte(""), stderr)
	assert.Equal(t, "GET", mc.Req.Method)
	assert.Equal(t, "/api/v1/projects/info/bash", mc.Req.URL.Path)
}

func TestCmdProjectsInfoDistro(t *testing.T) {
	// Test the "modules info --distro=test-distro" command
	mc := root.SetupCmdTest(func(request *http.Request) (*http.Response, error) {
		json := `{
    "projects": [
        {
            "builds": [
                {
                    "arch": "x86_64",
                    "build_config_ref": "BUILD_CONFIG_REF",
                    "build_env_ref": "BUILD_ENV_REF",
                    "build_time": "2020-07-27T13:17:35",
                    "changelog": "CHANGELOG_NEEDED",
                    "epoch": 0,
                    "metadata": {},
                    "release": "2.fc33",
                    "source": {
                        "license": "GPLv3+",
                        "metadata": {},
                        "source_ref": "SOURCE_REF",
                        "version": "5.0.17"
                    }
                }
            ],
            "description": "The GNU Bourne Again shell (Bash) is a shell or command language\ninterpreter that is compatible with the Bourne shell (sh). Bash\nincorporates useful features from the Korn shell (ksh) and the C shell\n(csh). Most sh scripts can be run by bash without modification.",
            "homepage": "https://www.gnu.org/software/bash",
            "name": "bash",
            "summary": "The GNU Bourne Again shell",
            "upstream_vcs": "UPSTREAM_VCS"
        }
    ]}`

		return &http.Response{
			StatusCode: 200,
			Body:       ioutil.NopCloser(bytes.NewReader([]byte(json))),
		}, nil
	})

	cmd, out, err := root.ExecuteTest("projects", "info", "--distro=test-distro", "bash")
	defer out.Close()
	require.Nil(t, err)
	require.NotNil(t, out.Stdout)
	require.NotNil(t, out.Stderr)
	require.NotNil(t, cmd)
	assert.Equal(t, cmd, infoCmd)
	stdout, err := ioutil.ReadAll(out.Stdout)
	assert.Nil(t, err)
	assert.Contains(t, string(stdout), "Summary: The GNU Bourne Again shell")
	assert.Contains(t, string(stdout), "             shell (sh). Bash")
	assert.Contains(t, string(stdout), "     5.0.17-2.fc33.x86_64 at")
	stderr, err := ioutil.ReadAll(out.Stderr)
	assert.Nil(t, err)
	assert.Equal(t, []byte(""), stderr)
	assert.Equal(t, "GET", mc.Req.Method)
	assert.Equal(t, "/api/v1/projects/info/bash", mc.Req.URL.Path)
}

func TestCmdProjectsInfoBadDistro(t *testing.T) {
	// Test the "projects info --distro=homer" command
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
	cmd, out, err := root.ExecuteTest("projects", "info", "--distro=homer", "bash")
	defer out.Close()
	require.NotNil(t, err)
	require.NotNil(t, out.Stdout)
	require.NotNil(t, out.Stderr)
	require.NotNil(t, cmd)
	assert.Equal(t, cmd, infoCmd)
	stdout, err := ioutil.ReadAll(out.Stdout)
	assert.Nil(t, err)
	assert.Equal(t, []byte(""), stdout)
	stderr, err := ioutil.ReadAll(out.Stderr)
	assert.Nil(t, err)
	assert.Contains(t, string(stderr), "DistroError: Invalid distro: homer")
	assert.Equal(t, "GET", mc.Req.Method)
}
