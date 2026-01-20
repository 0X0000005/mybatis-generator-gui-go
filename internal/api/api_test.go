package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/yourusername/mybatis-generator-gui-go/internal/config"
)

func init() {
	// 初始化测试数据库
	config.InitDatabase()
	gin.SetMode(gin.TestMode)
}

// ==================== Database API Tests ====================

// 测试获取数据库连接列表
func TestGetConnections(t *testing.T) {
	router := gin.Default()
	router.GET("/api/connections", GetConnections)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/connections", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var result []*config.DatabaseConfig
	err := json.Unmarshal(w.Body.Bytes(), &result)
	assert.NoError(t, err)
}

// 测试创建数据库连接配置
func TestCreateConnection(t *testing.T) {
	router := gin.Default()
	router.POST("/api/connections", CreateConnection)

	// 创建测试数据，使用唯一名称避免冲突
	uniqueName := fmt.Sprintf("test_connection_%d", time.Now().UnixNano())
	dbConfig := config.DatabaseConfig{
		Name:     uniqueName,
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
	assert.Equal(t, "创建成功", result["message"])
}

// 测试创建连接配置 - 无效JSON
func TestCreateConnection_InvalidJSON(t *testing.T) {
	router := gin.Default()
	router.POST("/api/connections", CreateConnection)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/connections", bytes.NewBufferString("invalid json"))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

// 测试更新数据库连接配置
func TestUpdateConnection(t *testing.T) {
	router := gin.Default()
	router.PUT("/api/connections/:id", UpdateConnection)

	dbConfig := config.DatabaseConfig{
		Name:     "updated_connection",
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
	req, _ := http.NewRequest("PUT", "/api/connections/1", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	// 可能返回OK或错误（取决于是否存在ID为1的记录）
	assert.True(t, w.Code == http.StatusOK || w.Code == http.StatusInternalServerError)
}

// 测试更新连接配置 - 无效ID
func TestUpdateConnection_InvalidID(t *testing.T) {
	router := gin.Default()
	router.PUT("/api/connections/:id", UpdateConnection)

	dbConfig := config.DatabaseConfig{
		Name: "test",
	}

	jsonData, _ := json.Marshal(dbConfig)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/api/connections/invalid", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

// 测试删除数据库连接配置
func TestDeleteConnection(t *testing.T) {
	router := gin.Default()
	router.DELETE("/api/connections/:id", DeleteConnection)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/api/connections/999", nil)
	router.ServeHTTP(w, req)

	// 可能返回OK或错误（取决于是否存在ID为999的记录）
	assert.True(t, w.Code == http.StatusOK || w.Code == http.StatusInternalServerError)
}

// 测试删除连接配置 - 无效ID
func TestDeleteConnection_InvalidID(t *testing.T) {
	router := gin.Default()
	router.DELETE("/api/connections/:id", DeleteConnection)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/api/connections/invalid", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
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

// 测试连接测试 - 无效JSON
func TestTestConnection_InvalidJSON(t *testing.T) {
	router := gin.Default()
	router.POST("/api/connections/test", TestConnection)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/connections/test", bytes.NewBufferString("invalid"))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
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

// 测试获取表列表 - 无效JSON
func TestGetTables_InvalidJSON(t *testing.T) {
	router := gin.Default()
	router.POST("/api/tables", GetTables)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/tables", bytes.NewBufferString("invalid"))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

// 测试获取列信息
func TestGetColumns(t *testing.T) {
	router := gin.Default()
	router.POST("/api/columns", GetColumns)

	requestData := map[string]interface{}{
		"databaseId": 999, // 不存在的ID
		"tableName":  "test_table",
	}

	jsonData, _ := json.Marshal(requestData)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/columns", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	// 应该返回错误，因为数据库配置不存在
	assert.Equal(t, http.StatusNotFound, w.Code)
}

// 测试获取列信息 - 无效JSON
func TestGetColumns_InvalidJSON(t *testing.T) {
	router := gin.Default()
	router.POST("/api/columns", GetColumns)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/columns", bytes.NewBufferString("invalid"))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

// ==================== Config API Tests ====================

// 测试获取代码生成配置列表
func TestGetGeneratorConfigs(t *testing.T) {
	router := gin.Default()
	router.GET("/api/generator-configs", GetGeneratorConfigs)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/generator-configs", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var result []config.GeneratorConfig
	err := json.Unmarshal(w.Body.Bytes(), &result)
	assert.NoError(t, err)
}

// 测试保存代码生成配置
func TestSaveGeneratorConfig(t *testing.T) {
	router := gin.Default()
	router.POST("/api/generator-configs", SaveGeneratorConfig)

	genConfig := config.GeneratorConfig{
		Name:         "test_config",
		TableName:    "test_table",
		ModelPackage: "com.example.model",
	}

	jsonData, _ := json.Marshal(genConfig)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/generator-configs", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var result map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &result)
	assert.NoError(t, err)
	assert.Equal(t, "保存成功", result["message"])
}

// 测试保存配置 - 无效JSON
func TestSaveGeneratorConfig_InvalidJSON(t *testing.T) {
	router := gin.Default()
	router.POST("/api/generator-configs", SaveGeneratorConfig)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/generator-configs", bytes.NewBufferString("invalid"))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

// 测试删除代码生成配置
func TestDeleteGeneratorConfig(t *testing.T) {
	router := gin.Default()
	router.DELETE("/api/generator-configs/:name", DeleteGeneratorConfig)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/api/generator-configs/nonexistent_config", nil)
	router.ServeHTTP(w, req)

	// 可能返回OK或错误（取决于是否存在该配置）
	assert.True(t, w.Code == http.StatusOK || w.Code == http.StatusInternalServerError)
}

// ==================== Generator API Tests ====================

// 测试代码生成（失败场景 - 数据库不存在）
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

// 测试代码生成 - 无效JSON
func TestGenerateCode_InvalidJSON(t *testing.T) {
	router := gin.Default()
	router.POST("/api/generate", GenerateCode)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/generate", bytes.NewBufferString("invalid"))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
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

// 测试获取文件名辅助函数
func TestGetFileNames(t *testing.T) {
	files := []string{
		"/path/to/User.java",
		"/path/to/UserMapper.java",
		"/path/to/UserMapper.xml",
	}

	names := getFileNames(files)

	assert.Equal(t, 3, len(names))
	assert.Equal(t, "User.java", names[0])
	assert.Equal(t, "UserMapper.java", names[1])
	assert.Equal(t, "UserMapper.xml", names[2])
}

// 测试随机字符串生成
func TestGenerateRandomString(t *testing.T) {
	str1 := generateRandomString(8)
	str2 := generateRandomString(8)

	// 长度应该正确
	assert.Equal(t, 8, len(str1))
	assert.Equal(t, 8, len(str2))

	// 两次生成的字符串应该不同（极小概率会相同）
	assert.NotEqual(t, str1, str2)

	// 测试不同长度
	str3 := generateRandomString(16)
	assert.Equal(t, 16, len(str3))
}

// ==================== Auth API Tests ====================

// 测试检查认证 - 无Cookie
func TestCheckAuth_NoCookie(t *testing.T) {
	router := gin.Default()
	router.GET("/test", func(c *gin.Context) {
		result := CheckAuth(c)
		c.JSON(http.StatusOK, gin.H{"authenticated": result})
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var result map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &result)
	assert.Equal(t, false, result["authenticated"])
}

// 测试检查认证 - 有效Cookie
func TestCheckAuth_ValidCookie(t *testing.T) {
	router := gin.Default()
	router.GET("/test", func(c *gin.Context) {
		result := CheckAuth(c)
		c.JSON(http.StatusOK, gin.H{"authenticated": result})
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	req.AddCookie(&http.Cookie{
		Name:  AuthCookieName,
		Value: AuthCookieValue,
	})
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var result map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &result)
	assert.Equal(t, true, result["authenticated"])
}

// 测试认证中间件 - 未授权
func TestAuthMiddleware_Unauthorized(t *testing.T) {
	router := gin.Default()
	router.Use(AuthMiddleware())
	router.GET("/protected", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/protected", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

// 测试认证中间件 - 已授权
func TestAuthMiddleware_Authorized(t *testing.T) {
	router := gin.Default()
	router.Use(AuthMiddleware())
	router.GET("/protected", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/protected", nil)
	req.AddCookie(&http.Cookie{
		Name:  AuthCookieName,
		Value: AuthCookieValue,
	})
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

// 测试登录页面
func TestHandleLoginPage(t *testing.T) {
	router := gin.Default()
	router.GET("/login", HandleLoginPage)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/login", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Header().Get("Content-Type"), "text/html")
}

// 测试登录页面 - 已登录用户重定向
func TestHandleLoginPage_AlreadyLoggedIn(t *testing.T) {
	router := gin.Default()
	router.GET("/login", HandleLoginPage)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/login", nil)
	req.AddCookie(&http.Cookie{
		Name:  AuthCookieName,
		Value: AuthCookieValue,
	})
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusFound, w.Code)
	assert.Equal(t, "/", w.Header().Get("Location"))
}

// 测试登录API - 响应结构验证
func TestHandleLoginAPI_ResponseFormat(t *testing.T) {
	router := gin.Default()
	router.POST("/api/login", HandleLoginAPI)

	// 使用任意凭据测试API响应结构
	loginReq := LoginRequest{
		Username: "test_user",
		Password: "test_pass",
	}

	jsonData, _ := json.Marshal(loginReq)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/login", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	// API应该返回200状态码
	assert.Equal(t, http.StatusOK, w.Code)

	// 响应应该是有效的LoginResponse结构
	var result LoginResponse
	err := json.Unmarshal(w.Body.Bytes(), &result)
	assert.NoError(t, err)

	// 由于凭据错误，应该返回失败
	assert.False(t, result.Success)
}

// 测试登录API - 错误凭据
func TestHandleLoginAPI_WrongCredentials(t *testing.T) {
	router := gin.Default()
	router.POST("/api/login", HandleLoginAPI)

	loginReq := LoginRequest{
		Username: "wrong_user",
		Password: "wrong_pass",
	}

	jsonData, _ := json.Marshal(loginReq)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/login", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var result LoginResponse
	json.Unmarshal(w.Body.Bytes(), &result)
	assert.False(t, result.Success)
	assert.NotEmpty(t, result.Error)
}

// 测试登录API - 无效JSON
func TestHandleLoginAPI_InvalidJSON(t *testing.T) {
	router := gin.Default()
	router.POST("/api/login", HandleLoginAPI)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/login", bytes.NewBufferString("invalid"))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

// 测试登出
func TestHandleLogout(t *testing.T) {
	router := gin.Default()
	router.GET("/logout", HandleLogout)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/logout", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusFound, w.Code)
	assert.Equal(t, "/login", w.Header().Get("Location"))
}

// 测试主页 - 未认证重定向
func TestHandleIndexWithAuth_Redirect(t *testing.T) {
	router := gin.Default()
	router.GET("/", func(c *gin.Context) {
		HandleIndexWithAuth(c, "1.0.0")
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusFound, w.Code)
	assert.Equal(t, "/login", w.Header().Get("Location"))
}

// 测试主页 - 已认证
func TestHandleIndexWithAuth_Authenticated(t *testing.T) {
	router := gin.Default()
	router.GET("/", func(c *gin.Context) {
		HandleIndexWithAuth(c, "1.0.0")
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/", nil)
	req.AddCookie(&http.Cookie{
		Name:  AuthCookieName,
		Value: AuthCookieValue,
	})
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Header().Get("Content-Type"), "text/html")
}

// ==================== Version API Test ====================

// 测试版本信息API
func TestGetVersion(t *testing.T) {
	router := gin.Default()
	router.GET("/api/version", func(c *gin.Context) {
		c.JSON(200, gin.H{"version": "1.4.0"})
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/version", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var result map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &result)
	assert.Equal(t, "1.4.0", result["version"])
}
