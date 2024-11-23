package httpserver

import (
	"time"
)

type Option func(*HttpServer)

func ReadTimeout(timeout time.Duration) Option {
	return func(hs *HttpServer) {
		hs.server.ReadTimeout = timeout
	}
}

func WriteTimeout(timeout time.Duration) Option {
	return func(hs *HttpServer) {
		hs.server.WriteTimeout = timeout
	}
}

func Address(address string) Option {
	return func(hs *HttpServer) {
		hs.server.Addr = address
	}
}

func ShutdownTimeout(timeout time.Duration) Option {
	return func(hs *HttpServer) {
		hs.shutdownTimeout = timeout
	}
}
