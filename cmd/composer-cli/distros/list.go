// Copyright 2021 by Red Hat, Inc. All rights reserved.
// Use of this source is goverend by the Apache License
// that can be found in the LICENSE file.

package distros

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/osbuild/weldr-client/cmd/composer-cli/root"
)

var (
	listCmd = &cobra.Command{
		Use:   "list",
		Short: "List the available distributions",
		Long:  "List the available distributions",
		RunE:  list,
		Args:  cobra.NoArgs,
	}
)

func init() {
	distrosCmd.AddCommand(listCmd)
}

func list(cmd *cobra.Command, args []string) error {
	distros, resp, err := root.Client.ListDistros()
	if root.JSONOutput {
		return nil
	}
	if err != nil {
		return root.ExecutionError(cmd, "Types Error: %s", err)
	}
	if resp != nil && !resp.Status {
		return root.ExecutionError(cmd, "Types Error: %s", resp.String())
	}

	for _, name := range distros {
		fmt.Println(name)
	}

	return nil
}
