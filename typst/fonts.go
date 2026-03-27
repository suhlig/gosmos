package typst

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// FontsCommand lists all discovered fonts in system and custom font paths
type FontsCommand struct {
	binPath               string
	fontPaths             []string
	isIgnoreSystemFonts   bool
	isIgnoreEmbeddedFonts bool
}

func NewFontsCommand(binPath string) (*FontsCommand, error) {
	f := FontsCommand{binPath: binPath}

	_, err := f.Help()

	if err != nil {
		return nil, err
	}

	fmt.Fprintf(os.Stderr, "returning %v\n", &f)

	return &f, nil
}

func (c *FontsCommand) WithFontPath(p string) *FontsCommand {
	c.fontPaths = append(c.fontPaths, p)
	return c
}

func (c *FontsCommand) WithoutSystemFonts() *FontsCommand {
	c.isIgnoreSystemFonts = true
	return c
}

func (c *FontsCommand) WithoutEmbeddedFonts() *FontsCommand {
	c.isIgnoreEmbeddedFonts = true
	return c
}

// TODO --variants Also lists style variants of each font family

func (c FontsCommand) Help() (string, error) {
	cmd := exec.Command(c.binPath, "fonts", "--help")

	var stdout bytes.Buffer
	cmd.Stdout = &stdout

	err := cmd.Run()

	if err != nil {
		return "", err
	}

	return stdout.String(), nil
}

func (c FontsCommand) Run() ([]string, error) {
	cmd := exec.Command(c.binPath, "fonts")

	if c.isIgnoreSystemFonts {
		cmd.Args = append(cmd.Args, "--ignore-system-fonts")
	}

	if c.isIgnoreEmbeddedFonts {
		cmd.Args = append(cmd.Args, "--ignore-embedded-fonts")
	}

	if len(c.fontPaths) > 0 {
		cmd.Args = append(cmd.Args, "--font-path", strings.Join(c.fontPaths, ":"))
	}

	var stdout bytes.Buffer
	cmd.Stdout = &stdout

	err := cmd.Run()

	if err != nil {
		return nil, err
	}

	fonts := strings.FieldsFunc(stdout.String(), func(r rune) bool {
		return r == '\n'
	})

	return fonts, nil
}
