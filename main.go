package main

import (
	"encoding/gob"
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
	router.Static("/image", "./images")
	router.Use(requestid.New()).Use(logger.SetLogger()) // TODO add request id into log output for easy debugging
	gob.Register(mysql.User{})
	store := cookie.NewStore([]byte("coePCZ7yNeBtWyWtCTbTw6ZOszyL3nYf"))
	store.Options(sessions.Options{MaxAge: 7200})
	router.Use(sessions.Sessions("session", store))

	registerRoutes(router)
	router.Run(":8080")
}

func registerRoutes(router *gin.Engine) {
	loginHandler := handler.NewLoginHandler()
	router.POST("/sign_up", loginHandler.SignUp)
	router.POST("/login", loginHandler.Login)
	router.Use(middleware.AuthRequired).GET("/logout", loginHandler.Logout)

	recruiterRoutes := router.Group("/recruiter")
	{
		profileRoutes := recruiterRoutes.Group("/profile")
		{
			recruiterHandler := handler.NewRecruiterProfileHandler()
			profileRoutes.Use(middleware.AuthRequired).GET("/:user_id", recruiterHandler.GetProfileInfo)
			profileRoutes.Use(middleware.AuthRequired).POST("/", recruiterHandler.UpdateProfileInfo)
			profileRoutes.Use(middleware.AuthRequired).POST("/photo", recruiterHandler.UploadPhoto)
		}
		endorseRoutes := recruiterRoutes.Group("/endorse")
		{
			endorseHandler := handler.NewEndorseHandler()
			endorseRoutes.Use(middleware.AuthRequired).GET("/draft/:user_id", endorseHandler.GetDraft)
			endorseRoutes.Use(middleware.AuthRequired).POST("/draft", endorseHandler.UpdateDraft)
			endorseRoutes.Use(middleware.AuthRequired).GET("/:user_id", endorseHandler.Get)
			endorseRoutes.Use(middleware.AuthRequired).GET("/:user_id/invite", endorseHandler.Invite)
			endorseRoutes.Use(middleware.AuthRequired).POST("/", endorseHandler.Update)
		}
	}
}
