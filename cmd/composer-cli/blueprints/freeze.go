// Copyright 2020 by Red Hat, Inc. All rights reserved.
// Use of this source is goverend by the Apache License
// that can be found in the LICENSE file.

package blueprints

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/spf13/cobra"

	"weldr-client/cmd/composer-cli/root"
)

var (
	freezeCmd = &cobra.Command{
		Use:   "freeze BLUEPRINT,...",
		Short: "Show the blueprints depsolved package and module versions",
		Long:  "Show the blueprints depsolved package and module versions",
		RunE:  freeze,
		Args:  cobra.MinimumNArgs(1),
	}
	freezeShowCmd = &cobra.Command{
		Use:   "show BLUEPRINT,...",
		Short: "Show the complete frozen blueprints TOML format",
		Long:  "Show the complete blueprints with their depsolved packages and modules in TOML format",
		RunE:  freezeShow,
		Args:  cobra.MinimumNArgs(1),
	}
	freezeSaveCmd = &cobra.Command{
		Use:   "save BLUEPRINT,...",
		Short: "Save the frozen blueprints to a TOML file",
		Long:  "Save the complete blueprints with their depsolved packages and modules in TOML formatted files named BLUEPRINT-NAME.frozen.toml",
		RunE:  freezeSave,
		Args:  cobra.MinimumNArgs(1),
	}
)

func init() {
	blueprintsCmd.AddCommand(freezeCmd)
	freezeCmd.AddCommand(freezeShowCmd)
	freezeCmd.AddCommand(freezeSaveCmd)
}

// blueprintParts is Used to decode the parts of the blueprint to display
type blueprintParts struct {
	Name    string
	Version string
	Modules []struct {
		Name    string
		Version string
	}
	Packages []struct {
		Name    string
		Version string
	}
}

func freeze(cmd *cobra.Command, args []string) (rcErr error) {
	names := root.GetCommaArgs(args)
	bps, errors, err := root.Client.GetFrozenBlueprintsJSON(names)
	if err != nil {
		return root.ExecutionError(cmd, "Save Error: %s", err)
	}
	if root.JSONOutput {
		return nil
	}
	if len(errors) > 0 {
		for _, e := range errors {
			fmt.Println(e.String())
		}
		rcErr = root.ExecutionError(cmd, "")
	}

	for _, bp := range bps {
		// Encode it using toml
		data := new(bytes.Buffer)
		if err := toml.NewEncoder(data).Encode(bp); err != nil {
			rcErr = root.ExecutionError(cmd, "Error converting blueprint: %s", err)
		}
		// Decode the parts we care about into blueprintParts
		var parts blueprintParts
		if _, err := toml.Decode(data.String(), &parts); err != nil {
			rcErr = root.ExecutionError(cmd, "Error decoding blueprint: %s", err)
		}

		if len(parts.Version) > 0 {
			fmt.Printf("blueprint: %s v%s\n", parts.Name, parts.Version)
		} else {
			fmt.Printf("blueprint: %s\n", parts.Name)
		}
		for _, m := range parts.Modules {
			fmt.Printf("    %s-%s\n", m.Name, m.Version)
		}
		for _, p := range parts.Packages {
			fmt.Printf("    %s-%s\n", p.Name, p.Version)
		}
	}

	// If there were any errors, even if other blueprints succeeded, it returns an error
	return rcErr
}

func freezeShow(cmd *cobra.Command, args []string) error {
	// TODO -- check root.JSONOutput and do a json request and output as a map with names as keys
	names := root.GetCommaArgs(args)
	blueprints, resp, err := root.Client.GetFrozenBlueprintsTOML(names)
	if resp != nil || err != nil {
		return root.ExecutionError(cmd, "Show Error: %s", err)
	}
	for _, bp := range blueprints {
		fmt.Println(bp)
	}

	return nil
}

func freezeSave(cmd *cobra.Command, args []string) (rcErr error) {
	names := root.GetCommaArgs(args)
	bps, errors, err := root.Client.GetFrozenBlueprintsJSON(names)
	if err != nil {
		return root.ExecutionError(cmd, "Save Error: %s", err)
	}
	if root.JSONOutput {
		return nil
	}
	if len(errors) > 0 {
		for _, e := range errors {
			fmt.Println(e.String())
		}
		rcErr = root.ExecutionError(cmd, "")
	}

	for _, bp := range bps {
		name, ok := bp.(map[string]interface{})["name"].(string)
		if !ok {
			rcErr = root.ExecutionError(cmd, "ERROR: no 'name' in blueprint")
			continue
		}

		// Save to a file in the current directory, replace spaces with - and
		// remove anything that looks like path separators or path traversal.
		filename := strings.ReplaceAll(name, " ", "-") + ".frozen.toml"
		filename = filepath.Base(filename)
		if filename == "/" || filename == "." || filename == ".." {
			rcErr = root.ExecutionError(cmd, "ERROR: Invalid blueprint filename: %s\n", name)
			continue
		}
		f, err := os.OpenFile(filename, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0600)
		if err != nil {
			rcErr = root.ExecutionError(cmd, "Error opening file %s: %s\n", "file.toml", err)
			continue
		}
		defer f.Close()
		err = toml.NewEncoder(f).Encode(bp)
		if err != nil {
			rcErr = root.ExecutionError(cmd, "Error encoding TOML file: %s\n", err)
		}
	}

	// If there were any errors, even if other blueprints succeeded, it returns an error
	return rcErr
}