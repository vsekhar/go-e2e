// +build js,wasm

// Package browser wraps syscall/js.
package browser

import (
	"syscall/js"
)

// Alert opens a modal dialog in the browser with message m.
func Alert(m string) {
	js.Global.Call("alert", m)
}

// On registers an event handler on an element
func On(event, element string, f func()) func() {
	cb := js.NewEventCallback(false, false, false, func(js.Value) {
		f()
	})
	js.Global.Get("document").Call("getElementById", element).Call("addEventListener", "click", js.ValueOf(cb))
	return cb.Close
}

// Set sets the value of property on element.
func Set(element, property string, value interface{}) {
	js.Global.Get("document").Call("getElementById", element).Set(property, js.ValueOf(value))
}

// Sometime arranges f to be called asynchronously.
func Sometime(f func()) {
	var cb js.Callback
	cb = js.NewCallback(func(_ []js.Value) {
		f()
		cb.Close()
	})
	js.ValueOf(cb).Invoke()
}

// ServeForever defers to the browser's event loop. The return value exists
// only to permit a call to ServeForever to be used as an argument to panic().
//
// ServeForever never returns.
func ServeForever() struct{} {
	select {}
}
