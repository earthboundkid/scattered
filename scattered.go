package main

import (
	"crypto/md5"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

type result struct {
	Input    string `json:"input"`
	basename string
	ext      string
	Output   string `json:"output"`
}

func link(paths []result) (err error) {
	for _, path := range paths {
		_, err := os.Stat(path.Output)
		if !os.IsNotExist(err) {
			return err
		} else if err == nil {
			if err = os.Remove(path.Output); err != nil {
				return err
			}
		}

		if err = os.Link(path.Input, path.Output); err != nil {
			return err
		}
	}

	return nil
}

func makeHash(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}

	h := md5.New()
	_, err = io.Copy(h, f)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", h.Sum(nil)), nil
}

func splitName(path string) (basename, ext string) {
	ext = filepath.Ext(path)
	basename = path[:len(path)-len(ext)]
	return
}

func isHashedPath(basename string) bool {
	innerExt := filepath.Ext(basename)
	return len(innerExt) == (md5.Size*2)+1
}

func getPaths(globs []string) (paths []result, err error) {
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
			basename, ext := splitName(path)
			if !isHashedPath(basename) {
				paths = append(paths, result{path, basename, ext, ""})
			}
		}
	}
	return paths, err
}

func die(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func main() {
	flag.Parse()
	paths, err := getPaths(flag.Args())
	die(err)

	for i := range paths {
		hash, err := makeHash(paths[i].Input)
		die(err)
		paths[i].Output = fmt.Sprintf("%s.%s%s", paths[i].basename, hash, paths[i].ext)
	}

	die(link(paths))

	b, err := json.MarshalIndent(&paths, "", "\t")
	die(err)

	os.Stdout.Write(b)
	// Trailing newline
	os.Stdout.WriteString("\n")
}
