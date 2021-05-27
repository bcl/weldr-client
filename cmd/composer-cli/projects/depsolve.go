// Copyright 2021 by Red Hat, Inc. All rights reserved.
// Use of this source is goverend by the Apache License
// that can be found in the LICENSE file.

package projects

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/osbuild/weldr-client/cmd/composer-cli/root"
)

var (
	depsolveCmd = &cobra.Command{
		Use:   "depsolve PROJECT,...",
		Short: "Show the dependencies of all of the listed projects",
		Long:  "Show the dependencies of all of the listed projects",
		RunE:  depsolve,
		Args:  cobra.MinimumNArgs(1),
	}
)

func init() {
	depsolveCmd.Flags().StringVarP(&distro, "distro", "", "", "Return results for distribution")
	projectsCmd.AddCommand(depsolveCmd)
}

type project struct {
	Name    string
	Epoch   int
	Version string
	Release string
	Arch    string
}

func (p project) String() string {
	if p.Epoch == 0 {
		return fmt.Sprintf("%s-%s-%s.%s", p.Name, p.Version, p.Release, p.Arch)
	}
	return fmt.Sprintf("%d:%s-%s-%s.%s", p.Epoch, p.Name, p.Version, p.Release, p.Arch)
}

func depsolve(cmd *cobra.Command, args []string) (rcErr error) {
	names := root.GetCommaArgs(args)

	deps, errors, err := root.Client.DepsolveProjects(names, distro)
	if err != nil {
		return root.ExecutionError(cmd, "Depsolve Error: %s", err)
	}
	if root.JSONOutput {
		return nil
	}
	if len(errors) > 0 {
		for _, e := range errors {
			fmt.Fprintf(os.Stderr, "ERROR: %s\n", e.String())
		}
		rcErr = root.ExecutionError(cmd, "")
	}

	// Encode it using json
	data := new(bytes.Buffer)
	if err := json.NewEncoder(data).Encode(deps); err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: converting deps: %s\n", err)
		return root.ExecutionError(cmd, "")
	}

	// Decode the dependencies
	var projects []project
	if err = json.Unmarshal(data.Bytes(), &projects); err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: decoding deps: %s\n", err)
		return root.ExecutionError(cmd, "")
	}

	for _, p := range projects {
		fmt.Printf("    %s\n", p)
	}

	return rcErr
}
