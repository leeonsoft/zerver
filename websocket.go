package zerver_rest

import (
	"net"
	"net/url"

	websocket "github.com/cosiner/zerver_websocket"
)

type (
	// webSocketConn represent an websocket connection
	// WebSocket connection is not be managed in server,
	// it's handler's responsibility to close connection
	WebSocketConn interface {
		URLVarIndexer
		net.Conn
		URL() *url.URL
	}

	// webSocketConn is the actual websocket connection
	webSocketConn struct {
		*websocket.Conn
		URLVarIndexer
	}

	// WebSocketHandlerFunc is the websocket connection handler
	WebSocketHandlerFunc func(WebSocketConn)

	// WebSocketHandler is the handler of websocket connection
	WebSocketHandler interface {
		Init(*Server) error
		Destroy()
		Handle(WebSocketConn)
	}
)

// WebSocketHandlerFunc is a function WebSocketHandler
func (WebSocketHandlerFunc) Init(*Server) error           { return nil }
func (fn WebSocketHandlerFunc) Handle(conn WebSocketConn) { fn(conn) }
func (WebSocketHandlerFunc) Destroy()                     {}

// newWebSocketConn wrap a exist websocket connection and url variables to a
// new webSocketConn
func newWebSocketConn(conn *websocket.Conn, varIndexer URLVarIndexer) *webSocketConn {
	return &webSocketConn{
		Conn:          conn,
		URLVarIndexer: varIndexer,
	}
}

// URL return client side url
func (wsc *webSocketConn) URL() *url.URL {
	return wsc.Config().Origin
}
