package httpserver

import (
	"context"
	"net/http"
	"time"
)

const (
	defaultReadTimeout     = 5 * time.Second
	defaultWriteTimeout    = 10 * time.Second
	defaultAddr            = "0.0.0.0:8080"
	defaultShutdownTimeout = 5 * time.Second
)

type HttpServer struct {
	server          *http.Server
	notify          chan error
	shutdownTimeout time.Duration
}

func New(handler http.Handler, opts ...Option) *HttpServer {
	httpServer := &http.Server{
		Handler:      handler,
		ReadTimeout:  defaultReadTimeout,
		WriteTimeout: defaultWriteTimeout,
		Addr:         defaultAddr,
	}

	server := &HttpServer{
		server: httpServer,
		notify: make(chan error, 1),
	}

	for _, opt := range opts {
		opt(server)
	}

	server.start()

	return server
}

func (hs *HttpServer) start() {
	go func() {
		hs.notify <- hs.server.ListenAndServe()
		close(hs.notify)
	}()
}

func (hs *HttpServer) Notify() <-chan error {
	return hs.notify
}

func (hs *HttpServer) Shutdown() error {
	ctx, cancel := context.WithTimeout(context.Background(), hs.shutdownTimeout)
	defer cancel()

	return hs.server.Shutdown(ctx)
}
