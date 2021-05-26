// Copyright 2020-2021 by Red Hat, Inc. All rights reserved.
// Use of this source is goverend by the Apache License
// that can be found in the LICENSE file.

// +build integration

package weldr

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestListComposes(t *testing.T) {
	composes, r, err := testState.client.ListComposes()
	require.Nil(t, err)
	require.Nil(t, r)
	require.NotNil(t, composes)
	assert.GreaterOrEqual(t, len(composes), 4)
}

func TestGetComposeTypes(t *testing.T) {
	types, r, err := testState.client.GetComposeTypes("")
	require.Nil(t, err)
	require.Nil(t, r)
	require.NotNil(t, types)
	assert.Equal(t, 6, len(types))
	assert.Contains(t, types, "openstack")
}

func TestGetComposeTypesDistro(t *testing.T) {
	distros, r, err := testState.client.ListDistros()
	require.Nil(t, err)
	require.Nil(t, r)
	assert.GreaterOrEqual(t, len(distros), 1)

	types, r, err := testState.client.GetComposeTypes(distros[0])
	require.Nil(t, err)
	require.Nil(t, r)
	require.NotNil(t, types)
	assert.Equal(t, 6, len(types))
	assert.Contains(t, types, "openstack")
}

func TestStartCompose(t *testing.T) {
	id, r, err := testState.client.StartComposeTest("cli-test-bp-1", "qcow2", 0, 2)
	require.Nil(t, err)
	require.Nil(t, r)
	assert.Greater(t, len(id), 0)
}

func TestStartComposeSize(t *testing.T) {
	id, r, err := testState.client.StartComposeTest("cli-test-bp-1", "qcow2", 998, 2)
	require.Nil(t, err)
	require.Nil(t, r)
	assert.Greater(t, len(id), 0)
}

func TestStartComposeUpload(t *testing.T) {
	// Need a temporary test file
	tmpProfile, err := ioutil.TempFile("", "test-profile-p*.toml")
	require.Nil(t, err)
	defer os.Remove(tmpProfile.Name())

	_, err = tmpProfile.Write([]byte(`provider = "aws"
[settings]
aws_access_key = "AWS Access Key"
aws_bucket = "AWS Bucket"
aws_region = "AWS Region"
aws_secret_key = "AWS Secret Key"
`))
	require.Nil(t, err)

	id, r, err := testState.client.StartComposeTestUpload("cli-test-bp-1", "qcow2", "test-image", tmpProfile.Name(), 0, 2)
	require.Nil(t, err)
	require.Nil(t, r)
	assert.Greater(t, len(id), 0)
}

func TestStartOSTreeCompose(t *testing.T) {
	id, r, err := testState.client.StartOSTreeComposeTest("cli-test-bp-1", "qcow2", "refid", "parent", "", 0, 2)
	require.Nil(t, err)
	require.Nil(t, r)
	assert.Greater(t, len(id), 0)
}

func TestStartOSTreeComposeUrl(t *testing.T) {
	id, r, err := testState.client.StartOSTreeComposeTest("cli-test-bp-1", "qcow2", "refid", "", "http://weldr.io", 0, 2)
	require.Nil(t, err)
	require.Nil(t, r)
	assert.Greater(t, len(id), 0)
}

func TestStartOSTreeComposeUrlError(t *testing.T) {
	// Sending both the parent url and the parent id should return an error
	id, r, err := testState.client.StartOSTreeComposeTest("cli-test-bp-1", "qcow2", "refid", "parent", "http://weldr.io", 0, 2)
	require.Nil(t, err)
	require.NotNil(t, r)
	assert.False(t, r.Status)
	assert.Equal(t, APIErrorMsg{"OSTreeOptionsError", "Supply at most one of Parent and URL"}, r.Errors[0])
	assert.Equal(t, len(id), 0)
}

func TestStartOSTreeComposeUpload(t *testing.T) {
	// Need a temporary test file
	tmpProfile, err := ioutil.TempFile("", "test-profile-p*.toml")
	require.Nil(t, err)
	defer os.Remove(tmpProfile.Name())

	_, err = tmpProfile.Write([]byte(`provider = "aws"
[settings]
aws_access_key = "AWS Access Key"
aws_bucket = "AWS Bucket"
aws_region = "AWS Region"
aws_secret_key = "AWS Secret Key"
`))
	require.Nil(t, err)

	id, r, err := testState.client.StartOSTreeComposeTestUpload("cli-test-bp-1", "qcow2", "test-image", tmpProfile.Name(), "refid", "", "http://weldr.io", 0, 2)
	require.Nil(t, err)
	require.Nil(t, r)
	assert.Greater(t, len(id), 0)
}

func TestStartComposeUnknownBlueprint(t *testing.T) {
	_, r, err := testState.client.StartCompose("thingy", "qcow2", 0)
	require.Nil(t, err)
	require.NotNil(t, r)
	assert.False(t, r.Status)
	assert.Equal(t, APIErrorMsg{"UnknownBlueprint", "Unknown blueprint name: thingy"}, r.Errors[0])
}

func TestStartComposeBadType(t *testing.T) {
	_, r, err := testState.client.StartCompose("cli-test-bp-1", "punchcard", 0)
	require.Nil(t, err)
	require.NotNil(t, r)
	assert.False(t, r.Status)
	assert.Equal(t, APIErrorMsg{"UnknownComposeType", "Unknown compose type for architecture: punchcard"}, r.Errors[0])
}

func TestStartComposeBadDepsolve(t *testing.T) {
	id, r, err := testState.client.StartCompose("cli-test-bp-3", "qcow2", 0)
	require.Nil(t, err)
	require.NotNil(t, r)
	assert.Equal(t, 0, len(id))
	assert.False(t, r.Status)
	require.Equal(t, 1, len(r.Errors))
	assert.Equal(t, "DepsolveError", r.Errors[0].ID)
}

func TestDeleteComposes(t *testing.T) {
	id, rs, err := testState.client.StartComposeTest("cli-test-bp-1", "qcow2", 0, 2)
	require.Nil(t, err)
	require.Nil(t, rs)
	assert.Greater(t, len(id), 0)

	status, r, err := testState.client.DeleteComposes([]string{id})
	require.Nil(t, err)
	require.Nil(t, r)
	require.NotNil(t, status)
	require.GreaterOrEqual(t, len(status), 1)
	assert.Equal(t, ComposeDeleteV0{ID: id, Status: true}, status[0])
}

func TestDeleteComposesMultiple(t *testing.T) {
	id, rs, err := testState.client.StartComposeTest("cli-test-bp-1", "qcow2", 0, 2)
	require.Nil(t, err)
	require.Nil(t, rs)
	assert.Greater(t, len(id), 0)

	status, r, err := testState.client.DeleteComposes([]string{id, "4b668b1a-e6b8-4dce-8828-4a8e3bef2345"})
	require.Nil(t, err)
	require.NotNil(t, r)
	require.NotNil(t, status)
	require.GreaterOrEqual(t, len(status), 1)
	require.GreaterOrEqual(t, len(r), 1)
	assert.Equal(t, ComposeDeleteV0{ID: id, Status: true}, status[0])
	assert.Equal(t, APIErrorMsg{"UnknownUUID", "compose 4b668b1a-e6b8-4dce-8828-4a8e3bef2345 doesn't exist"}, r[0])
}

func TestCancelFinishedCompose(t *testing.T) {
	id, rs, err := testState.client.StartComposeTest("cli-test-bp-1", "qcow2", 0, 2)
	require.Nil(t, err)
	require.Nil(t, rs)
	assert.Greater(t, len(id), 0)

	status, r, err := testState.client.CancelCompose(id)
	require.Nil(t, err)
	require.NotNil(t, r)
	require.NotNil(t, status)
	assert.False(t, status.Status)
	require.GreaterOrEqual(t, len(r), 1)
	assert.Equal(t, APIErrorMsg{"InternalServerError", "Internal server error: job does not exist"}, r[0])
}

func TestCancelComposeUnknown(t *testing.T) {
	status, r, err := testState.client.CancelCompose("ac188b76-138a-452c-82fb-5cc651986991")
	require.Nil(t, err)
	require.NotNil(t, r)
	require.NotNil(t, status)
	assert.Equal(t, APIErrorMsg{ID: "UnknownUUID", Msg: "Compose ac188b76-138a-452c-82fb-5cc651986991 doesn't exist"}, r[0])
}

func TestComposeLogUnknown(t *testing.T) {
	// This is a difficult one to test, you would have to catch it in the running state, with logs.
	// Test errors instead
	log, r, err := testState.client.ComposeLog("ac188b76-138a-452c-82fb-5cc651986991", 1024)
	require.Nil(t, err)
	require.NotNil(t, r)
	require.Greater(t, len(r.Errors), 0)
	require.NotNil(t, log)
	assert.Equal(t, "", log)
	assert.False(t, r.Status)
	assert.Equal(t, APIErrorMsg{ID: "UnknownUUID", Msg: "Compose ac188b76-138a-452c-82fb-5cc651986991 doesn't exist"}, r.Errors[0])
}

func MakeFinishedCompose(t *testing.T) string {
	// We need a finished compose to download from
	id, r, err := testState.client.StartComposeTest("cli-test-bp-1", "qcow2", 0, 2)
	require.Nil(t, err)
	require.Nil(t, r)
	require.Greater(t, len(id), 0)

	// It should be available immediately, but just in case, try 3 times with a delay
	var found bool
	for i := 0; i < 3 && !found; i++ {
		composes, r, err := testState.client.ListComposes()
		require.Nil(t, err)
		require.Nil(t, r)
		require.NotNil(t, composes)
		for _, c := range composes {
			if c.ID == id && c.Status == "FINISHED" {
				found = true
				break
			}
		}
		time.Sleep(5 * time.Second)
	}
	require.True(t, found)

	return id
}

func TestComposeLogs(t *testing.T) {
	id := MakeFinishedCompose(t)

	// Download the log file
	tf, fn, ct, r, err := testState.client.ComposeLogs(id)
	require.Nil(t, err)
	require.Nil(t, r)
	assert.Equal(t, "application/x-tar", ct)
	assert.Equal(t, fmt.Sprintf("%s-logs.tar", id), fn)
	require.Greater(t, len(tf), 0)
	_, err = os.Stat(tf)
	require.Nil(t, err)
	os.Remove(tf)
}

func TestComposeLogsUnknown(t *testing.T) {
	// Test handling of unknown uuid
	tf, fn, ct, r, err := testState.client.ComposeLogs("90eafe5a-00f3-40f8-8416-d6809a94e25d")
	require.Nil(t, err)
	require.NotNil(t, r)
	assert.Equal(t, false, r.Status)
	assert.Equal(t, 1, len(r.Errors))
	assert.Equal(t, APIErrorMsg{"UnknownUUID", "Compose 90eafe5a-00f3-40f8-8416-d6809a94e25d doesn't exist"}, r.Errors[0])
	assert.Equal(t, "", ct)
	assert.Equal(t, "", fn)
	assert.Equal(t, "", tf)
}

func TestComposeMetadata(t *testing.T) {
	id := MakeFinishedCompose(t)

	// Download the metadata file
	tf, fn, ct, r, err := testState.client.ComposeMetadata(id)
	require.Nil(t, err)
	require.Nil(t, r)
	assert.Equal(t, "application/x-tar", ct)
	assert.Equal(t, fmt.Sprintf("%s-metadata.tar", id), fn)
	require.Greater(t, len(tf), 0)
	_, err = os.Stat(tf)
	require.Nil(t, err)
	os.Remove(tf)
}

func TestComposeMetadataUnknown(t *testing.T) {
	// Test handling of unknown uuid
	tf, fn, ct, r, err := testState.client.ComposeMetadata("90eafe5a-00f3-40f8-8416-d6809a94e25d")
	require.Nil(t, err)
	require.NotNil(t, r)
	assert.Equal(t, false, r.Status)
	assert.Equal(t, 1, len(r.Errors))
	assert.Equal(t, APIErrorMsg{"UnknownUUID", "Compose 90eafe5a-00f3-40f8-8416-d6809a94e25d doesn't exist"}, r.Errors[0])
	assert.Equal(t, "", ct)
	assert.Equal(t, "", fn)
	assert.Equal(t, "", tf)
}

func TestComposeResults(t *testing.T) {
	id := MakeFinishedCompose(t)

	// Download the results file
	tf, fn, ct, r, err := testState.client.ComposeResults(id)
	require.Nil(t, err)
	require.Nil(t, r)
	assert.Equal(t, "application/x-tar", ct)
	assert.Equal(t, fmt.Sprintf("%s.tar", id), fn)
	require.Greater(t, len(tf), 0)
	_, err = os.Stat(tf)
	require.Nil(t, err)
	os.Remove(tf)
}

func TestComposeResultsUnknown(t *testing.T) {
	// Test handling of unknown uuid
	tf, fn, ct, r, err := testState.client.ComposeResults("90eafe5a-00f3-40f8-8416-d6809a94e25d")
	require.Nil(t, err)
	require.NotNil(t, r)
	assert.Equal(t, false, r.Status)
	assert.Equal(t, 1, len(r.Errors))
	assert.Equal(t, APIErrorMsg{"UnknownUUID", "Compose 90eafe5a-00f3-40f8-8416-d6809a94e25d doesn't exist"}, r.Errors[0])
	assert.Equal(t, "", ct)
	assert.Equal(t, "", fn)
	assert.Equal(t, "", tf)
}

func TestComposeImageError(t *testing.T) {
	id := MakeFinishedCompose(t)

	// Test composes don't actually have an image file, so this is going to fail with an error.
	// test that instead.

	// Download the image file
	tf, fn, ct, r, err := testState.client.ComposeImage(id)
	require.Nil(t, err)
	require.NotNil(t, r)
	assert.Equal(t, false, r.Status)
	assert.Equal(t, 1, len(r.Errors))
	assert.Equal(t, "InternalServerError", r.Errors[0].ID)
	assert.Contains(t, r.Errors[0].Msg, "Error accessing image file for compose")
	assert.Equal(t, "", ct)
	assert.Equal(t, "", fn)
	assert.Equal(t, "", tf)
}

func TestComposeImageUnknown(t *testing.T) {
	// Test handling of unknown uuid
	tf, fn, ct, r, err := testState.client.ComposeImage("90eafe5a-00f3-40f8-8416-d6809a94e25d")
	require.Nil(t, err)
	require.NotNil(t, r)
	assert.Equal(t, false, r.Status)
	assert.Equal(t, 1, len(r.Errors))
	assert.Equal(t, APIErrorMsg{"UnknownUUID", "Compose 90eafe5a-00f3-40f8-8416-d6809a94e25d doesn't exist"}, r.Errors[0])
	assert.Equal(t, "", ct)
	assert.Equal(t, "", fn)
	assert.Equal(t, "", tf)
}

func TestComposeInfo(t *testing.T) {
	id := MakeFinishedCompose(t)

	// Get the details about the compose
	info, r, err := testState.client.ComposeInfo(id)
	require.Nil(t, err)
	require.Nil(t, r)
	require.NotNil(t, info)
	assert.Equal(t, id, info.ID)
	assert.Equal(t, "qcow2", info.ComposeType)
	assert.Equal(t, "FINISHED", info.QueueStatus)
	assert.Equal(t, "cli-test-bp-1", info.Blueprint.Name)
	require.Greater(t, len(info.Blueprint.Packages), 0)
}

func TestComposeInfoUnknown(t *testing.T) {
	// Get the details about the compose
	info, r, err := testState.client.ComposeInfo("fcb032c5-5734-4cda-bc60-c4e72c0f76fd")
	require.Nil(t, err)
	require.NotNil(t, r)
	assert.Equal(t, false, r.Status)
	assert.Equal(t, 1, len(r.Errors))
	assert.Equal(t, APIErrorMsg{"UnknownUUID", "fcb032c5-5734-4cda-bc60-c4e72c0f76fd is not a valid build uuid"}, r.Errors[0])
	require.Equal(t, ComposeInfoV0{}, info)
}
