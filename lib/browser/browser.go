// +build js,wasm

// Package browser wraps syscall/js.
package browser

import (
	"syscall/js"
)

type window struct {
}

func Window() window {
	return window{}
}

// Alert opens a modal dialog in the browser with message m.
func (window) Alert(m string) {
	js.Global.Call("alert", m)
}

type console struct {
	c js.Value
}

func Console() console {
	return console{c: js.Global.Get("console")}
}

func (c *console) Info(m string) {
	c.c.Call("info", m)
}

// OnClick registers a function to be called when an element is clicked.
func OnClick(id string, f func(js.Value)) func() {
	cb := js.NewEventCallback(false, false, false, f)
	js.Global.Get("document").Call("getElementById", id).Call("addEventListener", "click", js.ValueOf(cb))
	return cb.Close
}

// Set sets the value of property on element.
func Set(element, property string, value interface{}) {
	js.Global.Get("document").Call("getElementById", element).Set(property, js.ValueOf(value))
}

// Sometime arranges f to be called asynchronously in the future.
func Sometime(f func([]js.Value)) func() {
	cb := js.NewCallback(f)
	js.ValueOf(cb).Invoke()
	return cb.Close
}

// ServeForever defers to the browser's event loop. The return value exists
// only to permit a call to ServeForever to be used as an argument to panic().
//
// ServeForever never returns.
func ServeForever() struct{} {
	select {}
}
