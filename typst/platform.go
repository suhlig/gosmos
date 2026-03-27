package typst

import (
	"fmt"
	"path/filepath"
	"runtime"
)

type Platform struct {
	OS   string
	Arch string
}

func (p Platform) String() string {
	return fmt.Sprintf("{%s, %s}", p.OS, p.Arch)
}

func Platforms(version string) []Platform {
	assets, ok := versions[version]

	if !ok {
		return nil
	}

	seen := make(map[Platform]bool)

	var out []Platform

	for _, a := range assets {
		if !seen[a.Platform] {
			out = append(out, a.Platform)
			seen[a.Platform] = true
		}
	}

	return out
}

// BinPath returns the path to the typst binary for the given platform under
// the given base directory, using the convention established by fetch-typst:
// <base>/<OS>/<Arch>/typst.
func BinPath(base string, p Platform) string {
	return filepath.Join(base, p.OS, p.Arch, "typst")
}

// CurrentBinPath returns the path to the typst binary for the current
// runtime platform (runtime.GOOS / runtime.GOARCH) under the given base
// directory.
func CurrentBinPath(base string) string {
	return BinPath(base, Platform{OS: runtime.GOOS, Arch: runtime.GOARCH})
}
