package server

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/rs/cors"

	"github.com/Kardainn/accountability/backend/config"
)

// Server is a server instance
type Server struct {
	ctx    context.Context
	server *http.Server
}

// Create creates a new Server instance
func Create(ctx context.Context) *Server {
	httpServer := &http.Server{
		Addr:         config.FromContext(ctx).HTTPListenAddress,
		ReadTimeout:  2 * time.Second,
		WriteTimeout: 2 * time.Minute,
	}
	server := &Server{
		ctx:    ctx,
		server: httpServer,
	}
	c := cors.New(cors.Options{
		AllowCredentials:   true,
		AllowedHeaders:     []string{"content-type"},
		AllowedMethods:     []string{"GET", "POST", "OPTIONS"},
		AllowedOrigins:     []string{"http://localhost:3000", "http://127.0.0.1:3000"},
		OptionsPassthrough: true,
	})
	r := mux.NewRouter()
	r.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method == http.MethodOptions {
				w.WriteHeader(200)
				return
			}
			next.ServeHTTP(w, r)
		})
	})
	r.HandleFunc("/", IsUp)
	r.HandleFunc("/createUser", server.userCreation).Methods("POST", "OPTIONS")
	// r.HandleFunc("/auth", server.auth).Methods("POST", "OPTIONS")
	handler := c.Handler(r)
	httpServer.Handler = handler
	return server
}

// Start start a server instance
func (s *Server) Start(signalErr chan error) {
	go func(signalErr chan error) {
		err := s.server.ListenAndServe()
		if err != nil {
			signalErr <- err
		}
		signalErr <- fmt.Errorf("server stopped")
	}(signalErr)
}

// IsUp is for test purpose
func IsUp(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)
	_, err := w.Write([]byte("the server is up and running"))
	if err != nil {
		w.WriteHeader(500)
	}
}
