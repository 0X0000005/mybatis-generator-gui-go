package main

import (
	"flag"
	"fmt"
	"html/template"
	"io/fs"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/yourusername/mybatis-generator-gui-go/internal/api"
	"github.com/yourusername/mybatis-generator-gui-go/internal/config"
	"github.com/yourusername/mybatis-generator-gui-go/internal/generator"
	"github.com/yourusername/mybatis-generator-gui-go/internal/web"
)

const version = "1.6.1"

func main() {
	// 解析命令行参数
	port := flag.Int("p", 8080, "服务器端口号")
	showVersion := flag.Bool("v", false, "显示版本号")
	showHelp := flag.Bool("h", false, "显示帮助信息")
	flag.Parse()

	// 显示版本号
	if *showVersion {
		fmt.Printf("MyBatis Generator GUI v%s\n", version)
		os.Exit(0)
	}

	// 显示帮助信息
	if *showHelp {
		fmt.Printf("MyBatis Generator GUI v%s\n\n", version)
		fmt.Println("用法: mybatis-generator-gui [选项]")
		fmt.Println()
		fmt.Println("选项:")
		fmt.Println("  -p <端口>   指定服务器端口号 (默认: 8080)")
		fmt.Println("  -v          显示版本号")
		fmt.Println("  -h          显示帮助信息")
		fmt.Println()
		fmt.Println("示例:")
		fmt.Println("  mybatis-generator-gui -p 9090")
		os.Exit(0)
	}

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
	addr := fmt.Sprintf(":%d", *port)
	log.Printf("MyBatis Generator GUI 启动成功!")
	log.Printf("访问地址: http://localhost%s", addr)
	log.Printf("版本: %s", version)

	if err := r.Run(addr); err != nil {
		log.Fatalf("服务器启动失败: %v", err)
	}
}
