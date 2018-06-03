package main

import (
	"fmt"
	"syscall/js"
	"time"

	"github.com/vsekhar/go-e2e/lib/browser"
)

func main() {
	window := browser.Window()
	console := browser.Console()
	console.Info("go main() started")

	// Run some code at some point.
	done := browser.Sometime(func(_ []js.Value) {
		println("browser.Sometime callback")
	})
	defer done()

	// Show modal alert when alert button is pressed.
	done = browser.OnClick("alert", func(e js.Value) {
		window.Alert("alert button pressed")
	})
	defer done()

	// Increment the clock.
	console.Info("starting clock goroutine")
	go func() {
		clock := 0
		for {
			time.Sleep(1 * time.Second)
			browser.Set("pings", "textContent", fmt.Sprintf("Clock goroutine has run for %d seconds", clock))
			clock++
		}
	}()
	console.Info("clock goroutine started")

	// Counter goroutine
	incCh := make(chan struct{})
	go func() {
		counter := 0
		for _ = range incCh {
			counter++
			browser.Set("counter", "textContent", fmt.Sprint(counter))
		}
	}()
	defer close(incCh)

	// Increment counter when increment button is pressed.
	done = browser.OnClick("increment", func(e js.Value) {
		incCh <- struct{}{}
	})
	defer done()

	// Open a socket
	s, err := browser.Dial("ws", "ws://html5rocks.websocket.org/echo")
	if err != nil {
		panic(err)
	}
	msg := []byte("hello")
	s.Write(msg)
	b := make([]byte, len(msg))
	s.Read(b)
	println(string(b))

	panic(browser.ServeForever())
}
