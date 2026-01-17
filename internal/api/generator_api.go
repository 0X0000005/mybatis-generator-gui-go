package api

import (
	"net/http"

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
	if err := gen.Generate(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "代码生成成功",
		"files": []string{
			req.Config.DomainObjectName + ".java",
			req.Config.MapperName + ".java",
			req.Config.MapperName + ".xml",
		},
	})
}
