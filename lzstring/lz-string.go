package lzstring

import (
	"fmt"
	"log"

	_ "embed"

	goja "github.com/golang-js/gojs"
)

// there's an encoding problem of go version LZString lib,
// so let's use the js version temporarily
// to build the project, you will need the lz-string.js
// https://www.npmjs.com/package/lz-string
// https://unpkg.com/lz-string@1.4.4/libs/lz-string.js
//go:embed lz-string.js
var JS []byte
var VM = goja.New()

func Encode(s string) string {
	js := string(JS) + fmt.Sprintf(`LZString.compressToBase64('%s')`, s)
	value, err := VM.RunString(js)
	if err != nil {
		log.Fatalln(err)
	}

	return value.String()
}

func Decode(s string) string {
	js := string(JS) + fmt.Sprintf(`LZString.decompressFromBase64('%s')`, s)
	value, err := VM.RunString(js)
	if err != nil {
		log.Fatalln(err)
	}

	return value.String()
}
