// +build js,wasm

package browser

import (
	"bufio"
	"io"
	"net"
	"syscall/js"
	"time"
)

// A Socket represents a connection to a remote host established by the browser.
type Socket struct {
	ws         js.Value
	r          io.Reader
	w          io.Writer
	ch         chan []byte
	connecting chan struct{}
	closed     chan struct{}
}

// Dial establishes a connection to a remote host with the specified protocol
// and returns the connection and an error.
//
// Currently only the WebSocket protocol ("ws") is supported.
func Dial(protocol, path string) (net.Conn, error) {
	if protocol != "ws" {
		panic("unsupported protocol: " + protocol)
	}
	s := &Socket{
		ws:         js.Global.Get("WebSocket").New(path),
		ch:         make(chan []byte, 1),
		connecting: make(chan struct{}),
		closed:     make(chan struct{}),
	}
	var w io.Writer
	s.r, w = io.Pipe()
	s.w = bufio.NewWriter(w)

	ocb := js.NewEventCallback(false, false, false, func(d js.Value) {
		close(s.connecting)
	})
	s.ws.Set("onopen", js.ValueOf(ocb))

	mcb := js.NewEventCallback(false, false, false, func(d js.Value) {
		data := BlobToString(d.Get("data"))
		println("message received: " + data)
		if _, err := s.w.Write([]byte(data)); err != nil {
			panic(err)
		}
	})
	s.ws.Set("onmessage", js.ValueOf(mcb))
	return s, nil
}

func (s *Socket) Close() error {
	s.ws.Call("close", 1000 /* normal */)
	close(s.closed)
	close(s.ch)
	return nil
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
	return s.r.Read(b)
}

func (s *Socket) waitUntilSent() {
	Every(50*time.Millisecond, func() bool {
		if s.ws.Get("bufferedAmount").Int() == 0 {
			return false
		}
		return true
	})
}

func (s *Socket) Write(b []byte) (int, error) {
	<-s.connecting
	s.waitUntilSent()
	s.ws.Call("send", b)
	s.waitUntilSent()
	return len(b), nil
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
