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
	router.Static("/file", "./files")
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
	fileHandler := handler.NewFileHandler()
	router.Use(middleware.AuthRequired).POST("/file/upload", fileHandler.Upload)

	recruiterRoutes := router.Group("/recruiter")
	{
		profileRoutes := recruiterRoutes.Group("/profile").Use(middleware.RecruiterAuth)
		{
			recruiterHandler := handler.NewRecruiterProfileHandler()
			profileRoutes.GET("/:user_id", recruiterHandler.GetProfileInfo)
			profileRoutes.POST("/", recruiterHandler.UpdateProfileInfo)
		}
		endorseRoutes := recruiterRoutes.Group("/endorse")
		{
			endorseHandler := handler.NewEndorseHandler()
			endorseRoutes.GET("/draft/:user_id", endorseHandler.GetDraft)
			endorseRoutes.POST("/draft", endorseHandler.UpdateDraft)
			endorseRoutes.GET("/:user_id", endorseHandler.Get)
			endorseRoutes.GET("/:user_id/invite", endorseHandler.Invite)
			endorseRoutes.POST("/", endorseHandler.Update)
		}
	}

	partnerRoutes := router.Group("/partner").Use(middleware.PartnerAuth)
	{
		profileHandler := handler.NewPartnerProfileHandler()
		partnerRoutes.GET("/profile/:user_id", profileHandler.Get)
		partnerRoutes.POST("/profile", profileHandler.Update)

		homepageHandler := handler.NewPartnerHomepageHandler()
		partnerRoutes.GET("/homepage/:user_id", homepageHandler.Get)
		partnerRoutes.POST("/homepage", homepageHandler.Update)

		ctpHandler := handler.NewPartnerCTPHandler()
		partnerRoutes.GET("/ctp/:user_id", ctpHandler.Get)
		partnerRoutes.POST("/ctp/new", ctpHandler.Create)
		partnerRoutes.POST("/ctp", ctpHandler.Update)
	}

	adminRoutes := router.Group("/admin").Use(middleware.AdminAuth)
	{
		jobHandler := handler.NewJobHandler()
		adminRoutes.GET("/jobs", jobHandler.Get)
		adminRoutes.POST("/jobs/new", jobHandler.Create)
		adminRoutes.POST("/jobs", jobHandler.Update)
	}
}
