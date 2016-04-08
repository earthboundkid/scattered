package scattered

import (
	"crypto/md5"
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

// HashPath opens the file at the provided filepath and returns a
// string containing the file's hash as part of its filename
func HashPath(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()

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
