package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/suhlig/gosmos/typst"

	"golang.org/x/sync/errgroup"
)

func main() {
	err := mainE(context.Background(), os.Args[1:]...)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error %s\n", err)
		os.Exit(1)
	}
}

const typstVersion = "0.14.2"

func mainE(ctx context.Context, args ...string) error {
	flags := flag.NewFlagSet("fetch-typst", flag.ContinueOnError)
	path := flags.String("path", "", "base directory for downloaded typst binaries")

	if err := flags.Parse(args); err != nil {
		return err
	}

	fetcher := typst.NewFetcher(typstVersion)

	g, ctx := errgroup.WithContext(ctx)

	for _, platform := range typst.Platforms(typstVersion) {
		g.Go(func() error {
			path := filepath.FromSlash(filepath.Join(*path, platform.OS, platform.Arch, "typst"))

			fileExists, err := fileExists(path)

			if err != nil {
				return err
			}

			if fileExists {
				// We trust what's there and do not perform additional checks, e.g. a checksum. We are not paranoid.
				return nil
			}

			bin, err := fetcher.Fetch(ctx, platform)

			if err != nil {
				return err
			}

			err = os.MkdirAll(filepath.Dir(path), 0755)

			if err != nil {
				return err
			}

			err = os.WriteFile(path, bin, 0755)

			if err != nil {
				return fmt.Errorf("writing %s: %w", path, err)
			}

			return nil
		})
	}

	return g.Wait()
}

func fileExists(path string) (bool, error) {
	_, err := os.Stat(path)

	if err == nil {
		return true, nil
	}

	if errors.Is(err, os.ErrNotExist) {
		return false, nil
	}

	return false, err // permission or other error
}
