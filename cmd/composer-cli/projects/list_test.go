// Copyright 2021 by Red Hat, Inc. All rights reserved.
// Use of this source is goverend by the Apache License
// that can be found in the LICENSE file.

package projects

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/osbuild/weldr-client/cmd/composer-cli/root"
)

func TestCmdProjectsList(t *testing.T) {
	// Test the "projects list" command
	mc := root.SetupCmdTest(func(request *http.Request) (*http.Response, error) {
		query := request.URL.Query()
		v := query.Get("limit")
		limit, _ := strconv.ParseUint(v, 10, 64)
		var json string
		if limit == 0 {
			json = `{ "limit": 0, "offset": 0, "projects": [], "total": 1 }`
		} else {
			json = `{
    "limit": 1,
    "offset": 0,
    "projects": [
        {
            "builds": [
                {
                    "arch": "x86_64",
                    "build_config_ref": "BUILD_CONFIG_REF",
                    "build_env_ref": "BUILD_ENV_REF",
                    "build_time": "2020-07-31T23:48:35",
                    "changelog": "CHANGELOG_NEEDED",
                    "epoch": 0,
                    "metadata": {},
                    "release": "21.fc33",
                    "source": {
                        "license": "GPLv2+ and BSD and MIT and IBM",
                        "metadata": {},
                        "source_ref": "SOURCE_REF",
                        "version": "0.0.23b"
                    }
                },
                {
                    "arch": "x86_64",
                    "build_config_ref": "BUILD_CONFIG_REF",
                    "build_env_ref": "BUILD_ENV_REF",
                    "build_time": "2021-03-02T12:09:34",
                    "changelog": "CHANGELOG_NEEDED",
                    "epoch": 0,
                    "metadata": {},
                    "release": "2.fc33",
                    "source": {
                        "license": "GPLv2+ and BSD and MIT and IBM and MPLv2.0",
                        "metadata": {},
                        "source_ref": "SOURCE_REF",
                        "version": "0.0.24b"
                    }
                }
            ],
            "description": "0 A.D. (pronounced \"zero ey-dee\") is a free, open-source, cross-platform\nreal-time strategy (RTS) game of ancient warfare. In short, it is a\nhistorically-based war/economy game that allows players to relive or rewrite\nthe history of Western civilizations, focusing on the years between 500 B.C.\nand 500 A.D. The project is highly ambitious, involving state-of-the-art 3D\ngraphics, detailed artwork, sound, and a flexible and powerful custom-built\ngame engine.\n\nThe game has been in development by Wildfire Games (WFG), a group of volunteer,\nhobbyist game developers, since 2001.",
            "homepage": "http://play0ad.com",
            "name": "0ad",
"summary": "Cross-Platform RTS Game of Ancient Warfare",
            "upstream_vcs": "UPSTREAM_VCS"
        }],
	"total": 1}`
		}
		return &http.Response{
			StatusCode: 200,
			Body:       ioutil.NopCloser(bytes.NewReader([]byte(json))),
		}, nil
	})

	cmd, out, err := root.ExecuteTest("projects", "list")
	defer out.Close()
	require.Nil(t, err)
	require.NotNil(t, out.Stdout)
	require.NotNil(t, out.Stderr)
	require.NotNil(t, cmd)
	assert.Equal(t, cmd, listCmd)
	stdout, err := ioutil.ReadAll(out.Stdout)
	assert.Nil(t, err)
	assert.Contains(t, string(stdout), "Name: 0ad\n")
	assert.Contains(t, string(stdout), "Homepage: http://play0ad.com\n")
	assert.Contains(t, string(stdout), "             open-source, cross-platform real-time")
	stderr, err := ioutil.ReadAll(out.Stderr)
	assert.Nil(t, err)
	assert.Equal(t, []byte(""), stderr)
	assert.Equal(t, "GET", mc.Req.Method)
	assert.Equal(t, "/api/v1/projects/list", mc.Req.URL.Path)
}

func TestCmdProjectsListDistro(t *testing.T) {
	// Test the "projects list --distro=test-distro" command
	mc := root.SetupCmdTest(func(request *http.Request) (*http.Response, error) {
		query := request.URL.Query()
		v := query.Get("limit")
		limit, _ := strconv.ParseUint(v, 10, 64)
		var json string
		if limit == 0 {
			json = `{ "limit": 0, "offset": 0, "projects": [], "total": 1 }`
		} else {
			json = `{
    "limit": 1,
    "offset": 0,
    "projects": [
        {
            "builds": [
                {
                    "arch": "x86_64",
                    "build_config_ref": "BUILD_CONFIG_REF",
                    "build_env_ref": "BUILD_ENV_REF",
                    "build_time": "2020-07-31T23:48:35",
                    "changelog": "CHANGELOG_NEEDED",
                    "epoch": 0,
                    "metadata": {},
                    "release": "21.fc33",
                    "source": {
                        "license": "GPLv2+ and BSD and MIT and IBM",
                        "metadata": {},
                        "source_ref": "SOURCE_REF",
                        "version": "0.0.23b"
                    }
                },
                {
                    "arch": "x86_64",
                    "build_config_ref": "BUILD_CONFIG_REF",
                    "build_env_ref": "BUILD_ENV_REF",
                    "build_time": "2021-03-02T12:09:34",
                    "changelog": "CHANGELOG_NEEDED",
                    "epoch": 0,
                    "metadata": {},
                    "release": "2.fc33",
                    "source": {
                        "license": "GPLv2+ and BSD and MIT and IBM and MPLv2.0",
                        "metadata": {},
                        "source_ref": "SOURCE_REF",
                        "version": "0.0.24b"
                    }
                }
            ],
            "description": "0 A.D. (pronounced \"zero ey-dee\") is a free, open-source, cross-platform\nreal-time strategy (RTS) game of ancient warfare. In short, it is a\nhistorically-based war/economy game that allows players to relive or rewrite\nthe history of Western civilizations, focusing on the years between 500 B.C.\nand 500 A.D. The project is highly ambitious, involving state-of-the-art 3D\ngraphics, detailed artwork, sound, and a flexible and powerful custom-built\ngame engine.\n\nThe game has been in development by Wildfire Games (WFG), a group of volunteer,\nhobbyist game developers, since 2001.",
            "homepage": "http://play0ad.com",
            "name": "0ad",
"summary": "Cross-Platform RTS Game of Ancient Warfare",
            "upstream_vcs": "UPSTREAM_VCS"
        }],
	"total": 1}`
		}
		return &http.Response{
			StatusCode: 200,
			Body:       ioutil.NopCloser(bytes.NewReader([]byte(json))),
		}, nil
	})

	cmd, out, err := root.ExecuteTest("projects", "list", "--distro=test-distro")
	defer out.Close()
	require.Nil(t, err)
	require.NotNil(t, out.Stdout)
	require.NotNil(t, out.Stderr)
	require.NotNil(t, cmd)
	assert.Equal(t, cmd, listCmd)
	stdout, err := ioutil.ReadAll(out.Stdout)
	assert.Nil(t, err)
	assert.Contains(t, string(stdout), "Name: 0ad\n")
	assert.Contains(t, string(stdout), "Homepage: http://play0ad.com\n")
	assert.Contains(t, string(stdout), "             open-source, cross-platform real-time")
	stderr, err := ioutil.ReadAll(out.Stderr)
	assert.Nil(t, err)
	assert.Equal(t, []byte(""), stderr)
	assert.Equal(t, "GET", mc.Req.Method)
	assert.Equal(t, "/api/v1/projects/list", mc.Req.URL.Path)
}

func TestCmdProjectsListBadDistro(t *testing.T) {
	// Test the "projects list --distro=homer" command
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
	cmd, out, err := root.ExecuteTest("projects", "list", "--distro=homer")
	defer out.Close()
	require.NotNil(t, err)
	require.NotNil(t, out.Stdout)
	require.NotNil(t, out.Stderr)
	require.NotNil(t, cmd)
	assert.Equal(t, cmd, listCmd)
	stdout, err := ioutil.ReadAll(out.Stdout)
	assert.Nil(t, err)
	assert.Equal(t, []byte(""), stdout)
	stderr, err := ioutil.ReadAll(out.Stderr)
	assert.Nil(t, err)
	assert.Contains(t, string(stderr), "DistroError: Invalid distro: homer")
	assert.Equal(t, "GET", mc.Req.Method)
}
