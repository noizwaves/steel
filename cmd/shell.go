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

// shellCmd represents the shell command
var shellCmd = &cobra.Command{
	Use:          "shell",
	Short:        "An interactive shell with dependencies initialized",
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

		return shellAction(workDir)
	},
}

func shellAction(workDir string) error {
	brewfilePath := filepath.Join(workDir, "Brewfile")
	brewfile, err := impl.LoadBrewfile(brewfilePath)
	if err != nil {
		return err
	}

	zshDotDir, err := prepareZshConfig(brewfile)
	if err != nil {
		return err
	}

	zshPath, err := lookupZsh()
	if err != nil {
		return err
	}

	// start interactive zsh
	zshCmd := exec.Command(zshPath)
	zshCmd.Dir = workDir
	zshCmd.Env = []string{
		// required by homebrew
		"HOME=" + os.Getenv("HOME"),
		fmt.Sprintf("ZDOTDIR=%s", zshDotDir),
	}

	zshCmd.Stdin = os.Stdin
	zshCmd.Stdout = os.Stdout
	zshCmd.Stderr = os.Stderr

	return zshCmd.Run()
}

func prepareZshConfig(brewfile *impl.Brewfile) (string, error) {
	zshRcContent := buildZshRc(brewfile)

	zshDotDir, err := os.MkdirTemp("", "steel_zsh_*")
	if err != nil {
		return "", err
	}
	zshRcPath := filepath.Join(zshDotDir, ".zshrc")
	err = os.WriteFile(zshRcPath, []byte(zshRcContent), 0666)
	if err != nil {
		return "", err
	}

	return zshDotDir, nil
}

func lookupZsh() (string, error) {
	return exec.LookPath("zsh")
}

func buildZshRc(brewfile *impl.Brewfile) string {
	content := bytes.Buffer{}
	// 1. Set TERM
	content.WriteString(`# Fix backspacing, etc
export TERM=xterm
`)

	// 2. Set some bling to differentiate shell
	content.WriteString(`# Some bling
PS1="🤘> "
`)

	// TODO: look up brew path dynamically
	content.WriteString(`# Initialize Homebrew
eval "$(/opt/homebrew/bin/brew shellenv)"
`)

	if brewfile.IncludesPackage("rbenv") {
		content.WriteString(`# Initialize rbenv
eval "$($HOMEBREW_PREFIX/bin/rbenv init - zsh)"
`)
	}

	return content.String()
}

func init() {
	rootCmd.AddCommand(shellCmd)
}
