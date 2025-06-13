package handler

import (
	"btaskee/domain"
	"btaskee/model/constant"
	reqModel "btaskee/model/request"
	"github.com/gin-gonic/gin"
	"net/http"
)

func Login(c *gin.Context) {
	var req *reqModel.User
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(err)
		c.Status(http.StatusBadRequest)
		return
	}

	c.Set(constant.DATA_CTX, domain.Login(c, req))
}

func Signup(c *gin.Context) {
	var req *reqModel.User
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(err)
		c.Status(http.StatusBadRequest)
		return
	}

	c.Set(constant.DATA_CTX, domain.Signup(c, req))
}

func GetMe(c *gin.Context) {
	c.Set(constant.DATA_CTX, domain.GetMe(c))
}
