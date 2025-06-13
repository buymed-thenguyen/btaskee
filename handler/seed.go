package handler

import (
	"btaskee/domain"
	"btaskee/model/constant"
	"github.com/gin-gonic/gin"
)

func SeedData(c *gin.Context) {
	c.Set(constant.DATA_CTX, domain.SeedData(c))
}
