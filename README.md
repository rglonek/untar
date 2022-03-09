# Simple untar

## Features:

* support for tar, tgz, tbz2
* simple UntarFile with src filename and dst destination
* support for Untar with io.Reader

## Usage (untar file):

This works on tar, tgz, tbz2 files. Discovery uses mimetype, so file extensions do not matter.

```go
err := untar.UntarFile("testsmall.tgz", "./")
if err != nil {
    log.Fatal(err)
}
```

## Full Usage Example with untar straight from web example:

```go
package main

import (
	"bytes"
	"log"
	"time"

	wget "github.com/rglonek/go-wget"
    "github.com/rglonek/untar"
)

func main() {
	timeout := time.Minute
	input := &wget.GetInput{
		Url:               "http://some/tgz/file.tgz",
		Auth: &wget.Auth{
			Username: "bob",
			Password: "test",
		},
		Timeout: &timeout,
	}
	output, err := wget.GetReader(input)
	if err != nil {
		log.Printf("Get Error: %s", err)
		return
	}
    defer output.R.Close()
    err = untar.Untar(output.R, "./")
	if err != nil {
		log.Printf("Untar Error: %s", err)
		return
	}
}
```
