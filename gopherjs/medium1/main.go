package main

import (
	"crypto/sha256"

	"github.com/gopherjs/gopherjs/js"
)

//go:generate gopherjs build --minify

func main() {
	js.Module.Get("exports").Set("hashit", hashit)
}

func hashit(s string) []byte {
	h := sha256.New()
	h.Write([]byte(s))
	return h.Sum(nil)
}
