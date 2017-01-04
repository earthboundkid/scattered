package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/carlmjohnson/scattered"
)

func link(paths map[string]string) (err error) {
	for src, dst := range paths {
		_, err := os.Stat(dst)
		if !os.IsNotExist(err) {
			return err
		} else if err == nil {
			if err = os.Remove(dst); err != nil {
				return err
			}
		}

		if err = os.Link(src, dst); err != nil {
			return err
		}
	}

	return nil
}

func getPaths(globs []string) (paths []string, err error) {
	var fileset = map[string]bool{}

	for _, glob := range globs {
		globpaths, err := filepath.Glob(glob)
		if err != nil {
			return nil, err
		}

		for _, path := range globpaths {
			if seen := fileset[path]; seen {
				continue
			}

			fileset[path] = true
			if !scattered.IsHashedPath(path) {
				paths = append(paths, path)
			}
		}
	}

	return paths, err
}

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	flag.Parse()
	paths, err := getPaths(flag.Args())
	if err != nil {
		return err
	}

	var pathsMap = map[string]string{}

	for _, src := range paths {
		dst, err := scattered.HashPath(src)
		if err != nil {
			return err
		}
		pathsMap[src] = dst
	}

	if err = link(pathsMap); err != nil {
		return err
	}

	b, err := json.MarshalIndent(&pathsMap, "", "\t")
	if err != nil {
		return err
	}

	os.Stdout.Write(b)
	// Trailing newline
	_, err = os.Stdout.WriteString("\n")
	return err
}
