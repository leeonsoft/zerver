package zerver

import (
	"io"
	"time"

	websocket "github.com/cosiner/zerver_websocket"
)

type (
	// WebSocketConn represent an websocket connection
	// WebSocket connection is not be managed in server,
	// it's handler's responsibility to close connection
	WebSocketConn interface {
		URLVarIndexer
		io.ReadWriteCloser
		SetDeadline(t time.Time) error
		SetReadDeadline(t time.Time) error
		SetWriteDeadline(t time.Time) error
		RemoteAddr() string
		serverGetter
	}

	// webSocketConn is the actual websocket connection
	webSocketConn struct {
		serverGetter
		*websocket.Conn
		URLVarIndexer
	}

	// WebSocketHandlerFunc is the websocket connection handler
	WebSocketHandlerFunc func(WebSocketConn)

	// WebSocketHandler is the handler of websocket connection
	WebSocketHandler interface {
		ServerInitializer
		Destroy()
		Handle(WebSocketConn)
	}
)

// newWebSocketConn wrap a exist websocket connection and url variables to a
// new webSocketConn
func newWebSocketConn(s serverGetter, conn *websocket.Conn, varIndexer URLVarIndexer) *webSocketConn {
	return &webSocketConn{
		serverGetter:  s,
		Conn:          conn,
		URLVarIndexer: varIndexer,
	}
}

func (wsc *webSocketConn) RemoteAddr() string {
	return wsc.Request().RemoteAddr
}

// WebSocketHandlerFunc is a function WebSocketHandler
func (WebSocketHandlerFunc) Init(*Server) error           { return nil }
func (fn WebSocketHandlerFunc) Handle(conn WebSocketConn) { fn(conn) }
func (WebSocketHandlerFunc) Destroy()                     {}
