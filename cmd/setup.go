/*
Copyright © 2023 Adam Neumann
*/
package cmd

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
)

var setupCmd = &cobra.Command{
	Use:          "setup",
	Short:        "Installs all dependencies required by the application",
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		workDir, err := cmd.Flags().GetString("workdir")
		if err != nil {
			return err
		}

		context, err := NewContext(workDir)
		if err != nil {
			return err
		}

		return setupAction(context)
	},
}

func setupAction(ctx *Context) error {
	err := runBrewBundleInstall(ctx.BrewPath, ctx.WorkDir)
	if err != nil {
		return err
	}

	if ctx.Brewfile.IncludesPackage("rbenv") {
		return runBrewBundleExecRbenvInstall(ctx.BrewPath)
	}

	if ctx.Brewfile.IncludesPackage("goenv") {
		return runBrewBundleExecGoenvInstall(ctx.BrewPath)
	}

	return nil
}

const dummyBrewArg = "brew"

func runBrewBundleExecRbenvInstall(brewPath string) error {
	cmdOut := bytes.Buffer{}
	cmdErr := bytes.Buffer{}
	cmd := exec.Cmd{
		Path:   brewPath,
		Args:   []string{dummyBrewArg, "bundle", "exec", "--", "rbenv", "install", "--skip-existing"},
		Dir:    workDir,
		Stdout: &cmdOut,
		Stderr: &cmdErr,
	}

	err := cmd.Run()
	if cmdOut.Len() > 0 {
		fmt.Fprintf(os.Stdout, "**** OUTPUT START ****\n%s\n**** OUTPUT END ****\n", cmdOut.String())
	}
	if err != nil {
		fmt.Fprintf(os.Stdout, "**** ERROR START ****\n%s\n**** ERROR END ****\n", cmdErr.String())
	}

	return err
}

func runBrewBundleExecGoenvInstall(brewPath string) error {
	cmdOut := bytes.Buffer{}
	cmdErr := bytes.Buffer{}
	cmd := exec.Cmd{
		Path:   brewPath,
		Args:   []string{dummyBrewArg, "bundle", "exec", "--", "goenv", "install", "--skip-existing"},
		Dir:    workDir,
		Stdout: &cmdOut,
		Stderr: &cmdErr,
	}

	err := cmd.Run()
	if cmdOut.Len() > 0 {
		fmt.Fprintf(os.Stdout, "**** OUTPUT START ****\n%s\n**** OUTPUT END ****\n", cmdOut.String())
	}
	if err != nil {
		fmt.Fprintf(os.Stdout, "**** ERROR START ****\n%s\n**** ERROR END ****\n", cmdErr.String())
	}

	return err
}

func runBrewBundleInstall(brewPath string, workDir string) error {
	cmdOut := bytes.Buffer{}
	cmdErr := bytes.Buffer{}
	cmd := exec.Cmd{
		Path:   brewPath,
		Args:   []string{dummyBrewArg, "bundle", "install", "--no-lock"},
		Dir:    workDir,
		Stdout: &cmdOut,
		Stderr: &cmdErr,
	}

	err := cmd.Run()
	if cmdOut.Len() > 0 {
		fmt.Fprintf(os.Stdout, "**** OUTPUT START ****\n%s\n**** OUTPUT END ****\n", cmdOut.String())
	}
	if err != nil {
		fmt.Fprintf(os.Stdout, "**** ERROR START ****\n%s\n**** ERROR END ****\n", cmdErr.String())
	}

	return err
}

func init() {
	rootCmd.AddCommand(setupCmd)
}
