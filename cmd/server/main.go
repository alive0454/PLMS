package main

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"log"
	"net/http"

	"PLMS/internal/config"
	"PLMS/internal/database"
	"PLMS/internal/handlers"
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

	// API 路由组
	api := router.Group("/api/v1")
	{
		// 用户相关路由
		userHandler := handlers.NewUserHandler(db)
		api.GET("/users", userHandler.GetUsers)
		api.GET("/users/:id", userHandler.GetUser)
		api.POST("/users", userHandler.CreateUser)
		api.PUT("/users/:id", userHandler.UpdateUser)
		api.DELETE("/users/:id", userHandler.DeleteUser)
		//住户相关api
		personHandler := handlers.NewPersonHandler(db)
		api.GET("/getBuildingNumbers", personHandler.GetBuildingNumbers)
		api.GET("/getUnitNumbersByBuildingNumber", personHandler.GetUnitNumbersByBuildingNumber)
		api.POST("/getPersons", personHandler.GetPersons)
		api.GET("/getPersonStatistics", personHandler.GetPersonStatistics)
		api.POST("/getRooms", personHandler.GetRooms)
		api.GET("/getPersonInfo", personHandler.GetPersonInfo)
		api.GET("/getPersonInfoByRoom", personHandler.GetPersonInfoByRoom)

		//User相关api
		api.GET("/currentUser", userHandler.GetCurrentUser)
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
