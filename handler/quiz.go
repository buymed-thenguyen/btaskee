package handler

import (
	"btaskee/domain"
	"btaskee/model/constant"
	"github.com/gin-gonic/gin"
)

func GetQuizDetail(c *gin.Context) {
	sessionCode := c.Param("code")
	c.Set(constant.DATA_CTX, domain.GetQuizDetail(c, sessionCode))
}

func GetListQuiz(c *gin.Context) {
	c.Set(constant.DATA_CTX, domain.GetListQuiz(c))
}
