// +build js,wasm

package browser

import (
	"sync"
	"syscall/js"
)

func BlobToString(b js.Value) string {
	fr := js.Global.Get("FileReader").New()
	l := new(sync.Mutex)
	l.Lock()
	lcb := js.NewEventCallback(false, false, false, func(js.Value) {
		println("reading done")
		l.Unlock()
	})
	defer lcb.Close()
	fr.Set("onload", js.ValueOf(lcb))
	fr.Call("readAsText", b)
	l.Lock()
	return fr.Get("result").String()
}
