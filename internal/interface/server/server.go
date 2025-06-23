package server

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type Logger interface {
	Infof(format string, args ...interface{})
	Errorf(format string, args ...interface{})
}

type Server struct {
	engine *gin.Engine
	http   *http.Server
	logger Logger
}

func NewServer(engine *gin.Engine, logger Logger) *Server {
	return &Server{
		engine: engine,
		http: &http.Server{
			Handler:      engine,
			ReadTimeout:  30 * time.Second,
			WriteTimeout: 120 * time.Second,
			IdleTimeout:  180 * time.Second,
		},
		logger: logger,
	}
}

func (s *Server) Start(port string) error {
	s.http.Addr = fmt.Sprintf(":%s", port)
	if s.logger != nil {
		s.logger.Infof("Server started on port %s", port)
	}
	return s.http.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	if s.logger != nil {
		s.logger.Infof("Shutting down server...")
	}
	return s.http.Shutdown(ctx)
}
