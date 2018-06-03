package console

import (
	"syscall/js"
)

var console = js.Global.Get("console")

// Info prints m to the console as an info message.
func Info(m string) {
	console.Call("info", m)
}
