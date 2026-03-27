package typst

import (
	"bytes"
	"os/exec"
	"strings"
)

// FontsCommand lists all discovered fonts in system and custom font paths
type FontsCommand struct {
	binPath     string
	fontOptions fontOptions
}

func NewFontsCommand(binPath string, opts ...FontOption) (*FontsCommand, error) {
	f := FontsCommand{
		binPath:     binPath,
		fontOptions: applyFontOptions(opts),
	}

	_, err := f.Help()

	if err != nil {
		return nil, err
	}

	return &f, nil
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

	cmd.Args = append(cmd.Args, c.fontOptions.args()...)

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
