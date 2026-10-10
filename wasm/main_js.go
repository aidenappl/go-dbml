//go:build js && wasm

package main

import (
	"runtime/debug"
	"syscall/js"

	dbml "github.com/aidenappl/go-dbml"
)

// str returns args[i] as a Go string ("" when missing or not a string).
func str(args []js.Value, i int) string {
	if i < len(args) && args[i].Type() == js.TypeString {
		return args[i].String()
	}
	return ""
}

// fn wraps a handler; js.FuncOf callbacks must never panic.
func fn(h func(args []js.Value) string) js.Func {
	return js.FuncOf(func(_ js.Value, args []js.Value) (out any) {
		defer func() {
			if r := recover(); r != nil {
				out = `{"ok":false,"error":"internal error"}`
			}
		}()
		return h(args)
	})
}

func main() {
	// The linear memory never shrinks: collect earlier than the default
	// (heap target 1.5x the live heap instead of 2x).
	debug.SetGCPercent(50)
	api := js.Global().Get("Object").New()
	api.Set("version", dbml.Version)
	api.Set("parse", fn(func(a []js.Value) string { return apiParse(str(a, 0), str(a, 1)) }))
	api.Set("exportSQL", fn(func(a []js.Value) string { return apiExportSQL(str(a, 0), str(a, 1), str(a, 2)) }))
	api.Set("importSQL", fn(func(a []js.Value) string { return apiImportSQL(str(a, 0), str(a, 1), str(a, 2)) }))
	api.Set("emit", fn(func(a []js.Value) string { return apiEmit(str(a, 0)) }))
	js.Global().Set("goDbml", api)
	if cb := js.Global().Get("goDbmlReady"); cb.Type() == js.TypeFunction {
		cb.Invoke()
	}
	select {}
}
