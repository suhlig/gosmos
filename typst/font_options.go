package typst

import "strings"

// fontOptions holds the font-related flags shared by all typst subcommands.
type fontOptions struct {
	fontPaths             []string
	isIgnoreSystemFonts   bool
	isIgnoreEmbeddedFonts bool
}

// FontOption is a functional option for configuring font-related flags.
type FontOption func(*fontOptions)

// WithFontPath adds a font search path.
func WithFontPath(p string) FontOption {
	return func(o *fontOptions) {
		o.fontPaths = append(o.fontPaths, p)
	}
}

// WithoutSystemFonts disables system font discovery.
func WithoutSystemFonts() FontOption {
	return func(o *fontOptions) {
		o.isIgnoreSystemFonts = true
	}
}

// WithoutEmbeddedFonts disables embedded font discovery.
func WithoutEmbeddedFonts() FontOption {
	return func(o *fontOptions) {
		o.isIgnoreEmbeddedFonts = true
	}
}

// args returns the CLI flags for the font options.
func (o fontOptions) args() []string {
	var args []string

	if o.isIgnoreSystemFonts {
		args = append(args, "--ignore-system-fonts")
	}

	if o.isIgnoreEmbeddedFonts {
		args = append(args, "--ignore-embedded-fonts")
	}

	if len(o.fontPaths) > 0 {
		args = append(args, "--font-path", strings.Join(o.fontPaths, ":"))
	}

	return args
}

func applyFontOptions(opts []FontOption) fontOptions {
	var o fontOptions

	for _, opt := range opts {
		opt(&o)
	}

	return o
}
