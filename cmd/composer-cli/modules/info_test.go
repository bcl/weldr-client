// Copyright 2021 by Red Hat, Inc. All rights reserved.
// Use of this source is goverend by the Apache License
// that can be found in the LICENSE file.

package modules

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/osbuild/weldr-client/cmd/composer-cli/root"
)

func TestCmdModulesInfo(t *testing.T) {
	// Test the "modules info" command
	mc := root.SetupCmdTest(func(request *http.Request) (*http.Response, error) {
		json := `{
    "modules": [
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
            "dependencies": [
                {
                    "arch": "noarch",
                    "check_gpg": true,
                    "checksum": "sha256:f4efaa5bc8382246d8230ece8bacebd3c29eb9fd52b509b1e6575e643953851b",
                    "epoch": 0,
                    "name": "basesystem",
                    "release": "10.fc33",
                    "remote_location": "http://mirror.web-ster.com/fedora/releases/33/Everything/x86_64/os/Packages/b/basesystem-11-10.fc33.noarch.rpm",
                    "version": "11"
                },
                {
                    "arch": "x86_64",
                    "check_gpg": true,
                    "checksum": "sha256:c59a621f3cdd5e073b3c1ef9cd8fd9d7e02d77d94be05330390eac05f77b5b60",
                    "epoch": 0,
                    "name": "bash",
                    "release": "2.fc33",
                    "remote_location": "http://mirror.web-ster.com/fedora/releases/33/Everything/x86_64/os/Packages/b/bash-5.0.17-2.fc33.x86_64.rpm",
                    "version": "5.0.17"
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

	cmd, out, err := root.ExecuteTest("modules", "info", "bash")
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
	assert.Contains(t, string(stdout), "     basesystem-11-10.fc33.noarch")
	stderr, err := ioutil.ReadAll(out.Stderr)
	assert.Nil(t, err)
	assert.Equal(t, []byte(""), stderr)
	assert.Equal(t, "GET", mc.Req.Method)
	assert.Equal(t, "/api/v1/modules/info/bash", mc.Req.URL.Path)
}

func TestCmdModulesInfoDistro(t *testing.T) {
	// Test the "modules info --distro=test-distro" command
	mc := root.SetupCmdTest(func(request *http.Request) (*http.Response, error) {
		json := `{
    "modules": [
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
            "dependencies": [
                {
                    "arch": "noarch",
                    "check_gpg": true,
                    "checksum": "sha256:f4efaa5bc8382246d8230ece8bacebd3c29eb9fd52b509b1e6575e643953851b",
                    "epoch": 0,
                    "name": "basesystem",
                    "release": "10.fc33",
                    "remote_location": "http://mirror.web-ster.com/fedora/releases/33/Everything/x86_64/os/Packages/b/basesystem-11-10.fc33.noarch.rpm",
                    "version": "11"
                },
                {
                    "arch": "x86_64",
                    "check_gpg": true,
                    "checksum": "sha256:c59a621f3cdd5e073b3c1ef9cd8fd9d7e02d77d94be05330390eac05f77b5b60",
                    "epoch": 0,
                    "name": "bash",
                    "release": "2.fc33",
                    "remote_location": "http://mirror.web-ster.com/fedora/releases/33/Everything/x86_64/os/Packages/b/bash-5.0.17-2.fc33.x86_64.rpm",
                    "version": "5.0.17"
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

	cmd, out, err := root.ExecuteTest("modules", "info", "--distro=test-distro", "bash")
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
	assert.Contains(t, string(stdout), "     basesystem-11-10.fc33.noarch")
	stderr, err := ioutil.ReadAll(out.Stderr)
	assert.Nil(t, err)
	assert.Equal(t, []byte(""), stderr)
	assert.Equal(t, "GET", mc.Req.Method)
	assert.Equal(t, "/api/v1/modules/info/bash", mc.Req.URL.Path)
}

func TestCmdModulesInfoBadDistro(t *testing.T) {
	// Test the "modules info --distro=homer" command
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
	cmd, out, err := root.ExecuteTest("modules", "info", "--distro=homer", "bash")
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
