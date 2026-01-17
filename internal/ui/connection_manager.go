package ui

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"

	"github.com/yourusername/mybatis-generator-gui-go/internal/config"
	"github.com/yourusername/mybatis-generator-gui-go/internal/database"
)

// ConnectionManager 连接管理器
type ConnectionManager struct {
	window fyne.Window
}

// NewConnectionManager 创建连接管理器
func NewConnectionManager(window fyne.Window) *ConnectionManager {
	return &ConnectionManager{
		window: window,
	}
}

// ShowNewConnectionDialog 显示新建连接对话框
func (cm *ConnectionManager) ShowNewConnectionDialog() {
	cm.showConnectionDialog(nil, false)
}

// ShowEditConnectionDialog 显示编辑连接对话框
func (cm *ConnectionManager) ShowEditConnectionDialog(dbConfig *config.DatabaseConfig) {
	cm.showConnectionDialog(dbConfig, true)
}

// showConnectionDialog 显示连接对话框
func (cm *ConnectionManager) showConnectionDialog(dbConfig *config.DatabaseConfig, isEdit bool) {
	// 创建表单
	nameEntry := widget.NewEntry()
	dbTypeSelect := widget.NewSelect([]string{config.DbTypeMySQL, config.DbTypePostgreSQL}, nil)
	hostEntry := widget.NewEntry()
	portEntry := widget.NewEntry()
	schemaEntry := widget.NewEntry()
	usernameEntry := widget.NewEntry()
	passwordEntry := widget.NewPasswordEntry()

	// 设置默认值
	if dbConfig != nil {
		nameEntry.SetText(dbConfig.Name)
		dbTypeSelect.SetSelected(dbConfig.DbType)
		hostEntry.SetText(dbConfig.Host)
		portEntry.SetText(dbConfig.Port)
		schemaEntry.SetText(dbConfig.Schema)
		usernameEntry.SetText(dbConfig.Username)
		passwordEntry.SetText(dbConfig.Password)
	} else {
		dbTypeSelect.SetSelected(config.DbTypeMySQL)
		hostEntry.SetText("localhost")
		portEntry.SetText("3306")
	}

	// 创建表单项
	form := &widget.Form{
		Items: []*widget.FormItem{
			{Text: "连接名称:", Widget: nameEntry},
			{Text: "数据库类型:", Widget: dbTypeSelect},
			{Text: "主机:", Widget: hostEntry},
			{Text: "端口:", Widget: portEntry},
			{Text: "数据库/Schema:", Widget: schemaEntry},
			{Text: "用户名:", Widget: usernameEntry},
			{Text: "密码:", Widget: passwordEntry},
		},
	}

	// 测试连接按钮
	var testBtn *widget.Button
	testBtn = widget.NewButton("测试连接", func() {
		testConfig := &config.DatabaseConfig{
			DbType:   dbTypeSelect.Selected,
			Host:     hostEntry.Text,
			Port:     portEntry.Text,
			Schema:   schemaEntry.Text,
			Username: usernameEntry.Text,
			Password: passwordEntry.Text,
		}

		if err := database.TestConnection(testConfig); err != nil {
			dialog.ShowError(fmt.Errorf("连接失败: %v", err), cm.window)
		} else {
			dialog.ShowInformation("成功", "数据库连接成功!", cm.window)
		}
	})

	content := container.NewBorder(
		nil,
		testBtn,
		nil,
		nil,
		form,
	)

	// 创建对话框
	var title string
	if isEdit {
		title = "编辑数据库连接"
	} else {
		title = "新建数据库连接"
	}

	d := dialog.NewCustomConfirm(title, "保存", "取消", content, func(save bool) {
		if !save {
			return
		}

		// 验证输入
		if nameEntry.Text == "" || dbTypeSelect.Selected == "" ||
			hostEntry.Text == "" || portEntry.Text == "" ||
			schemaEntry.Text == "" || usernameEntry.Text == "" {
			dialog.ShowError(fmt.Errorf("请填写所有必填项"), cm.window)
			return
		}

		// 创建配置
		newConfig := &config.DatabaseConfig{
			Name:     nameEntry.Text,
			DbType:   dbTypeSelect.Selected,
			Host:     hostEntry.Text,
			Port:     portEntry.Text,
			Schema:   schemaEntry.Text,
			Username: usernameEntry.Text,
			Password: passwordEntry.Text,
			Encoding: "utf8mb4",
		}

		if isEdit && dbConfig != nil {
			newConfig.ID = dbConfig.ID
		}

		// 保存配置
		if err := config.SaveDatabaseConfig(newConfig, isEdit); err != nil {
			dialog.ShowError(fmt.Errorf("保存失败: %v", err), cm.window)
			return
		}

		dialog.ShowInformation("成功", "保存成功!", cm.window)
	}, cm.window)

	d.Resize(fyne.NewSize(500, 400))
	d.Show()
}
