package ui

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"

	"github.com/yourusername/mybatis-generator-gui-go/internal/config"
	"github.com/yourusername/mybatis-generator-gui-go/internal/generator"
	"github.com/yourusername/mybatis-generator-gui-go/internal/utils"
)

// ConfigPanel 配置面板
type ConfigPanel struct {
	window    fyne.Window
	dbTree    *DBTreeView
	genConfig *config.GeneratorConfig
	dbConfig  *config.DatabaseConfig
	tableName string

	// UI组件
	projectFolderEntry      *widget.Entry
	modelPackageEntry       *widget.Entry
	daoPackageEntry         *widget.Entry
	mapperPackageEntry      *widget.Entry
	modelTargetFolderEntry  *widget.Entry
	daoTargetFolderEntry    *widget.Entry
	mapperTargetFolderEntry *widget.Entry
	tableNameEntry          *widget.Entry
	domainObjectNameEntry   *widget.Entry
	mapperNameEntry         *widget.Entry
	generateKeysEntry       *widget.Entry
	encodingSelect          *widget.Select

	// 选项复选框
	offsetLimitCheck *widget.Check
	commentCheck     *widget.Check
	overrideXMLCheck *widget.Check
	lombokCheck      *widget.Check
	jsr310Check      *widget.Check
}

// NewConfigPanel 创建配置面板
func NewConfigPanel(window fyne.Window, dbTree *DBTreeView) *ConfigPanel {
	cp := &ConfigPanel{
		window: window,
		dbTree: dbTree,
	}

	cp.genConfig = &config.GeneratorConfig{
		Encoding: "UTF-8",
	}

	// 设置表选择回调
	dbTree.SetOnTableSelect(func(dbConfig *config.DatabaseConfig, tableName string) {
		cp.onTableSelect(dbConfig, tableName)
	})

	return cp
}

// Build 构建UI
func (cp *ConfigPanel) Build() fyne.CanvasObject {
	cp.buildFormComponents()

	form := &widget.Form{
		Items: []*widget.FormItem{
			{Text: "项目目录:", Widget: container.NewBorder(nil, nil, nil, widget.NewButton("选择", cp.chooseProjectFolder), cp.projectFolderEntry)},
			{Text: "Model包名:", Widget: cp.modelPackageEntry},
			{Text: "Model目标文件夹:", Widget: cp.modelTargetFolderEntry},
			{Text: "DAO包名:", Widget: cp.daoPackageEntry},
			{Text: "DAO目标文件夹:", Widget: cp.daoTargetFolderEntry},
			{Text: "Mapper包名:", Widget: cp.mapperPackageEntry},
			{Text: "Mapper目标文件夹:", Widget: cp.mapperTargetFolderEntry},
			{Text: "表名:", Widget: cp.tableNameEntry},
			{Text: "实体类名:", Widget: cp.domainObjectNameEntry},
			{Text: "Mapper名:", Widget: cp.mapperNameEntry},
			{Text: "主键字段:", Widget: cp.generateKeysEntry},
			{Text: "编码格式:", Widget: cp.encodingSelect},
		},
	}

	options := container.NewVBox(
		cp.offsetLimitCheck,
		cp.commentCheck,
		cp.overrideXMLCheck,
		cp.lombokCheck,
		cp.jsr310Check,
	)

	generateBtn := widget.NewButton("生成代码", cp.generateCode)
	saveConfigBtn := widget.NewButton("保存配置", cp.saveConfig)

	buttons := container.NewHBox(
		generateBtn,
		saveConfigBtn,
	)

	return container.NewBorder(
		nil,
		buttons,
		nil,
		nil,
		container.NewVSplit(
			container.NewScroll(form),
			container.NewScroll(options),
		),
	)
}

// buildFormComponents 构建表单组件
func (cp *ConfigPanel) buildFormComponents() {
	cp.projectFolderEntry = widget.NewEntry()
	cp.modelPackageEntry = widget.NewEntry()
	cp.daoPackageEntry = widget.NewEntry()
	cp.mapperPackageEntry = widget.NewEntry()
	cp.modelTargetFolderEntry = widget.NewEntry()
	cp.daoTargetFolderEntry = widget.NewEntry()
	cp.mapperTargetFolderEntry = widget.NewEntry()
	cp.tableNameEntry = widget.NewEntry()
	cp.domainObjectNameEntry = widget.NewEntry()
	cp.mapperNameEntry = widget.NewEntry()
	cp.generateKeysEntry = widget.NewEntry()
	cp.encodingSelect = widget.NewSelect([]string{"UTF-8", "GBK"}, nil)
	cp.encodingSelect.SetSelected("UTF-8")

	// 设置默认值
	cp.modelTargetFolderEntry.SetText("src/main/java")
	cp.daoTargetFolderEntry.SetText("src/main/java")
	cp.mapperTargetFolderEntry.SetText("src/main/resources")

	cp.offsetLimitCheck = widget.NewCheck("生成分页查询", nil)
	cp.commentCheck = widget.NewCheck("生成注释", nil)
	cp.commentCheck.SetChecked(true)
	cp.overrideXMLCheck = widget.NewCheck("覆盖XML文件", nil)
	cp.lombokCheck = widget.NewCheck("使用Lombok", nil)
	cp.jsr310Check = widget.NewCheck("使用JSR310日期类型", nil)
}

// chooseProjectFolder 选择项目目录
func (cp *ConfigPanel) chooseProjectFolder() {
	dialog.ShowFolderOpen(func(dir fyne.ListableURI, err error) {
		if err != nil || dir == nil {
			return
		}
		cp.projectFolderEntry.SetText(dir.Path())
	}, cp.window)
}

// onTableSelect 表选择事件
func (cp *ConfigPanel) onTableSelect(dbConfig *config.DatabaseConfig, tableName string) {
	cp.dbConfig = dbConfig
	cp.tableName = tableName

	cp.tableNameEntry.SetText(tableName)
	cp.domainObjectNameEntry.SetText(utils.DBStringToPascalCase(tableName))
	cp.mapperNameEntry.SetText(utils.DBStringToPascalCase(tableName) + "Mapper")
}

// generateCode 生成代码
func (cp *ConfigPanel) generateCode() {
	// 更新配置
	cp.updateConfigFromForm()

	// 验证
	if err := cp.validateConfig(); err != nil {
		dialog.ShowError(err, cp.window)
		return
	}

	// 创建生成器
	gen := generator.NewGenerator(cp.genConfig, cp.dbConfig)

	// 生成代码
	if err := gen.Generate(); err != nil {
		dialog.ShowError(fmt.Errorf("生成代码失败: %v", err), cp.window)
		return
	}

	dialog.ShowInformation("成功", "代码生成成功!", cp.window)
}

// saveConfig 保存配置
func (cp *ConfigPanel) saveConfig() {
	// TODO: 弹出对话框让用户输入配置名称
	dialog.ShowInformation("提示", "保存配置功能待实现", cp.window)
}

// ShowConfigManager 显示配置管理器
func (cp *ConfigPanel) ShowConfigManager() {
	dialog.ShowInformation("提示", "配置管理器待实现", cp.window)
}

// updateConfigFromForm 从表单更新配置
func (cp *ConfigPanel) updateConfigFromForm() {
	cp.genConfig.ProjectFolder = cp.projectFolderEntry.Text
	cp.genConfig.ModelPackage = cp.modelPackageEntry.Text
	cp.genConfig.ModelPackageTargetFolder = cp.modelTargetFolderEntry.Text
	cp.genConfig.DaoPackage = cp.daoPackageEntry.Text
	cp.genConfig.DaoTargetFolder = cp.daoTargetFolderEntry.Text
	cp.genConfig.MappingXMLPackage = cp.mapperPackageEntry.Text
	cp.genConfig.MappingXMLTargetFolder = cp.mapperTargetFolderEntry.Text
	cp.genConfig.TableName = cp.tableNameEntry.Text
	cp.genConfig.DomainObjectName = cp.domainObjectNameEntry.Text
	cp.genConfig.MapperName = cp.mapperNameEntry.Text
	cp.genConfig.GenerateKeys = cp.generateKeysEntry.Text
	cp.genConfig.Encoding = cp.encodingSelect.Selected

	cp.genConfig.OffsetLimit = cp.offsetLimitCheck.Checked
	cp.genConfig.Comment = cp.commentCheck.Checked
	cp.genConfig.OverrideXML = cp.overrideXMLCheck.Checked
	cp.genConfig.UseLombokPlugin = cp.lombokCheck.Checked
	cp.genConfig.JSR310Support = cp.jsr310Check.Checked
}

// validateConfig 验证配置
func (cp *ConfigPanel) validateConfig() error {
	if cp.genConfig.ProjectFolder == "" {
		return fmt.Errorf("项目目录不能为空")
	}
	if cp.genConfig.TableName == "" {
		return fmt.Errorf("表名不能为空")
	}
	if cp.genConfig.DomainObjectName == "" {
		return fmt.Errorf("实体类名不能为空")
	}
	if cp.genConfig.ModelPackage == "" || cp.genConfig.DaoPackage == "" || cp.genConfig.MappingXMLPackage == "" {
		return fmt.Errorf("包名不能为空")
	}
	if cp.dbConfig == nil {
		return fmt.Errorf("请先选择数据库表")
	}
	return nil
}
