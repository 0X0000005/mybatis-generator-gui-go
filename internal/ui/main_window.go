package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

// MainWindow 主窗口
type MainWindow struct {
	window        fyne.Window
	dbTree        *DBTreeView
	configPanel   *ConfigPanel
	connectionMgr *ConnectionManager
}

// NewMainWindow 创建主窗口
func NewMainWindow(window fyne.Window) *MainWindow {
	mw := &MainWindow{
		window: window,
	}

	mw.connectionMgr = NewConnectionManager(window)
	mw.dbTree = NewDBTreeView(window, mw.connectionMgr)
	mw.configPanel = NewConfigPanel(window, mw.dbTree)

	return mw
}

// Build 构建主界面
func (mw *MainWindow) Build() fyne.CanvasObject {
	// 顶部工具栏
	toolbar := mw.buildToolbar()

	// 左侧数据库树
	leftPanel := mw.dbTree.Build()

	// 右侧配置面板
	rightPanel := mw.configPanel.Build()

	// 分割布局
	splitContent := container.NewHSplit(
		leftPanel,
		rightPanel,
	)
	splitContent.Offset = 0.25 // 左侧占25%

	// 主布局
	mainContent := container.NewBorder(
		toolbar,      // 顶部
		nil,          // 底部
		nil,          // 左侧
		nil,          // 右侧
		splitContent, // 中间
	)

	return mainContent
}

// buildToolbar 构建工具栏
func (mw *MainWindow) buildToolbar() fyne.CanvasObject {
	newConnBtn := widget.NewButton("新建连接", func() {
		mw.connectionMgr.ShowNewConnectionDialog()
	})

	configBtn := widget.NewButton("配置管理", func() {
		mw.configPanel.ShowConfigManager()
	})

	toolbar := container.NewHBox(
		newConnBtn,
		configBtn,
	)

	return container.NewVBox(
		toolbar,
		widget.NewSeparator(),
	)
}
