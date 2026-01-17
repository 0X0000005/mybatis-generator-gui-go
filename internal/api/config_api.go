package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/yourusername/mybatis-generator-gui-go/internal/config"
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
