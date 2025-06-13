package domain

import (
	"btaskee/db"
	"btaskee/model/response"
	"github.com/gin-gonic/gin"
	"net/http"
)

func SeedData(c *gin.Context) *response.DefaultResponse {
	if err := db.SeedQuizzes(); err != nil {
		c.Error(err)
		c.Status(http.StatusInternalServerError)
	}
	return &response.DefaultResponse{Message: "ok"}
}
