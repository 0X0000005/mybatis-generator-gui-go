package api

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/yourusername/mybatis-generator-gui-go/internal/config"
	"github.com/yourusername/mybatis-generator-gui-go/internal/generator"
)

// 存储生成的ZIP文件映射（线程安全）
var (
	generatedZips   = make(map[string]string)
	generatedZipsMu sync.RWMutex
)

// GenerateCode 生成代码
func GenerateCode(c *gin.Context) {
	var req struct {
		DatabaseID int                    `json:"databaseId"`
		Config     config.GeneratorConfig `json:"config"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("ERROR: 解析请求失败: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	log.Printf("INFO: 开始生成代码 - DatabaseID: %d, Table: %s", req.DatabaseID, req.Config.TableName)

	// 如果ProjectFolder为空，使用临时目录
	if req.Config.ProjectFolder == "" {
		req.Config.ProjectFolder = filepath.Join(os.TempDir(), "mybatis-gen")
		log.Printf("INFO: 使用默认临时目录: %s", req.Config.ProjectFolder)
	}

	// 加载数据库配置
	configs, err := config.LoadDatabaseConfigs()
	if err != nil {
		log.Printf("ERROR: 加载数据库配置失败: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var dbConfig *config.DatabaseConfig
	for _, cfg := range configs {
		if cfg.ID == req.DatabaseID {
			dbConfig = cfg
			break
		}
	}

	if dbConfig == nil {
		log.Printf("ERROR: 数据库配置不存在 - ID: %d", req.DatabaseID)
		c.JSON(http.StatusNotFound, gin.H{"error": "数据库配置不存在"})
		return
	}

	log.Printf("INFO: 使用数据库配置: %s (%s)", dbConfig.Name, dbConfig.DbType)

	// 创建生成器并生成代码
	gen := generator.NewGenerator(&req.Config, dbConfig)
	files, err := gen.Generate()
	if err != nil {
		log.Printf("ERROR: 生成代码失败: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	log.Printf("INFO: 成功生成 %d 个文件", len(files))

	// 打包成ZIP
	zipPath, err := generator.CreateZipArchive(files, req.Config.ProjectFolder, req.Config.TableName)
	if err != nil {
		log.Printf("ERROR: 打包失败: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "打包失败: " + err.Error()})
		return
	}

	log.Printf("INFO: ZIP文件已创建: %s", zipPath)

	// 生成唯一的下载ID并存储映射
	downloadID := fmt.Sprintf("%d_%s", req.DatabaseID, req.Config.TableName)

	generatedZipsMu.Lock()
	generatedZips[downloadID] = zipPath
	generatedZipsMu.Unlock()

	log.Printf("INFO: 下载ID已创建: %s -> %s", downloadID, zipPath)

	c.JSON(http.StatusOK, gin.H{
		"success":    true,
		"message":    "代码生成成功",
		"downloadId": downloadID,
		"files":      getFileNames(files),
	})
}

// DownloadCode 下载生成的代码ZIP
func DownloadCode(c *gin.Context) {
	downloadID := c.Param("id")
	log.Printf("INFO: 请求下载 - ID: %s", downloadID)

	generatedZipsMu.RLock()
	zipPath, exists := generatedZips[downloadID]
	generatedZipsMu.RUnlock()

	if !exists {
		log.Printf("ERROR: 下载ID不存在: %s", downloadID)
		c.JSON(http.StatusNotFound, gin.H{"error": "下载文件不存在"})
		return
	}

	log.Printf("INFO: 找到ZIP文件: %s", zipPath)

	// 检查文件是否存在
	if _, err := os.Stat(zipPath); os.IsNotExist(err) {
		log.Printf("ERROR: ZIP文件已被删除: %s", zipPath)
		generatedZipsMu.Lock()
		delete(generatedZips, downloadID)
		generatedZipsMu.Unlock()
		c.JSON(http.StatusNotFound, gin.H{"error": "文件已过期"})
		return
	}

	log.Printf("INFO: 开始发送文件: %s", zipPath)

	// 设置响应头
	c.Header("Content-Description", "File Transfer")
	c.Header("Content-Transfer-Encoding", "binary")
	c.Header("Content-Disposition", "attachment; filename="+filepath.Base(zipPath))
	c.Header("Content-Type", "application/zip")

	// 发送文件
	c.File(zipPath)

	log.Printf("INFO: 文件发送成功，文件将在5分钟后自动清理: %s", zipPath)
}

// getFileNames 从完整路径中提取文件名
func getFileNames(files []string) []string {
	names := make([]string, len(files))
	for i, file := range files {
		names[i] = filepath.Base(file)
	}
	return names
}
