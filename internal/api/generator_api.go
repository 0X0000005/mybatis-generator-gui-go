package api

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

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
		TableNames []string               `json:"tableNames"`
		Config     config.GeneratorConfig `json:"config"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("ERROR: 解析请求失败: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if len(req.TableNames) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请选择至少一张表"})
		return
	}

	log.Printf("INFO: 开始生成代码 - DatabaseID: %d, Tables: %v", req.DatabaseID, req.TableNames)

	// 获取当前工作目录
	currentDir, err := os.Getwd()
	if err != nil {
		log.Printf("ERROR: 获取当前目录失败: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取当前目录失败"})
		return
	}

	// 创建temp目录下的随机子目录（避免多用户冲突）
	timestamp := time.Now().Format("20060102_150405")
	randomSuffix := generateRandomString(8)
	projectSubDir := fmt.Sprintf("gen_%s_%s", timestamp, randomSuffix)

	tempBaseDir := filepath.Join(currentDir, "temp")
	req.Config.ProjectFolder = filepath.Join(tempBaseDir, projectSubDir)

	log.Printf("INFO: 使用临时目录: %s", req.Config.ProjectFolder)

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

	// 为每张表生成代码
	var allFiles []string
	for _, tableName := range req.TableNames {
		// 复制配置并设置当前表
		tableConfig := req.Config
		tableConfig.TableName = tableName
		tableConfig.DomainObjectName = toPascalCase(tableName)
		tableConfig.MapperName = toPascalCase(tableName) + "Mapper"

		log.Printf("INFO: 生成表 %s 的代码", tableName)

		gen := generator.NewGenerator(&tableConfig, dbConfig)
		files, err := gen.Generate()
		if err != nil {
			log.Printf("ERROR: 生成表 %s 代码失败: %v", tableName, err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("生成表 %s 失败: %v", tableName, err)})
			return
		}
		allFiles = append(allFiles, files...)
	}

	log.Printf("INFO: 成功生成 %d 张表, 共 %d 个文件", len(req.TableNames), len(allFiles))

	// 打包成ZIP
	zipName := fmt.Sprintf("generated_%d_tables", len(req.TableNames))
	zipPath, err := generator.CreateZipArchive(allFiles, req.Config.ProjectFolder, zipName)
	if err != nil {
		log.Printf("ERROR: 打包失败: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "打包失败: " + err.Error()})
		return
	}

	log.Printf("INFO: ZIP文件已创建: %s", zipPath)

	// 生成唯一的下载ID并存储映射
	downloadID := fmt.Sprintf("%d_multi_%s", req.DatabaseID, generateRandomString(8))

	generatedZipsMu.Lock()
	generatedZips[downloadID] = zipPath
	generatedZipsMu.Unlock()

	log.Printf("INFO: 下载ID已创建: %s -> %s", downloadID, zipPath)

	c.JSON(http.StatusOK, gin.H{
		"success":    true,
		"message":    "代码生成成功",
		"downloadId": downloadID,
		"files":      getFileNames(allFiles),
		"tableCount": len(req.TableNames),
	})
}

// toPascalCase 将下划线命名转为帕斯卡命名
func toPascalCase(s string) string {
	parts := strings.Split(s, "_")
	for i, part := range parts {
		if len(part) > 0 {
			parts[i] = strings.ToUpper(part[:1]) + strings.ToLower(part[1:])
		}
	}
	return strings.Join(parts, "")
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

// generateRandomString 生成随机字符串
func generateRandomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyz0123456789"
	result := make([]byte, length)
	for i := range result {
		result[i] = charset[rand.Intn(len(charset))]
	}
	return string(result)
}
