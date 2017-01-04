# Scattered [![GoDoc](https://godoc.org/github.com/carlmjohnson/scattered?status.svg)](https://godoc.org/github.com/carlmjohnson/scattered)
Scattered is a command line tool for asset hashing. (It would be called [“scattered, covered, and smothered,”][waho] but that name is too long.) It is useful as a stand-alone tool for hashing web assets. Given a shell path or glob, for each file it makes an MD5 hash and hard-links basename.HASH.ext to the file. Finally, it returns a JSON object mapping input to output paths for use as a file manifest by some other tool.

[waho]: https://en.wikipedia.org/wiki/Waffle_House

##Screenshots
```bash
$ scattered -h
Usage of scattered:

        scattered [options] <globs>...

Given a shell path or glob, for each file it makes an MD5 hash and
hard-links basename.HASH.ext to the file. Finally, it returns a JSON
object mapping input to output paths for use as a file manifest by
some other tool.

Options:

  -dryrun
        Just create the JSON manifest; don't link files
$ cat hello.txt
world
$ scattered *.txt
{
        "hello.txt": "hello.591785b794601e212b260e25925636fd.txt"
}
$ ls -1
hello.591785b794601e212b260e25925636fd.txt
hello.txt
```
