// +build js

package main

import (
	"github.com/gopherjs/gopherjs/js"
)

// Naive workaround for flag.Parse in browser of flag.Parse.
func init() {
	flags := js.Global.Get("document").Call("getElementById", "flags")
	if flags == nil {
		return
	}

	addrFromFlag := flags.Call("getAttribute", "data-addr").String()

	if addrFromFlag != "" {
		addr = &addrFromFlag
	}
}
