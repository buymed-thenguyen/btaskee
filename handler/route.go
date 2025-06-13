package handler

import (
	"btaskee/config"
	"btaskee/handler/middleware"
	"btaskee/ws"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func SetupRouter(authCfg *config.Auth) *gin.Engine {
	r := gin.New()
	r.Use(gin.Recovery())

	api := r.Group("/", middleware.ResponseWrapper())
	api.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Authorization", "Content-Type"},
		AllowCredentials: true,
	}))

	api.POST("/seed", SeedData)
	api.POST("/user/sign-up", Signup)
	api.POST("/user/log-in", Login)

	user := api.Group("/user", middleware.AuthMiddleware(authCfg))
	user.GET("/me", GetMe)

	session := api.Group("/session", middleware.AuthMiddleware(authCfg))
	session.POST("/create", CreateSessionWithQuizID)
	session.POST("/:code/join", JoinSessionByCode)
	session.POST("/:code/submit", SubmitAnswer)
	session.PUT("/:code/start", StartSession)
	session.GET("/:code/leaderboard", GetLeaderboardBySession)
	session.GET("/:code", GetSessionDetail)
	session.GET("/:code/quiz", GetQuizDetail)
	session.GET("/:code/participants", GetSessionParticipants)
	session.GET("/:code/participants/answers", GetSessionParticipantAnswers)

	quiz := api.Group("/quiz", middleware.AuthMiddleware(authCfg))
	quiz.GET("/", GetListQuiz)

	// ws
	api.GET("/ws/:code", ws.HandleWS)

	// html
	r.Static("/template", "./template")
	r.StaticFile("/", "./template/index.html")
	r.StaticFile("/log-out", "./template/logout.html")
	r.StaticFile("/quizzes", "./template/quizzes.html")
	r.StaticFile("/session", "./template/session.html")
	r.StaticFile("/question", "./template/question.html")
	r.StaticFile("/leaderboard", "./template/leaderboard.html")

	return r
}
