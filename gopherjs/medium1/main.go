package main

import (
	"crypto/sha256"

	"github.com/gopherjs/gopherjs/js"
)

//go:generate gopherjs build --minify

type UserInfo struct {
	*js.Object
	Name    string `js:"name"`
	Address string `js:"address"`
	City    string `js:"city"`
	State   string `js:"state"`
	Hash    []byte `js:"hash"`
}

func main() {
	js.Module.Get("exports").Set("hashit", hashit)
	js.Module.Get("exports").Set("hashobj", hashobj)
	js.Module.Get("exports").Set("hashobj2", hashobj2)
}

func hashit(s string) []byte {
	h := sha256.New()
	h.Write([]byte(s))
	return h.Sum(nil)
}

func hashobj(obj map[string]interface{}) []byte {
	h := sha256.New()
	for _, k := range []string{"name", "address", "city"} {
		if s, ok := obj[k]; ok {
			switch st := s.(type) {
			case string:
				h.Write([]byte(st))
			case []byte:
				h.Write(st)
			default:
			}
		}
	}
	return h.Sum(nil)
}

func hashobj2(obj *js.Object) *js.Object {
	ui := &UserInfo{Object: obj}

	h := sha256.New()
	h.Write([]byte(ui.Name))
	h.Write([]byte(ui.Address))
	h.Write([]byte(ui.City))
	ui.Hash = h.Sum(nil)
	return ui.Object
}
