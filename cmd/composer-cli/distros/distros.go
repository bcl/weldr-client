// Copyright 2021 by Red Hat, Inc. All rights reserved.
// Use of this source is goverend by the Apache License
// that can be found in the LICENSE file.

package distros

import (
	"github.com/spf13/cobra"

	"github.com/osbuild/weldr-client/cmd/composer-cli/root"
)

var (
	distrosCmd = &cobra.Command{
		Use:   "distros ...",
		Short: "Manage distributions",
		Long:  "Manage supported distributions on the server",
	}
)

func init() {
	root.AddRootCommand(distrosCmd)
}
