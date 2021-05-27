// Copyright 2021 by Red Hat, Inc. All rights reserved.
// Use of this source is goverend by the Apache License
// that can be found in the LICENSE file.

package projects

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/osbuild/weldr-client/cmd/composer-cli/root"
)

var (
	infoCmd = &cobra.Command{
		Use:   "info PROJECT,...",
		Short: "Show detailed info about the listed projects",
		Long:  "Show detailed info about the listed projects",
		RunE:  info,
		Args:  cobra.MinimumNArgs(1),
	}
)

func init() {
	infoCmd.Flags().StringVarP(&distro, "distro", "", "", "Return results for distribution")
	projectsCmd.AddCommand(infoCmd)
}

func info(cmd *cobra.Command, args []string) error {
	names := root.GetCommaArgs(args)

	projects, resp, err := root.Client.ProjectsInfo(names, distro)
	if root.JSONOutput {
		return nil
	}
	if err != nil {
		return root.ExecutionError(cmd, "Info Error: %s", err)
	}
	if resp != nil && !resp.Status {
		return root.ExecutionError(cmd, strings.Join(resp.AllErrors(), "\n"))
	}

	for _, p := range projects {
		root.PrintWrap(6, 80, fmt.Sprintf("Name: %s", p.Name))
		root.PrintWrap(9, 80, fmt.Sprintf("Summary: %s", p.Summary))
		root.PrintWrap(10, 80, fmt.Sprintf("Homepage: %s", p.Homepage))
		root.PrintWrap(13, 80, fmt.Sprintf("Description: %s", p.Description))
		fmt.Println("Builds: ")
		for _, b := range p.Builds {
			fmt.Println("    ", b)
		}
		fmt.Printf("\n\n")
	}

	return nil
}
