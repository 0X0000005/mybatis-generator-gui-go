package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/yourusername/mybatis-generator-gui-go/internal/config"
)

func init() {
	// 初始化测试数据库
	config.InitDatabase()
	gin.SetMode(gin.TestMode)
}

// 测试获取数据库连接列表
func TestGetDatabaseConfigs(t *testing.T) {
	router := gin.Default()
	router.GET("/api/connections", GetDatabaseConfigs)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/connections", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var result []config.DatabaseConfig
	err := json.Unmarshal(w.Body.Bytes(), &result)
	assert.NoError(t, err)
}

// 测试保存数据库连接配置
func TestSaveDatabaseConfig(t *testing.T) {
	router := gin.Default()
	router.POST("/api/connections", SaveDatabaseConfig)

	// 创建测试数据
	dbConfig := config.DatabaseConfig{
		Name:     "test_connection",
		DbType:   "mysql",
		Host:     "localhost",
		Port:     "3306",
		Schema:   "test_db",
		Username: "root",
		Password: "password",
		Encoding: "utf8mb4",
	}

	jsonData, _ := json.Marshal(dbConfig)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/connections", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var result map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &result)
	assert.NoError(t, err)
	assert.Equal(t, "保存成功", result["message"])
}

// 测试连接测试
func TestTestConnection(t *testing.T) {
	router := gin.Default()
	router.POST("/api/connections/test", TestConnection)

	// 创建测试数据（使用无效连接来避免真正连接数据库）
	dbConfig := config.DatabaseConfig{
		DbType:   "mysql",
		Host:     "invalid_host",
		Port:     "3306",
		Schema:   "test_db",
		Username: "root",
		Password: "password",
	}

	jsonData, _ := json.Marshal(dbConfig)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/connections/test", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var result map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &result)
	assert.NoError(t, err)
	// 连接应该失败，因为是无效主机
	assert.Equal(t, false, result["success"])
}

// 测试获取表列表
func TestGetTables(t *testing.T) {
	router := gin.Default()
	router.POST("/api/tables", GetTables)

	requestData := map[string]interface{}{
		"databaseId": 999, // 不存在的ID
		"filter":     "",
	}

	jsonData, _ := json.Marshal(requestData)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/tables", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	// 应该返回错误，因为数据库配置不存在
	assert.Equal(t, http.StatusNotFound, w.Code)
}

// 测试代码生成（失败场景）
func TestGenerateCode_InvalidDatabase(t *testing.T) {
	router := gin.Default()
	router.POST("/api/generate", GenerateCode)

	requestData := map[string]interface{}{
		"databaseId": 999, // 不存在的ID
		"config": map[string]interface{}{
			"tableName":        "test_table",
			"domainObjectName": "TestEntity",
		},
	}

	jsonData, _ := json.Marshal(requestData)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/generate", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

// 测试下载不存在的文件
func TestDownloadCode_NotFound(t *testing.T) {
	router := gin.Default()
	router.GET("/api/download/:id", DownloadCode)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/download/invalid_id", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}
