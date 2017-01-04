package scattered

import (
	"crypto/md5"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// HashReader reads the provided io.Reader and returns its MD5 hash as a string or an error.
func HashReader(r io.Reader) (string, error) {
	h := md5.New()
	_, err := io.Copy(h, r)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", h.Sum(nil)), nil
}

// ErrIsDir is returned by HashPath when attempting to hash a directory
var ErrIsDir = errors.New("Tried to hash the path of a directory")

// HashPath opens the file at the provided filepath and returns a
// string containing the file's hash as part of its filename
func HashPath(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()

	if stat, err := f.Stat(); err != nil {
		return "", err
	} else if stat.IsDir() {
		return "", ErrIsDir
	}

	hash, err := HashReader(f)
	if err != nil {
		return "", err
	}

	basename, ext := splitName(path)

	return fmt.Sprintf("%s.%s%s", basename, hash, ext), nil
}

// IsHashedPath returns true if a filepath matches the
// pattern path/name.HASH.ext, where HASH is an MD5 hash.
func IsHashedPath(path string) bool {
	basename, _ := splitName(path)
	innerExt := filepath.Ext(basename)
	return len(innerExt) == (md5.Size*2)+1
}

func splitName(path string) (basename, ext string) {
	ext = filepath.Ext(path)
	basename = path[:len(path)-len(ext)]
	return
}

// HashFileGlobs returns a map from filepaths to their HashPath equivalent for
// all files whoses parents match the dirpat glob and themselves match one
// of the filepat globs.
func HashFileGlobs(dirpat string, filepats ...string) (map[string]string, error) {
	paths, err := getPaths(dirpat, filepats...)
	if err != nil {
		return nil, err
	}

	var pathsMap = map[string]string{}

	for _, src := range paths {
		dst, err := HashPath(src)
		if err == ErrIsDir {
			continue
		}
		if err != nil {
			return nil, err
		}
		pathsMap[src] = dst
	}

	return pathsMap, nil
}

func getPaths(recpat string, globs ...string) (paths []string, err error) {
	err = filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		finaldir := filepath.Base(filepath.Dir(path))
		if matched, err := filepath.Match(recpat, finaldir); err != nil {
			return err
		} else if !matched && finaldir != "." {
			return filepath.SkipDir
		}

		if info.IsDir() {
			return nil
		}

		for _, glob := range globs {
			if matched, err := filepath.Match(glob, filepath.Base(path)); err != nil {
				return err
			} else if matched && !IsHashedPath(path) {
				paths = append(paths, path)
			}
		}
		return nil
	})

	return paths, err
}
