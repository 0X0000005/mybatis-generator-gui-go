package main

import (
	"html/template"
	"io/fs"
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/yourusername/mybatis-generator-gui-go/internal/api"
	"github.com/yourusername/mybatis-generator-gui-go/internal/config"
	"github.com/yourusername/mybatis-generator-gui-go/internal/generator" // Added by user instruction
	"github.com/yourusername/mybatis-generator-gui-go/internal/web"
)

const version = "1.2.0"

func main() {
	// 初始化随机种子
	rand.Seed(time.Now().UnixNano())

	// 启动ZIP文件清理任务
	generator.StartCleanupScheduler()

	// 设置Gin模式初始化配置数据库
	if err := config.InitDatabase(); err != nil {
		log.Fatalf("初始化数据库失败: %v", err)
	}
	defer config.CloseDatabase()

	// 创建Gin路由
	r := gin.Default()

	// 加载嵌入的HTML模板
	tmpl, err := template.ParseFS(web.TemplatesFS, "templates/*.html")
	if err != nil {
		log.Fatalf("加载模板失败: %v", err)
	}
	r.SetHTMLTemplate(tmpl)

	// 设置静态文件服务（使用嵌入的文件系统）
	staticSub, _ := fs.Sub(web.StaticFS, "static")
	r.StaticFS("/static", http.FS(staticSub))

	// 登录相关路由
	r.GET("/login", api.HandleLoginPage)
	r.POST("/api/login", api.HandleLoginAPI)
	r.GET("/logout", api.HandleLogout)

	// 主页（需要认证）
	r.GET("/", func(c *gin.Context) {
		api.HandleIndexWithAuth(c, version)
	})

	// API路由组（需要认证）
	apiGroup := r.Group("/api")
	apiGroup.Use(api.AuthMiddleware())
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
		apiGroup.GET("/download/:id", api.DownloadCode)

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
