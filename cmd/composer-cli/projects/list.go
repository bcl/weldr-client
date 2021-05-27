// Copyright 2021 by Red Hat, Inc. All rights reserved.
// Use of this source is goverend by the Apache License
// that can be found in the LICENSE file.

package projects

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/osbuild/weldr-client/cmd/composer-cli/root"
)

var (
	listCmd = &cobra.Command{
		Use:   "list",
		Short: "List available projects",
		Long:  "List available projects",
		RunE:  list,
	}
	distro string
)

func init() {
	listCmd.Flags().StringVarP(&distro, "distro", "", "", "Return results for distribution")
	projectsCmd.AddCommand(listCmd)
}

func list(cmd *cobra.Command, args []string) error {
	projects, resp, err := root.Client.ListProjects(distro)
	if root.JSONOutput {
		return nil
	}
	if err != nil {
		return root.ExecutionError(cmd, "List Error: %s", err)
	}
	if resp != nil && !resp.Status {
		return root.ExecutionError(cmd, "List Error: %s", resp.String())
	}

	for _, p := range projects {
		root.PrintWrap(6, 80, fmt.Sprintf("Name: %s", p.Name))
		root.PrintWrap(9, 80, fmt.Sprintf("Summary: %s", p.Summary))
		root.PrintWrap(10, 80, fmt.Sprintf("Homepage: %s", p.Homepage))
		root.PrintWrap(13, 80, fmt.Sprintf("Description: %s", p.Description))
		fmt.Printf("\n\n")
	}

	return nil
}
