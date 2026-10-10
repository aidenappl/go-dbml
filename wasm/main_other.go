//go:build !(js && wasm)

package main

import (
	"fmt"
	"os"
)

func main() {
	fmt.Fprintln(os.Stderr, "go-dbml/wasm: build with GOOS=js GOARCH=wasm (see Devfile.yaml `dev wasm`)")
	os.Exit(2)
}
