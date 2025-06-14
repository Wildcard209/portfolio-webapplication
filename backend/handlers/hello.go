package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type HelloResponse struct {
	Message string `json:"message" example:"Hello from Go backend 2!"`
}

// Hello handles the hello endpoint
// @Summary Hello endpoint for testing
// @Description Returns a greeting message from the Go backend with hot reload support
// @Tags hello
// @Accept json
// @Produce json
// @Success 200 {object} HelloResponse
// @Router /test [get]
func Hello(c *gin.Context) {
	response := HelloResponse{
		Message: "Hello from Go backend 54353!",
	}
	c.JSON(http.StatusOK, response)
}
