// Package transport is the broker's HTTP boundary: routing, the public
// error translation, and the domain-correlation logging helper.
package transport

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"time"
)

// readHeaderTimeout caps how long a client may take to send its request
// headers. Generous for localhost — it exists to reclaim connections from
// clients that died mid-request (runtime/resilience), not to fend off an
// attacker. The write and idle budgets stay unset until the WS-out exists
// and can inform them.
const readHeaderTimeout = 5 * time.Second

// Server wraps the broker's single HTTP boundary: the mux, its bound
// listener, and lifecycle control (Listen/Serve/Shutdown), kept separate so
// callers can bind — and learn the actual address — before Serve blocks.
type Server struct {
	httpServer *http.Server
	addr       string
	listener   net.Listener
}

// New builds a Server bound to addr, exposing GET /health and translating
// any unmatched route through the transport error boundary.
func New(addr string, storeStatus func() string) *Server {
	return &Server{
		httpServer: &http.Server{
			Handler:           NewMux(storeStatus),
			ReadHeaderTimeout: readHeaderTimeout,
		},
		addr: addr,
	}
}

// NewMux builds the broker's HTTP routes. Exported so tests can exercise
// routes directly via httptest without going through Listen/Serve.
func NewMux(storeStatus func() string) *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", HealthHandler(storeStatus))
	mux.HandleFunc("/", func(w http.ResponseWriter, _ *http.Request) {
		WriteError(w, http.StatusNotFound, ErrNotFound)
	})
	return mux
}

// Listen binds the configured address. Kept separate from Serve so callers
// know binding succeeded — and, for addr ":0", the actual assigned port —
// before anything is exposed to traffic.
func (s *Server) Listen() error {
	ln, err := net.Listen("tcp", s.addr)
	if err != nil {
		return fmt.Errorf("transport: listen: %w", err)
	}
	s.listener = ln
	return nil
}

// Addr returns the actual bound address; only meaningful after Listen
// succeeds.
func (s *Server) Addr() string {
	if s.listener == nil {
		return s.addr
	}
	return s.listener.Addr().String()
}

// Serve blocks accepting connections until Shutdown closes the listener.
func (s *Server) Serve() error {
	if err := s.httpServer.Serve(s.listener); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return fmt.Errorf("transport: serve: %w", err)
	}
	return nil
}

// Shutdown stops accepting new connections and waits for in-flight ones to
// finish (or ctx to expire) before returning.
func (s *Server) Shutdown(ctx context.Context) error {
	if err := s.httpServer.Shutdown(ctx); err != nil {
		return fmt.Errorf("transport: shutdown: %w", err)
	}
	return nil
}
