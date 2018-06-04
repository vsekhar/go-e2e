package main

import (
	"fmt"
	"time"

	"github.com/vsekhar/go-e2e/lib/browser"
	"github.com/vsekhar/go-e2e/lib/browser/console"
)

func main() {
	console.Info("go main() started")

	// Run some code asynchronously.
	browser.Sometime(func() {
		println("browser.Sometime callback")
	})

	// Show modal alert when alert button is pressed.
	done := browser.On("click", "alert", func() {
		browser.Alert("alert button pressed")
	})
	defer done() // cleanup event handler when it is no longer needed.

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
	done = browser.On("click", "increment", func() {
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
	_, err = s.Read(b)
	if err != nil {
		panic(err)
	}
	println("message: " + string(b))

	panic(browser.ServeForever())
}
