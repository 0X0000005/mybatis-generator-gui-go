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

// DBTreeView 数据库树形视图
type DBTreeView struct {
	window        fyne.Window
	connMgr       *ConnectionManager
	tree          *widget.Tree
	configs       []*config.DatabaseConfig
	selectedTable string
	selectedDB    *config.DatabaseConfig
	onTableSelect func(dbConfig *config.DatabaseConfig, tableName string)
}

// NewDBTreeView 创建数据库树形视图
func NewDBTreeView(window fyne.Window, connMgr *ConnectionManager) *DBTreeView {
	dbt := &DBTreeView{
		window:  window,
		connMgr: connMgr,
	}

	dbt.loadDatabaseConfigs()
	dbt.buildTree()

	return dbt
}

// Build 构建UI
func (dbt *DBTreeView) Build() fyne.CanvasObject {
	refreshBtn := widget.NewButton("刷新", func() {
		dbt.Refresh()
	})

	filterEntry := widget.NewEntry()
	filterEntry.SetPlaceHolder("过滤表名...")

	return container.NewBorder(
		container.NewVBox(
			refreshBtn,
			filterEntry,
		),
		nil,
		nil,
		nil,
		container.NewScroll(dbt.tree),
	)
}

// Refresh 刷新树
func (dbt *DBTreeView) Refresh() {
	dbt.loadDatabaseConfigs()
	dbt.tree.Refresh()
}

// SetOnTableSelect 设置表选择回调
func (dbt *DBTreeView) SetOnTableSelect(callback func(*config.DatabaseConfig, string)) {
	dbt.onTableSelect = callback
}

// loadDatabaseConfigs 加载数据库配置
func (dbt *DBTreeView) loadDatabaseConfigs() {
	configs, err := config.LoadDatabaseConfigs()
	if err != nil {
		dialog.ShowError(fmt.Errorf("加载数据库配置失败: %v", err), dbt.window)
		return
	}
	dbt.configs = configs
}

// buildTree 构建树
func (dbt *DBTreeView) buildTree() {
	dbt.tree = widget.NewTree(
		func(uid widget.TreeNodeID) []widget.TreeNodeID {
			// 根节点
			if uid == "" {
				var ids []string
				for _, cfg := range dbt.configs {
					ids = append(ids, fmt.Sprintf("db_%d", cfg.ID))
				}
				return ids
			}
			// 数据库节点的子节点暂时为空，需要展开时加载
			return []string{}
		},
		func(uid widget.TreeNodeID) bool {
			// 根节点和数据库节点是分支
			return uid == "" || len(uid) > 0 && uid[0:3] == "db_"
		},
		func(branch bool) fyne.CanvasObject {
			return widget.NewLabel("Template")
		},
		func(uid widget.TreeNodeID, branch bool, obj fyne.CanvasObject) {
			label := obj.(*widget.Label)
			if uid == "" {
				label.SetText("数据库连接")
			} else if uid[0:3] == "db_" {
				// 查找对应的数据库配置
				for _, cfg := range dbt.configs {
					if fmt.Sprintf("db_%d", cfg.ID) == uid {
						label.SetText(cfg.Name)
						break
					}
				}
			} else {
				// 表名
				label.SetText(uid)
			}
		},
	)

	// 双击展开数据库连接时加载表
	dbt.tree.OnSelected = func(uid widget.TreeNodeID) {
		if uid[0:3] == "db_" {
			// 查找对应的数据库配置
			for _, cfg := range dbt.configs {
				if fmt.Sprintf("db_%d", cfg.ID) == uid {
					dbt.selectedDB = cfg
					dbt.loadTables(cfg)
					break
				}
			}
		} else if dbt.selectedDB != nil {
			// 选中了表
			dbt.selectedTable = uid
			if dbt.onTableSelect != nil {
				dbt.onTableSelect(dbt.selectedDB, uid)
			}
		}
	}
}

// loadTables 加载表列表
func (dbt *DBTreeView) loadTables(dbConfig *config.DatabaseConfig) {
	connector := database.NewConnector(dbConfig)
	if err := connector.Connect(); err != nil {
		dialog.ShowError(fmt.Errorf("连接数据库失败: %v", err), dbt.window)
		return
	}
	defer connector.Close()

	tables, err := connector.GetTableNames("")
	if err != nil {
		dialog.ShowError(fmt.Errorf("获取表列表失败: %v", err), dbt.window)
		return
	}

	// 这里需要一个更复杂的实现来动态更新树节点
	// 暂时简化处理，显示信息
	dialog.ShowInformation("表列表", fmt.Sprintf("数据库 %s 包含 %d 个表", dbConfig.Name, len(tables)), dbt.window)
}
