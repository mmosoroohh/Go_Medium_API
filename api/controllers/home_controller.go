package controllers

import (
	"github.com/mmosoroohh/Go_Medium_API/api/responses"
	"net/http"
)

func (server *Server) Home(w http.ResponseWriter, r *http.Request)  {
	responses.JSON(w, http.StatusOK, "Welcome to Medium API")
}