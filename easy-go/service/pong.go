package service

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func HandlePing(context *gin.Context) {
	context.String(http.StatusOK, "pong")
}
