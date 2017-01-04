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
  -recurse string
        Glob for directories to recurse into (default "[^.]*")
$ tree
.
|-- css
|   `-- site.css
|-- hello.txt
|-- img
|   `-- example.png
`-- js
    `-- menus.js

3 directories, 4 files
$ cat hello.txt
world
$ scattered *.txt
{
        "hello.txt": "hello.591785b794601e212b260e25925636fd.txt"
}
$ ls -1 hello*
hello.591785b794601e212b260e25925636fd.txt
hello.txt
$ scattered '*.css' '*.js' '*.png'
{
        "css/site.css": "css/site.d41d8cd98f00b204e9800998ecf8427e.css",
        "img/example.png": "img/example.d41d8cd98f00b204e9800998ecf8427e.png",
        "js/menus.js": "js/menus.d41d8cd98f00b204e9800998ecf8427e.js"
}
```
