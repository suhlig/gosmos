package typst

import (
	"strings"
)

// DiagnosticLevel indicates the severity of a compiler diagnostic.
type DiagnosticLevel int

const (
	DiagnosticLevelInfo    DiagnosticLevel = iota
	DiagnosticLevelWarning DiagnosticLevel = iota
	DiagnosticLevelError   DiagnosticLevel = iota
)

func (l DiagnosticLevel) String() string {
	switch l {
	case DiagnosticLevelInfo:
		return "info"
	case DiagnosticLevelWarning:
		return "warning"
	case DiagnosticLevelError:
		return "error"
	default:
		return "unknown"
	}
}

// Diagnostic is a message emitted by the typst compiler.
type Diagnostic struct {
	Level   DiagnosticLevel
	Message string
}

// CompileResult holds the output of a successful compilation along with any
// diagnostics (informational messages, warnings, and errors) emitted during
// compilation.
type CompileResult struct {
	// PDF contains the compiled PDF bytes. It may be non-nil even when
	// Diagnostics contains warnings.
	PDF []byte

	// Diagnostics contains all messages emitted by the compiler. Warnings and
	// errors come from the compiler's stderr; informational messages come from
	// stdout.
	Diagnostics []Diagnostic
}

// parseStdout treats each non-empty line from the compiler's stdout as an
// informational diagnostic.
func parseStdout(output string) []Diagnostic {
	var diagnostics []Diagnostic

	for _, line := range strings.Split(output, "\n") {
		if strings.TrimSpace(line) == "" {
			continue
		}

		diagnostics = append(diagnostics, Diagnostic{
			Level:   DiagnosticLevelInfo,
			Message: line,
		})
	}

	return diagnostics
}

// parseStderr splits the compiler's stderr into diagnostic blocks. Each block
// starts with a line prefixed by "error:" or "warning:"; subsequent lines
// (indented context, source snippets) belong to the same block until the next
// blank line.
func parseStderr(output string) []Diagnostic {
	var diagnostics []Diagnostic

	// Blocks are separated by blank lines.
	blocks := strings.Split(output, "\n\n")

	for _, block := range blocks {
		block = strings.TrimSpace(block)
		if block == "" {
			continue
		}

		var level DiagnosticLevel

		switch {
		case strings.HasPrefix(block, "error:"):
			level = DiagnosticLevelError
		case strings.HasPrefix(block, "warning:"):
			level = DiagnosticLevelWarning
		default:
			// Unrecognised prefix — treat as informational.
			level = DiagnosticLevelInfo
		}

		diagnostics = append(diagnostics, Diagnostic{
			Level:   level,
			Message: block,
		})
	}

	return diagnostics
}
