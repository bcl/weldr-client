// Copyright 2020 by Red Hat, Inc. All rights reserved.
// Use of this source is goverend by the Apache License
// that can be found in the LICENSE file.

// Package main is the starting point for execution of composer-cli
package main

import (
	"os"

	_ "github.com/osbuild/weldr-client/cmd/composer-cli/blueprints"
	_ "github.com/osbuild/weldr-client/cmd/composer-cli/compose"
	_ "github.com/osbuild/weldr-client/cmd/composer-cli/distros"
	_ "github.com/osbuild/weldr-client/cmd/composer-cli/modules"
	_ "github.com/osbuild/weldr-client/cmd/composer-cli/projects"
	"github.com/osbuild/weldr-client/cmd/composer-cli/root"
	_ "github.com/osbuild/weldr-client/cmd/composer-cli/sources"
	_ "github.com/osbuild/weldr-client/cmd/composer-cli/status"
)

func main() {
	root.Init()

	// Printing errors is handled by the commands or ExecutionError(), just return 1
	if err := root.Execute(); err != nil {
		os.Exit(1)
	}
}
