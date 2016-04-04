# Scattered [![GoDoc](https://godoc.org/github.com/carlmjohnson/scattered?status.svg)](https://godoc.org/github.com/carlmjohnson/scattered)
Scattered is a command line tool for asset hashing. (It would be called “scattered, covered, and smothered,” but that name is too long.) It is useful as a stand-alone tool for hashing web assets. Given a shell path or glob, for each file it makes an MD5 hash and hard-links basename.HASH.ext to the file. Finally, it returns a JSON array of input and output paths for use as a file manifest by some other tool.

##Screenshots
```bash
$ cat hello.txt
world
$ scattered *.txt
[
        {
                "input": "hello.txt",
                "output": "hello.591785b794601e212b260e25925636fd.txt"
        }
]
$ ls -1
hello.591785b794601e212b260e25925636fd.txt
hello.txt
```
