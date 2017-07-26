package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/carlmjohnson/scattered"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(2)
	}
}

func run() error {
	dryrun := flag.Bool("dryrun", false, "Just create the JSON manifest; don't create files")
	basepath := flag.String("basepath", ".", "Base directory to process from")
	dirpat := flag.String("dirpat", "^[^.].*", "Regex for directories to process files in")
	output := flag.String("output", "", "File to save manifest (stdout if unset)")
	link := flag.Bool("link", false, "Use hardlinks instead of copying files")

	flag.Usage = func() {
		fmt.Fprint(os.Stderr, `Usage of scattered:

	scattered [options] <globs>...

Given a shell path or glob, for each file it makes an MD5 hash and
copies the file to basename.HASH.ext. Finally, it returns a JSON
object mapping input to output paths for use as a file manifest by
some other tool.

Options:

`)
		flag.PrintDefaults()
	}
	flag.Parse()
	pathsMap, err := scattered.HashFileGlobs(*basepath, *dirpat, flag.Args()...)
	if err != nil {
		return err
	}

	if !*dryrun {
		fileaction := scattered.Copy
		if *link {
			fileaction = scattered.Link
		}

		if err = fileaction(*basepath, pathsMap); err != nil {
			return err
		}
	}

	b, err := json.MarshalIndent(&pathsMap, "", "\t")
	if err != nil {
		return err
	}

	var fout io.Writer = os.Stdout
	if *output != "" {
		f, err := os.Create(*output)
		if err != nil {
			return err
		}
		defer f.Close()
		fout = f
	}

	if _, err = fout.Write(b); err != nil {
		return err
	}

	// Trailing newline
	_, err = io.WriteString(fout, "\n")
	return err
}
