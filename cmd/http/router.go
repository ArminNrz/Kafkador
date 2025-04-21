package http

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"kafkador/cmd/service"
	util "kafkador/cmd/utils"
)

type Server struct {
	Config  util.Config
	Service service.BikerService
}

func NewServer(config util.Config, loop service.BikerService) (*Server, error) {
	server := &Server{
		Config:  config,
		Service: loop,
	}

	server.setupRouter()
	return server, nil
}

func (s *Server) setupRouter() {
	r := gin.Default()

	r.POST("/create", s.createBikerHandler)

	r.Run(fmt.Sprintf(":%s", s.Config.HTTPServerAddress))
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
