// Copyright 2020 by Red Hat, Inc. All rights reserved.
// Use of this source is goverend by the Apache License
// that can be found in the LICENSE file.

package compose

import (
	"fmt"
	"sort"

	"github.com/spf13/cobra"

	"github.com/osbuild/weldr-client/cmd/composer-cli/root"
)

var (
	typesCmd = &cobra.Command{
		Use:   "types",
		Short: "List the available compose types",
		Long:  "List the available compose types",
		RunE:  types,
		Args:  cobra.NoArgs,
	}
	distro string
)

func init() {
	typesCmd.Flags().StringVarP(&distro, "distro", "", "", "Distribution")
	composeCmd.AddCommand(typesCmd)
}

func types(cmd *cobra.Command, args []string) error {
	types, resp, err := root.Client.GetComposeTypes(distro)
	if root.JSONOutput {
		return nil
	}
	if err != nil {
		return root.ExecutionError(cmd, "Types Error: %s", err)
	}
	if resp != nil && !resp.Status {
		return root.ExecutionError(cmd, "Types Error: %s", resp.String())
	}

	sort.Strings(types)
	for i := range types {
		fmt.Println(types[i])
	}

	return nil
}
