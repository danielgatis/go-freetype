[![Go Report Card](https://goreportcard.com/badge/github.com/danielgatis/go-freetype?style=flat-square)](https://goreportcard.com/report/github.com/danielgatis/go-freetype)
[![License MIT](https://img.shields.io/badge/license-MIT-blue.svg)](https://raw.githubusercontent.com/danielgatis/go-freetype/master/LICENSE)
[![Go Doc](https://img.shields.io/badge/godoc-reference-blue.svg?style=flat-square)](https://godoc.org/github.com/danielgatis/go-freetype)

# go-freetype

Go bindings for the FreeType library. Only the high-level API is bound.

## Install

```bash
go get -u github.com/danielgatis/go-freetype
```

And then import the package in your code:

```go
import "github.com/danielgatis/go-freetype/freetype"
```

### Example

An example described below is one of the use cases.

```go
package main

import (
	"fmt"
	"image/png"
	"io/ioutil"
	"os"

	"github.com/danielgatis/go-findfont/findfont"
	"github.com/danielgatis/go-freetype/freetype"
)

func checkErr(err error) {
	if err != nil {
		fmt.Printf("Failed: %v\n", err)
		os.Exit(1)
	}
}

func main() {
	fonts, err := findfont.Find("Arial", findfont.FontRegular)
	checkErr(err)

	data, err := ioutil.ReadFile(fonts[0][2])
	checkErr(err)

	lib, err := freetype.NewLibrary()
	checkErr(err)

	face, err := freetype.NewFace(lib, data, 0)
	checkErr(err)

	err = face.Pt(32, 72)
	checkErr(err)

	img, m, err := face.Glyph('A')
	checkErr(err)

	f, err := os.Create("image.png")
	checkErr(err)

	err = png.Encode(f, img)
	checkErr(err)

	fmt.Printf("Image: %v\n", f.Name())
	fmt.Printf("Metrics: %v\n", m)

	err = face.Done()
	checkErr(err)

	err = lib.Done()
	checkErr(err)
}
```


```
❯ go run main.go
Image: image.png
Metrics: &{23 23 -1 23 21 -11 3 30}
```

image.png

![image.png](image.png)

## License

Copyright (c) 2020-present [Daniel Gatis](https://github.com/danielgatis)

Licensed under [MIT License](./LICENSE)
