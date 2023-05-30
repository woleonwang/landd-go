package main

import (
	"github.com/gin-contrib/logger"
	"github.com/gin-contrib/requestid"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"landd.co/landd/pkg/handler"
	"landd.co/landd/pkg/middleware"
	"landd.co/landd/pkg/mysql"
	"landd.co/landd/pkg/utils"
)

func main() {
	mysql.Init()
	utils.Init()

	router := gin.Default()
	router.Use(requestid.New()).Use(logger.SetLogger()) // TODO add request id into log output for easy debugging
	router.Use(sessions.Sessions("session", cookie.NewStore([]byte("coePCZ7yNeBtWyWtCTbTw6ZOszyL3nYf"))))

	registerRoutes(router)
	router.Run(":8080")
}

func registerRoutes(router *gin.Engine) {
	loginHandler := handler.NewLoginHandler()
	router.POST("/sign_up", loginHandler.SignUp)
	router.POST("/login", loginHandler.Login)
	router.Use(middleware.AuthRequired).GET("/logout", loginHandler.Logout)
}
