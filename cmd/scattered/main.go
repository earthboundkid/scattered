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

func getPaths(recpat string, globs []string) (paths []string, err error) {
	err = filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}

		finaldir := filepath.Base(filepath.Dir(path))
		if matched, err := filepath.Match(recpat, finaldir); err != nil {
			return err
		} else if !matched && finaldir != "." {
			return filepath.SkipDir
		}

		for _, glob := range globs {
			if matched, err := filepath.Match(glob, filepath.Base(path)); err != nil {
				return err
			} else if matched && !scattered.IsHashedPath(path) {
				paths = append(paths, path)
			}
		}
		return nil
	})

	return paths, err
}

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	dryrun := flag.Bool("dryrun", false, "Just create the JSON manifest; don't link files")
	recurse := flag.String("recurse", "[^.]*", "Glob for directories to recurse into")
	flag.Usage = func() {
		fmt.Fprint(os.Stderr, `Usage of scattered:

	scattered [options] <globs>...

Given a shell path or glob, for each file it makes an MD5 hash and
hard-links basename.HASH.ext to the file. Finally, it returns a JSON
object mapping input to output paths for use as a file manifest by
some other tool.

Options:

`)
		flag.PrintDefaults()
	}
	flag.Parse()
	paths, err := getPaths(*recurse, flag.Args())
	if err != nil {
		return err
	}

	var pathsMap = map[string]string{}

	for _, src := range paths {
		dst, err := scattered.HashPath(src)
		if err == scattered.ErrIsDir {
			continue
		}
		if err != nil {
			return err
		}
		pathsMap[src] = dst
	}

	if !*dryrun {
		if err = link(pathsMap); err != nil {
			return err
		}
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
