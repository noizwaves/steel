package impl

import (
	"bufio"
	"fmt"
	"os"
	"slices"
	"strings"
)

type Brewfile struct {
	lines []string
}

func LoadBrewfile(path string) (*Brewfile, error) {
	readFile, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer readFile.Close()

	scanner := bufio.NewScanner(readFile)
	scanner.Split(bufio.ScanLines)

	lines := make([]string, 0)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "brew ") {
			lines = append(lines, line)
		}
	}

	return &Brewfile{
		lines: lines,
	}, nil
}

func (b *Brewfile) IncludesPackage(name string) bool {
	return slices.Contains(b.lines, fmt.Sprintf("brew '%s'", name))
}
