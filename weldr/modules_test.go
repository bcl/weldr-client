// Copyright 2021 by Red Hat, Inc. All rights reserved.
// Use of this source is goverend by the Apache License
// that can be found in the LICENSE file.

// +build integration

package weldr

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestListModules(t *testing.T) {
	modules, r, err := testState.client.ListModules("")
	require.Nil(t, err)
	require.Nil(t, r)
	require.NotNil(t, modules)
	assert.GreaterOrEqual(t, len(modules), 2)
	assert.Equal(t, "rpm", modules[0].Type)
}

func TestListModulesDistro(t *testing.T) {
	modules, r, err := testState.client.ListModules(testState.distros[0])
	require.Nil(t, err)
	require.Nil(t, r)
	require.NotNil(t, modules)
	assert.GreaterOrEqual(t, len(modules), 2)
	assert.Equal(t, "rpm", modules[0].Type)
}

func TestModulesInfo(t *testing.T) {
	modules, r, err := testState.client.ModulesInfo([]string{"bash"}, "")
	require.Nil(t, err)
	require.Nil(t, r)
	require.NotNil(t, modules)
	assert.Equal(t, 1, len(modules))
	assert.GreaterOrEqual(t, len(modules[0].Dependencies), 1)
}

func TestModulesInfoDistro(t *testing.T) {
	modules, r, err := testState.client.ModulesInfo([]string{"bash"}, testState.distros[0])
	require.Nil(t, err)
	require.Nil(t, r)
	require.NotNil(t, modules)
	assert.Equal(t, 1, len(modules))
	assert.GreaterOrEqual(t, len(modules[0].Dependencies), 1)
}

func TestModulesInfoMultiple(t *testing.T) {
	modules, r, err := testState.client.ModulesInfo([]string{"bash", "filesystem", "tmux"}, "")
	require.Nil(t, err)
	require.Nil(t, r)
	require.NotNil(t, modules)
	assert.Equal(t, 3, len(modules))
	for i := range modules {
		assert.GreaterOrEqual(t, len(modules[i].Dependencies), 1)
	}
}

func TestModulesInfoMultipleDistro(t *testing.T) {
	modules, r, err := testState.client.ModulesInfo([]string{"bash", "filesystem", "tmux"}, testState.distros[0])
	require.Nil(t, err)
	require.Nil(t, r)
	require.NotNil(t, modules)
	assert.Equal(t, 3, len(modules))
	for i := range modules {
		assert.GreaterOrEqual(t, len(modules[i].Dependencies), 1)
	}
}

func TestModulesInfoOneError(t *testing.T) {
	modules, r, err := testState.client.ModulesInfo([]string{"bart"}, "")
	require.Nil(t, err)
	require.NotNil(t, r)
	require.Nil(t, modules)
	assert.Equal(t, false, r.Status)
	assert.Equal(t, 1, len(r.Errors))
	assert.Equal(t, "UnknownModule", r.Errors[0].ID)
	assert.Equal(t, "No packages have been found.", r.Errors[0].Msg)
}

func TestModulesInfoMultipleOneError(t *testing.T) {
	modules, r, err := testState.client.ModulesInfo([]string{"bash", "filesystem", "bart"}, "")
	require.Nil(t, err)
	require.Nil(t, r)
	require.NotNil(t, modules)
	assert.Equal(t, 2, len(modules))
	for i := range modules {
		assert.GreaterOrEqual(t, len(modules[i].Dependencies), 1)
	}
}

func TestModulesInfoMultipleOneErrorDistro(t *testing.T) {
	modules, r, err := testState.client.ModulesInfo([]string{"bash", "filesystem", "bart"}, testState.distros[0])
	require.Nil(t, err)
	require.Nil(t, r)
	require.NotNil(t, modules)
	assert.Equal(t, 2, len(modules))
	for i := range modules {
		assert.GreaterOrEqual(t, len(modules[i].Dependencies), 1)
	}
}
