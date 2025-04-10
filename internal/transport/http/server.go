package http 

import (
	"github.com/gin-gonic/gin"
	"context"
	"net/http"
)

type Server struct {
	server *gin.Engine
	port string 
} 

func NewServer(port string,  router *gin.Engine) *Server {
	return &Server{
		server: router,
		port: port,
	}
}

func (s *Server) Run() error {
	return s.server.Run(s.port)
}

func (s *Server) Shutdown(ctx context.Context) error {
	httpServer := &http.Server{
		Addr: s.port,
		Handler: s.server,
	}
	return httpServer.Shutdown(ctx)
}

