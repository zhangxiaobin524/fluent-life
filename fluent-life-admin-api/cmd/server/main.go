package main

import (
	"log"

	"fluent-life-admin-api/internal/config"
	"fluent-life-admin-api/internal/handlers"
	"fluent-life-admin-api/internal/middleware"
	"fluent-life-admin-api/internal/models"
	"fluent-life-admin-api/pkg/response"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	db, err := config.InitDB(cfg)
	if err != nil {
		log.Fatalf("Failed to connect database: %v", err)
	}

	if err := models.AutoMigrate(db); err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	if cfg.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.New()
	r.Use(middleware.CORS())

	r.GET("/health", func(c *gin.Context) {
		response.Success(c, gin.H{"status": "ok"}, "服务运行正常")
	})

	adminHandler := handlers.NewAdminHandler(db)

	api := r.Group("/api/v1")
	{
		// 管理员登录
		api.POST("/admin/login", adminHandler.Login)

		// 测试根路由
		api.GET("/test-root", adminHandler.TestRoute)

		// 需要认证的管理接口（简化版，实际应该使用JWT中间件）
		admin := api.Group("/admin")
		// admin.Use(middleware.AdminAuthMiddleware())
		{
			// 用户管理
			admin.GET("/users", adminHandler.GetUsers)
			admin.GET("/users/:id", adminHandler.GetUser)
			admin.DELETE("/users/:id", adminHandler.DeleteUser)

			// 帖子管理
			admin.GET("/posts", adminHandler.GetPosts)
			admin.GET("/posts/:id", adminHandler.GetPost)
			admin.POST("/posts/delete-batch", adminHandler.DeletePost)

			// 房间管理
			admin.GET("/rooms", adminHandler.GetRooms)
			admin.GET("/rooms/:id", adminHandler.GetRoom)
			admin.DELETE("/rooms/:id", adminHandler.DeleteRoom)
			admin.POST("/rooms/delete-batch", adminHandler.DeleteRoom)
			admin.PATCH("/rooms/:id/toggle", adminHandler.ToggleRoom)

			// 训练统计
			admin.GET("/training/stats", adminHandler.GetTrainingStats)
			admin.GET("/training/records", adminHandler.GetTrainingRecords)

			// 绕口令管理
			admin.GET("/tongue-twisters", adminHandler.GetTongueTwisters)
			admin.GET("/tongue-twisters/:id", adminHandler.GetTongueTwister)
			admin.POST("/tongue-twisters", adminHandler.CreateTongueTwister)
			admin.POST("/tongue-twisters/batch-create", adminHandler.BatchCreateTongueTwisters)
			admin.PUT("/tongue-twisters/:id", adminHandler.UpdateTongueTwister)
			admin.POST("/tongue-twisters/delete-batch", adminHandler.DeleteTongueTwister)
			admin.DELETE("/tongue-twisters/all", adminHandler.DeleteAllTongueTwisters)
			admin.POST("/tongue-twisters/clean", adminHandler.CleanTongueTwisters)

			// 每日朗诵文案管理
			admin.GET("/daily-expressions", adminHandler.GetDailyExpressions)
			admin.GET("/daily-expressions/:id", adminHandler.GetDailyExpression)
			admin.POST("/daily-expressions", adminHandler.CreateDailyExpression)
			admin.POST("/daily-expressions/batch-create", adminHandler.BatchCreateDailyExpressions)
			admin.PUT("/daily-expressions/:id", adminHandler.UpdateDailyExpression)
			admin.POST("/daily-expressions/delete-batch", adminHandler.DeleteDailyExpression)

			// 测试路由
			admin.GET("/test", adminHandler.TestRoute)
		}
	}

	addr := ":" + cfg.Port
	log.Printf("Admin API Server starting on %s", addr)
	if err := r.Run(addr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
