package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/yourusername/mybatis-generator-gui-go/internal/api"
	"github.com/yourusername/mybatis-generator-gui-go/internal/config"
)

const version = "1.0.0"

func main() {
	// 初始化配置数据库
	if err := config.InitDatabase(); err != nil {
		log.Fatalf("初始化数据库失败: %v", err)
	}
	defer config.CloseDatabase()

	// 创建Gin路由
	r := gin.Default()

	// 设置静态文件目录
	r.Static("/static", "./internal/web/static")
	r.LoadHTMLGlob("internal/web/templates/*")

	// 主页
	r.GET("/", func(c *gin.Context) {
		c.HTML(200, "index.html", gin.H{
			"version": version,
		})
	})

	// API路由组
	apiGroup := r.Group("/api")
	{
		// 数据库连接管理
		apiGroup.GET("/connections", api.GetConnections)
		apiGroup.POST("/connections", api.CreateConnection)
		apiGroup.PUT("/connections/:id", api.UpdateConnection)
		apiGroup.DELETE("/connections/:id", api.DeleteConnection)
		apiGroup.POST("/connections/test", api.TestConnection)

		// 数据库表操作
		apiGroup.POST("/tables", api.GetTables)
		apiGroup.POST("/columns", api.GetColumns)

		// 代码生成配置
		apiGroup.GET("/generator-configs", api.GetGeneratorConfigs)
		apiGroup.POST("/generator-configs", api.SaveGeneratorConfig)
		apiGroup.DELETE("/generator-configs/:name", api.DeleteGeneratorConfig)

		// 代码生成
		apiGroup.POST("/generate", api.GenerateCode)

		// 版本信息
		apiGroup.GET("/version", func(c *gin.Context) {
			c.JSON(200, gin.H{"version": version})
		})
	}

	// 启动服务器
	log.Printf("MyBatis Generator GUI 启动成功!")
	log.Printf("访问地址: http://localhost:8080")
	log.Printf("版本: %s", version)

	if err := r.Run(":8080"); err != nil {
		log.Fatalf("服务器启动失败: %v", err)
	}
}
