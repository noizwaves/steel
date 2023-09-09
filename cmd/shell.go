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

		userValue, err := cmd.Flags().GetBool("user")
		if err != nil {
			return err
		}

		context, err := NewContext(workDirValue)
		if err != nil {
			return err
		}

		return shellAction(context, userValue)
	},
}

func shellAction(ctx *Context, user bool) error {
	zshDotDir, err := prepareZshConfig(ctx.BrewPath, ctx.Brewfile, user)
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
		"_STEEL_SHELL_ACTIVE=true",
		// required by homebrew
		envvar("HOME", os.Getenv("HOME")),
		// inject steel zshrc here
		envvar("ZDOTDIR", zshDotDir),
	}

	if value, found := os.LookupEnv("TERM"); found {
		zshCmd.Env = append(zshCmd.Env, envvar("TERM", value))
	}

	if value, found := os.LookupEnv("HISTFILE"); found {
		zshCmd.Env = append(zshCmd.Env, envvar("HISTFILE", value))
	}

	zshCmd.Stdin = os.Stdin
	zshCmd.Stdout = os.Stdout
	zshCmd.Stderr = os.Stderr

	return zshCmd.Run()
}

func prepareZshConfig(brewPath string, brewfile *impl.Brewfile, user bool) (string, error) {
	zshRcContent := buildZshRc(brewPath, brewfile, user)

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

func envvar(name string, value string) string {
	return fmt.Sprintf("%s=%s", name, value)
}

func lookupZsh() (string, error) {
	return exec.LookPath("zsh")
}

func buildZshRc(brewPath string, brewfile *impl.Brewfile, user bool) string {
	content := bytes.Buffer{}
	// Ensure TERM is set
	content.WriteString(`# Fix backspacing, etc
export TERM=${TERM:-xterm}

`)

	// Pass through HISTFILE via rc file, as ZSH does not accept value from environment
	if value, found := os.LookupEnv("HISTFILE"); found {
		content.WriteString(fmt.Sprintf("export HISTFILE=%s\n\n", value))
	}

	// 2. Set some bling to differentiate shell
	content.WriteString(`# Some bling
PS1="🤘> "

`)

	content.WriteString("# Initialize Homebrew\n")
	content.WriteString(fmt.Sprintf("eval \"$(%s shellenv)\"\n\n", brewPath))

	if brewfile.IncludesPackage("rbenv") {
		content.WriteString(`# Initialize rbenv
eval "$($HOMEBREW_PREFIX/bin/rbenv init - zsh)"

`)
	}

	if brewfile.IncludesPackage("goenv") {
		content.WriteString(`# Initialize goenv
eval "$($HOMEBREW_PREFIX/bin/goenv init -)"

`)
	}

	if user {
		content.WriteString(`
# Load user .zshrc into shell
source ~/.zshrc
`)
	}

	return content.String()
}

func init() {
	rootCmd.AddCommand(shellCmd)

	shellCmd.Flags().BoolP("user", "u", false, "Load ~/.zshrc into environment, making it impure")
}
