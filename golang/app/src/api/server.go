package api

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"test/utils"
)

// Server provides an http.Server.
type Server struct {
	*http.Server
}

// NewServer creates and configures an APIServer serving all application routes.
func NewServer(dbConnection utils.DbConnection, test bool) (*Server, error) {
	log.Println("configuring server...")
	api, err := NewApp(dbConnection)
	if err != nil {
		return nil, err
	}

	var addr string
	var port string
	if !test {
		port = os.Getenv("PORT")
	} else {
		port = os.Getenv("TEST_PORT")
	}
	if strings.Contains(port, ":") {
		addr = port
	} else {
		addr = ":" + port
	}

	srv := http.Server{
		Addr:    addr,
		Handler: api,
	}

	return &Server{&srv}, nil
}

// Start runs ListenAndServe on the http.Server with graceful shutdown.
func (srv *Server) Start() {
	log.Println("starting server...")
	go func() {
		if err := srv.ListenAndServe(); err != http.ErrServerClosed {
			panic(err)
		}
	}()
	log.Printf("Listening on %s\n", srv.Addr)

	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	sig := <-quit
	log.Println("Shutting down server... Reason:", sig)
	// teardown logic...

	if err := srv.Shutdown(context.Background()); err != nil {
		panic(err)
	}
	log.Println("Server gracefully stopped")
}
