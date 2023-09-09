/*
Copyright © 2023 Adam Neumann
*/
package cmd

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/noizwaves/steel/impl"
	"github.com/spf13/cobra"
)

var setupCmd = &cobra.Command{
	Use:          "setup",
	Short:        "Installs all dependencies required by the application",
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		workDirValue, err := cmd.Flags().GetString("workdir")
		if err != nil {
			return err
		}

		workDir, err := filepath.Abs(workDirValue)
		if err != nil {
			return err
		}

		// TODO: validate system requirements met (Homebrew)
		// TODO: validate workDir requirements met (Brewfile, etc)

		return setupAction(workDir)
	},
}

func setupAction(workDir string) error {
	brewfilePath := filepath.Join(workDir, "Brewfile")
	brewfile, err := impl.LoadBrewfile(brewfilePath)
	if err != nil {
		return err
	}

	err = runBrewBundleInstall(workDir)
	if err != nil {
		return err
	}

	if brewfile.IncludesPackage("rbenv") {
		return runBrewBundleExecRbenvInstall()
	}

	return nil
}

const dummyBrewArg = "brew"

func runBrewBundleExecRbenvInstall() error {
	brewPath, err := lookupBrew()
	if err != nil {
		return err
	}

	cmdOut := bytes.Buffer{}
	cmdErr := bytes.Buffer{}
	cmd := exec.Cmd{
		Path:   brewPath,
		Args:   []string{dummyBrewArg, "bundle", "exec", "--", "rbenv", "install", "--skip-existing"},
		Dir:    workDir,
		Stdout: &cmdOut,
		Stderr: &cmdErr,
	}

	err = cmd.Run()
	if cmdOut.Len() > 0 {
		fmt.Fprintf(os.Stdout, "**** OUTPUT START ****\n%s\n**** OUTPUT END ****\n", cmdOut.String())
	}
	if err != nil {
		fmt.Fprintf(os.Stdout, "**** ERROR START ****\n%s\n**** ERROR END ****\n", cmdErr.String())
	}

	return err
}

func runBrewBundleInstall(workDir string) error {
	brewPath, err := lookupBrew()
	if err != nil {
		return err
	}

	cmdOut := bytes.Buffer{}
	cmdErr := bytes.Buffer{}
	cmd := exec.Cmd{
		Path:   brewPath,
		Args:   []string{dummyBrewArg, "bundle", "install", "--no-lock"},
		Dir:    workDir,
		Stdout: &cmdOut,
		Stderr: &cmdErr,
	}

	err = cmd.Run()
	if cmdOut.Len() > 0 {
		fmt.Fprintf(os.Stdout, "**** OUTPUT START ****\n%s\n**** OUTPUT END ****\n", cmdOut.String())
	}
	if err != nil {
		fmt.Fprintf(os.Stdout, "**** ERROR START ****\n%s\n**** ERROR END ****\n", cmdErr.String())
	}

	return err
}

func lookupBrew() (string, error) {
	return exec.LookPath("brew")
}

func init() {
	rootCmd.AddCommand(setupCmd)
}
