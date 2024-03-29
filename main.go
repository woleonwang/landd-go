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
	recruiterRoutes := router.Group("/recruiter")
	{
		profileRoutes := recruiterRoutes.Group("/profile")
		{
			recruiterHandler := handler.NewRecruiterProfileHandler()
			profileRoutes.GET("/:user_id", recruiterHandler.GetProfileInfo)
			profileRoutes.Use(middleware.RecruiterAuth).POST("/", recruiterHandler.UpdateProfileInfo)
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

	partnerRoutes := router.Group("/partner")
	{
		profileHandler := handler.NewPartnerProfileHandler()
		partnerRoutes.Use(middleware.PartnerAuth).GET("/profile/:user_id", profileHandler.Get)
		partnerRoutes.Use(middleware.PartnerAuth).POST("/profile", profileHandler.Update)

		homepageHandler := handler.NewPartnerHomepageHandler()
		partnerRoutes.Use(middleware.PartnerAuth).GET("/homepage/:user_id", homepageHandler.Get)
		partnerRoutes.Use(middleware.PartnerAuth).POST("/homepage", homepageHandler.Update)

		ctpHandler := handler.NewPartnerCTPHandler()
		partnerRoutes.Use(middleware.PartnerAuth).GET("/ctp/:user_id", ctpHandler.Get)
		partnerRoutes.Use(middleware.PartnerAuth).POST("/ctp/new", ctpHandler.Create)
		partnerRoutes.Use(middleware.PartnerAuth).POST("/ctp", ctpHandler.Update)
	}

	adminRoutes := router.Group("/admin")
	{
		jobHandler := handler.NewJobHandler()
		adminRoutes.Use(middleware.AdminAuth).GET("/jobs", jobHandler.Get)
		adminRoutes.Use(middleware.AdminAuth).POST("/jobs/new", jobHandler.Create)
		adminRoutes.Use(middleware.AdminAuth).POST("/jobs", jobHandler.Update)
	}

	loginHandler := handler.NewLoginHandler()
	router.POST("/sign_up", loginHandler.SignUp)
	router.POST("/login", loginHandler.Login)
	router.Use(middleware.AuthRequired).GET("/logout", loginHandler.Logout)
	fileHandler := handler.NewFileHandler()
	router.Use(middleware.AuthRequired).POST("/file/upload", fileHandler.Upload)
}
