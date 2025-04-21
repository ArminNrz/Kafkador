package http

import (
	"github.com/gin-gonic/gin"
	"kafkador/cmd/api"
	"net/http"
)

func (s *Server) createBikerHandler(ctx *gin.Context) {
	var req api.CreateBikerRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	res, err := s.Service.CreateBiker(ctx, req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"data": res})
}
