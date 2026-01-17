package main

import (
	"log"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"

	"github.com/yourusername/mybatis-generator-gui-go/internal/config"
	"github.com/yourusername/mybatis-generator-gui-go/internal/ui"
)

const version = "1.0.0"

func main() {
	// 初始化配置数据库
	if err := config.InitDatabase(); err != nil {
		log.Fatalf("初始化数据库失败: %v", err)
	}
	defer config.CloseDatabase()

	// 创建Fyne应用
	myApp := app.New()
	myWindow := myApp.NewWindow("MyBatis Generator GUI - v" + version)

	// 创建主界面
	mainUI := ui.NewMainWindow(myWindow)
	content := mainUI.Build()

	myWindow.SetContent(content)
	myWindow.Resize(fyne.NewSize(1200, 800))
	myWindow.ShowAndRun()
}
