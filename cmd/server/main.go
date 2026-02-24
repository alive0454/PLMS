package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"PLMS/internal/config"
	"PLMS/internal/database"
	"PLMS/internal/handlers"
	"PLMS/internal/middleware"
)

func main() {
	// 加载配置
	cfg := config.LoadConfig()

	// 初始化数据库连接
	db, err := database.InitDB(cfg.Database)
	if err != nil {
		log.Fatal("数据库连接失败:", err)
	}

	// 设置 Gin 模式
	if cfg.App.Env == "production" {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}

	// 创建 Gin 实例
	router := setupRouter(db, cfg)

	// 启动服务器
	log.Printf("服务器启动在 %s", cfg.App.Port)
	if err := router.Run(":" + cfg.App.Port); err != nil {
		log.Fatal("服务器启动失败:", err)
	}
}

func setupRouter(db *gorm.DB, cfg *config.Config) *gin.Engine {
	router := gin.Default()

	// 中间件
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.Use(middleware.CORSMiddleware())

	// 静态文件服务
	router.Static("/static", "./web/static")

	// 健康检查
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "ok",
			"service": cfg.App.Name,
			"version": cfg.App.Version,
		})
	})

	// API v1 路由组
	api := router.Group("/api/v1")
	{
		// ==================== 认证相关路由（无需认证） ====================
		authHandler := handlers.NewAuthHandler(db)
		{
			// 登录接口 - 无需认证
			api.POST("/auth/login", authHandler.Login)
		}

		// ==================== 需要认证的路由 ====================
		// 使用认证中间件
		authorized := api.Group("")
		authorized.Use(middleware.AuthMiddleware())
		// 添加默认密码检查中间件
		authorized.Use(middleware.DefaultPasswordCheckMiddleware(db))
		{
			// 认证相关 - 需要登录
			authorized.GET("/auth/current-user", authHandler.GetCurrentUser)
			authorized.POST("/auth/change-password", authHandler.ChangePassword)
			authorized.POST("/auth/logout", authHandler.Logout)
			authorized.GET("/auth/check-default-password", authHandler.CheckDefaultPassword)

			// 用户相关路由 - 需要登录（原接口，保留兼容）
			userHandler := handlers.NewUserHandler(db)
			authorized.GET("/currentUser", userHandler.GetCurrentUser) // 保留原接口
			authorized.GET("/users", userHandler.GetUsers)
			authorized.GET("/users/:id", userHandler.GetUser)
			authorized.POST("/users", userHandler.CreateUser)
			authorized.PUT("/users/:id", userHandler.UpdateUser)
			authorized.DELETE("/users/:id", userHandler.DeleteUser)

			// 住户相关api - 需要登录
			personHandler := handlers.NewPersonHandler(db)
			authorized.GET("/getBuildingNumbers", personHandler.GetBuildingNumbers)
			authorized.GET("/getUnitNumbersByBuildingNumber", personHandler.GetUnitNumbersByBuildingNumber)
			authorized.POST("/getPersons", personHandler.GetPersons)
			authorized.GET("/getPersonStatistics", personHandler.GetPersonStatistics)
			authorized.POST("/getRooms", personHandler.GetRooms)
			authorized.GET("/getPersonInfo", personHandler.GetPersonInfo)
			authorized.GET("/getPersonInfoByRoom", personHandler.GetPersonInfoByRoom)
		}
	}

	// 默认路由
	router.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "接口不存在",
			"path":  c.Request.URL.Path,
		})
	})

	return router
}
