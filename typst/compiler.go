package typst

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
)

type CompileCommand struct {
	binPath     string
	fontOptions fontOptions
}

func NewCompileCommand(binPath string, opts ...FontOption) (*CompileCommand, error) {
	t := CompileCommand{
		binPath:     binPath,
		fontOptions: applyFontOptions(opts),
	}

	_, err := t.Help()

	if err != nil {
		return nil, err
	}

	return &t, nil
}

// TODO TYPST_ROOT
// TODO TYPST_PACKAGE_PATH
// TODO TYPST_PACKAGE_CACHE_PATH
// TODO SOURCE_DATE_EPOCH
// TODO TYPST_FEATURES

func (c CompileCommand) Help() (string, error) {
	cmd := exec.Command(c.binPath, "help")

	var stdout bytes.Buffer
	cmd.Stdout = &stdout

	err := cmd.Run()

	if err != nil {
		return "", err
	}

	return stdout.String(), nil
}

// Run compiles the given typst source and returns a CompileResult containing
// the PDF bytes and any diagnostics emitted by the compiler.
func (c CompileCommand) Run(input string) (*CompileResult, error) {
	pdfPath, err := os.CreateTemp("", "out.pdf")

	if err != nil {
		return nil, fmt.Errorf("unable to create temp file: %w", err)
	}

	defer os.Remove(pdfPath.Name())

	cmd := exec.Command(
		c.binPath,
		"compile",
		"-",
		"--format", "pdf",
		pdfPath.Name(),
	)

	cmd.Args = append(cmd.Args, c.fontOptions.args()...)

	stdin, err := cmd.StdinPipe()

	if err != nil {
		return nil, fmt.Errorf("unable to get stdin pipe: %s", err.Error())
	}

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	go func() {
		stdin.Write([]byte(input))
		stdin.Close()
	}()

	err = cmd.Run()

	diagnostics := append(parseStdout(stdout.String()), parseStderr(stderr.String())...)

	if err != nil {
		return &CompileResult{Diagnostics: diagnostics}, fmt.Errorf("failed to run command %s: %w", cmd, err)
	}

	pdf, err := os.ReadFile(pdfPath.Name())

	if err != nil {
		return nil, fmt.Errorf("unable to read result PDF: %w", err)
	}

	return &CompileResult{PDF: pdf, Diagnostics: diagnostics}, nil
}
