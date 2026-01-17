package api

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/yourusername/mybatis-generator-gui-go/internal/config"
	"github.com/yourusername/mybatis-generator-gui-go/internal/generator"
)

// GetGeneratorConfigs 获取所有代码生成配置
func GetGeneratorConfigs(c *gin.Context) {
	configs, err := config.LoadGeneratorConfigs()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, configs)
}

// SaveGeneratorConfig 保存代码生成配置
func SaveGeneratorConfig(c *gin.Context) {
	var genConfig config.GeneratorConfig
	if err := c.ShouldBindJSON(&genConfig); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := config.SaveGeneratorConfig(&genConfig); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "保存成功"})
}

// DeleteGeneratorConfig 删除代码生成配置
func DeleteGeneratorConfig(c *gin.Context) {
	name := c.Param("name")
	if err := config.DeleteGeneratorConfig(name); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "删除成功"})
}

// 存储生成的ZIP文件映射
var generatedZips = make(map[string]string)

// GenerateCode 生成代码
func GenerateCode(c *gin.Context) {
	var req struct {
		DatabaseID int                    `json:"databaseId"`
		Config     config.GeneratorConfig `json:"config"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 如果ProjectFolder为空，使用临时目录
	if req.Config.ProjectFolder == "" {
		req.Config.ProjectFolder = filepath.Join(os.TempDir(), "mybatis-gen")
	}

	// 加载数据库配置
	configs, err := config.LoadDatabaseConfigs()
	if err != nil {
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
		c.JSON(http.StatusNotFound, gin.H{"error": "数据库配置不存在"})
		return
	}

	// 创建生成器并生成代码
	gen := generator.NewGenerator(&req.Config, dbConfig)
	files, err := gen.Generate()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 打包成ZIP
	zipPath, err := generator.CreateZipArchive(files, req.Config.ProjectFolder)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "打包失败: " + err.Error()})
		return
	}

	// 生成唯一的下载ID
	downloadID := fmt.Sprintf("%d_%s", req.DatabaseID, req.Config.TableName)
	generatedZips[downloadID] = zipPath

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

	zipPath, exists := generatedZips[downloadID]
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "下载文件不存在"})
		return
	}

	// 检查文件是否存在
	if _, err := os.Stat(zipPath); os.IsNotExist(err) {
		delete(generatedZips, downloadID)
		c.JSON(http.StatusNotFound, gin.H{"error": "文件已过期"})
		return
	}

	// 设置响应头
	c.Header("Content-Description", "File Transfer")
	c.Header("Content-Transfer-Encoding", "binary")
	c.Header("Content-Disposition", "attachment; filename="+filepath.Base(zipPath))
	c.Header("Content-Type", "application/zip")

	// 发送文件
	c.File(zipPath)

	// 发送后删除临时文件
	go func() {
		os.Remove(zipPath)
		delete(generatedZips, downloadID)
	}()
}

// getFileNames 从完整路径中提取文件名
func getFileNames(files []string) []string {
	names := make([]string, len(files))
	for i, file := range files {
		names[i] = filepath.Base(file)
	}
	return names
}
