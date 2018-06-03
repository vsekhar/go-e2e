// +build js,wasm

package browser

import (
	"bytes"
	"errors"
	"net"
	"sync"
	"syscall/js"
	"time"
)

// A Socket represents a connection to a remote host established by the browser.
type Socket struct {
	ws         js.Value
	buf        bytes.Buffer
	err        error
	l          *sync.Mutex
	c          *sync.Cond
	connecting chan struct{}
}

// Dial establishes a connection to a remote host with the specified protocol
// and returns the connection and an error.
//
// Currently only the WebSocket protocol ("ws") is supported.
func Dial(protocol, path string) (net.Conn, error) {
	if protocol != "ws" {
		panic("unsupported protocol: " + protocol)
	}
	var s *Socket
	s = &Socket{
		ws:         js.Global.Get("WebSocket").New(path),
		connecting: make(chan struct{}),
	}
	s.l = new(sync.Mutex)
	s.c = sync.NewCond(s.l)
	ocb := js.NewEventCallback(false, false, false, func(d js.Value) {
		close(s.connecting)
	})
	s.ws.Set("onopen", js.ValueOf(ocb))

	// Message EventListener
	mcb := js.NewEventCallback(false, false, false, func(d js.Value) {
		fr := js.Global.Get("FileReader").New()
		lcb := js.NewEventCallback(false, false, false, func(js.Value) {
			// loading done, write to buffered pipe
			s.l.Lock()
			defer s.l.Unlock()
			data := fr.Get("result").String()
			if _, err := s.buf.Write([]byte(data)); err != nil {
				panic(err)
			}
			s.c.Signal()
		})
		fr.Set("onload", js.ValueOf(lcb))
		fr.Call("readAsText", d.Get("data"))
	})
	s.ws.Set("onmessage", js.ValueOf(mcb))

	// Error EventListener
	ecb := js.NewEventCallback(false, false, false, func(d js.Value) {
		s.l.Lock()
		defer s.l.Unlock()
		if s.err == nil {
			s.err = errors.New(d.String())
		}
	})
	s.ws.Set("onerror", js.ValueOf(ecb))

	return s, nil
}

func (s *Socket) Close() error {
	s.ws.Call("close", 1000 /* normal */)
	s.l.Lock()
	defer s.l.Unlock()
	return s.err
}

type socketAddr struct {
	network string
	addr    string
}

func (sa socketAddr) Network() string {
	return sa.network
}

func (sa socketAddr) String() string {
	return sa.addr
}

func (s *Socket) LocalAddr() net.Addr {
	return socketAddr{
		network: "ws",
		addr:    "browser",
	}
}

func (s *Socket) RemoteAddr() net.Addr {
	return socketAddr{
		network: "ws",
		addr:    s.ws.Get("url").String(),
	}
}

func (s *Socket) Read(b []byte) (int, error) {
	<-s.connecting
	s.l.Lock()
	defer s.l.Unlock()
	for s.buf.Len() == 0 {
		s.c.Wait()
	}
	n, err := s.buf.Read(b)
	if err != nil {
		return n, err
	}
	return n, s.err
}

func (s *Socket) Write(b []byte) (int, error) {
	<-s.connecting
	s.ws.Call("send", b)
	Every(50*time.Millisecond, func() bool {
		if s.ws.Get("bufferedAmount").Int() == 0 {
			return false
		}
		return true
	})
	s.l.Lock()
	defer s.l.Unlock()
	return len(b), s.err
}

func (*Socket) SetDeadline(time.Time) error {
	panic("unimplemented")
}

func (*Socket) SetReadDeadline(time.Time) error {
	panic("unimplemented")
}

func (*Socket) SetWriteDeadline(time.Time) error {
	panic("unimplemented")
}
