package cmd

import (
	"os/exec"
	"path/filepath"

	"github.com/noizwaves/steel/impl"
)

type Context struct {
	WorkDir  string
	Brewfile *impl.Brewfile
	BrewPath string
}

func NewContext(workDir string) (*Context, error) {
	workDir, err := filepath.Abs(workDir)
	if err != nil {
		return nil, err
	}

	brewfile, err := impl.LoadBrewfile(filepath.Join(workDir, "Brewfile"))
	if err != nil {
		return nil, err
	}

	brewPath, err := exec.LookPath("brew")
	if err != nil {
		return nil, err
	}

	return &Context{
		WorkDir:  workDir,
		Brewfile: brewfile,
		BrewPath: brewPath,
	}, nil
}
