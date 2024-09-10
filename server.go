package main

import (
	"context"
	"fmt"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"io"
	"net"
	"net/http"
)

type EchoHandler struct {
	log *zap.Logger
}

type HelloHandler struct {
	log *zap.Logger
}
type Route interface {
	http.Handler
	Pattern() string
}

func NewHTTPServer(lc fx.Lifecycle, mux *http.ServeMux, log *zap.Logger) *http.Server {
	srv := &http.Server{Addr: ":3000", Handler: mux}
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			listener, err := net.Listen("tcp", srv.Addr)
			if err != nil {
				return err
			}
			log.Info("...starting the server", zap.String("addr", srv.Addr))
			go srv.Serve(listener)
			return nil
		},
		OnStop: func(ctx context.Context) error {
			return srv.Shutdown(ctx)
		},
	})

	return srv
}

func NewEchoHandler(log *zap.Logger) *EchoHandler {
	return &EchoHandler{log}
}

func NewHelloHandler(log *zap.Logger) *HelloHandler {
	return &HelloHandler{log}
}

// example handler
func (e *EchoHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	if _, err := io.Copy(w, r.Body); err != nil {
		e.log.Warn("unable to handle the request", zap.Error(err))
	}
}
func (h *HelloHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		h.log.Error("Failed to read request", zap.Error(err))
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	if _, err := fmt.Fprintf(w, "Hello, %s\n", body); err != nil {
		h.log.Error("Failed to write response", zap.Error(err))
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}
func (*HelloHandler) Pattern() string {
	return "/hello"
}

func (*EchoHandler) Pattern() string {

	return "/echo"
}

// handle routing
func NewServeMux(route1, route2 Route) *http.ServeMux {
	router := http.NewServeMux()

	router.Handle(route1.Pattern(), route1)
	router.Handle(route2.Pattern(), route2)
	return router

}
