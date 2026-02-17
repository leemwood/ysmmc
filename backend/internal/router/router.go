package router

import (
	"github.com/gin-gonic/gin"
	"github.com/ysmmc/backend/internal/handler"
	"github.com/ysmmc/backend/internal/middleware"
)

func Setup(r *gin.Engine) {
	r.Use(middleware.CORS())
	r.Use(middleware.Logger())

	authHandler := handler.NewAuthHandler()
	userHandler := handler.NewUserHandler()
	modelHandler := handler.NewModelHandler()
	adminHandler := handler.NewAdminHandler()
	favoriteHandler := handler.NewFavoriteHandler()
	announcementHandler := handler.NewAnnouncementHandler()
	uploadHandler := handler.NewUploadHandler()

	api := r.Group("/api")
	{
		auth := api.Group("/auth")
		{
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)
			auth.POST("/refresh", authHandler.RefreshToken)
			auth.POST("/forgot-password", authHandler.ForgotPassword)
			auth.POST("/reset-password", authHandler.ResetPassword)
			auth.GET("/verify", authHandler.VerifyEmail)
			auth.POST("/change-email", middleware.Auth(), authHandler.ChangeEmail)
			auth.GET("/verify-email-change", authHandler.VerifyEmailChange)
			auth.GET("/me", middleware.Auth(), authHandler.Me)
		}

		users := api.Group("/users")
		{
			users.GET("/me", middleware.Auth(), userHandler.GetMe)
			users.PUT("/me", middleware.Auth(), userHandler.UpdateMe)
			users.PUT("/me/password", middleware.Auth(), userHandler.ChangePassword)
			users.GET("/:id", userHandler.GetByID)
			users.GET("/:id/models", userHandler.GetUserModels)
		}

		models := api.Group("/models")
		{
			models.GET("", modelHandler.List)
			models.GET("/:id", modelHandler.GetByID)
			models.POST("", middleware.Auth(), modelHandler.Create)
			models.PUT("/:id", middleware.Auth(), modelHandler.Update)
			models.DELETE("/:id", middleware.Auth(), modelHandler.Delete)
			models.POST("/:id/download", modelHandler.Download)
			models.POST("/:id/favorite", middleware.Auth(), modelHandler.AddFavorite)
			models.DELETE("/:id/favorite", middleware.Auth(), modelHandler.RemoveFavorite)
			models.GET("/:id/favorite", middleware.Auth(), modelHandler.CheckFavorite)
		}

		favorites := api.Group("/favorites")
		favorites.Use(middleware.Auth())
		{
			favorites.GET("", favoriteHandler.List)
		}

		announcements := api.Group("/announcements")
		{
			announcements.GET("", announcementHandler.List)
			announcements.GET("/all", announcementHandler.ListAll)
			announcements.GET("/:id", announcementHandler.GetByID)
		}

		upload := api.Group("/upload")
		upload.Use(middleware.Auth())
		{
			upload.POST("/model", uploadHandler.UploadModel)
			upload.POST("/image", uploadHandler.UploadImage)
		}

		admin := api.Group("/admin")
		admin.Use(middleware.Auth())
		admin.Use(middleware.AdminOnly())
		{
			admin.GET("/stats", adminHandler.GetStats)
			admin.GET("/super-admin", adminHandler.GetSuperAdmin)
			admin.GET("/models/pending", adminHandler.ListPendingModels)
			admin.GET("/models/pending-updates", adminHandler.ListPendingUpdates)
			admin.PUT("/models/:id/approve", adminHandler.ApproveModel)
			admin.PUT("/models/:id/reject", adminHandler.RejectModel)
			admin.GET("/users", adminHandler.ListUsers)
			admin.PUT("/users/:id/role", adminHandler.UpdateUserRole)
			admin.PUT("/users/:id/admin", middleware.SuperAdminOnly(), adminHandler.SetAdmin)
			admin.DELETE("/users/:id/admin", middleware.SuperAdminOnly(), adminHandler.RemoveAdmin)
			admin.PUT("/users/:id/ban", adminHandler.BanUser)
			admin.PUT("/users/:id/unban", adminHandler.UnbanUser)
			admin.GET("/profiles/pending", adminHandler.ListPendingProfiles)
			admin.PUT("/profiles/:id/approve", adminHandler.ApproveProfile)
			admin.PUT("/profiles/:id/reject", adminHandler.RejectProfile)
			admin.POST("/announcements", adminHandler.CreateAnnouncement)
			admin.PUT("/announcements/:id", adminHandler.UpdateAnnouncement)
			admin.DELETE("/announcements/:id", adminHandler.DeleteAnnouncement)
		}
	}

	r.Static("/uploads", "../uploads")

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})
}
