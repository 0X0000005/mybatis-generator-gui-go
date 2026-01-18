package api

import (
	"html/template"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/yourusername/mybatis-generator-gui-go/internal/web"
	"golang.org/x/crypto/bcrypt"
)

const (
	// AuthCookieName Cookie 名称
	AuthCookieName = "mgg_session"
	// AuthCookieValue 基于版本号，每次更新版本后自动失效旧 session
	AuthCookieValue = "auth_v1.3.1"
	// UsernameHash 用户名的 bcrypt hash
	UsernameHash = "$2a$10$knnRQBicLZhhk5DXHmxHR.z9/z2MNxPd1GUtmqZZhKso6dDcgIj2K"
	// PasswordHash 密码的 bcrypt hash
	PasswordHash = "$2a$10$1eF.gu83pocTbjZeULFuY.t.4.AqxIVnf6qxR3Fd4aKwX4Vb1u51i"
)

// LoginRequest 登录请求结构
type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// LoginResponse 登录响应结构
type LoginResponse struct {
	Success bool   `json:"success"`
	Error   string `json:"error,omitempty"`
}

// CheckAuth 检查是否已登录
func CheckAuth(c *gin.Context) bool {
	cookie, err := c.Cookie(AuthCookieName)
	if err != nil {
		return false
	}
	return cookie == AuthCookieValue
}

// AuthMiddleware 认证中间件
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if !CheckAuth(c) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "未授权访问"})
			c.Abort()
			return
		}
		c.Next()
	}
}

// HandleLoginPage 渲染登录页面
func HandleLoginPage(c *gin.Context) {
	log.Printf("[访问] %s %s - 登录页", c.Request.Method, c.Request.URL.Path)

	// 如果已登录，重定向到主页
	if CheckAuth(c) {
		log.Printf("[授权] 用户已登录，重定向到主页")
		c.Redirect(http.StatusFound, "/")
		return
	}

	// 读取登录页面模板
	data, err := web.TemplatesFS.ReadFile("templates/login.html")
	if err != nil {
		log.Printf("[错误] 读取 login.html 失败: %v", err)
		c.String(http.StatusInternalServerError, "无法加载登录页面")
		return
	}

	c.Data(http.StatusOK, "text/html; charset=utf-8", data)
}

// HandleLoginAPI 处理登录 API
func HandleLoginAPI(c *gin.Context) {
	log.Printf("[访问] %s %s - 登录验证接口", c.Request.Method, c.Request.URL.Path)

	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("[错误] 解析登录请求失败: %v", err)
		c.JSON(http.StatusBadRequest, LoginResponse{
			Success: false,
			Error:   "解析请求失败",
		})
		return
	}

	// 使用 bcrypt 验证用户名和密码
	usernameMatch := bcrypt.CompareHashAndPassword([]byte(UsernameHash), []byte(req.Username)) == nil
	passwordMatch := bcrypt.CompareHashAndPassword([]byte(PasswordHash), []byte(req.Password)) == nil

	if usernameMatch && passwordMatch {
		log.Printf("[成功] 登录验证通过")
		// 设置 Cookie
		c.SetCookie(
			AuthCookieName,
			AuthCookieValue,
			3600*24, // 24小时
			"/",
			"",
			false,
			true, // HttpOnly
		)
		c.JSON(http.StatusOK, LoginResponse{Success: true})
	} else {
		log.Printf("[失败] 登录验证失败")
		c.JSON(http.StatusOK, LoginResponse{
			Success: false,
			Error:   "账号或密码不正确",
		})
	}
}

// HandleLogout 处理登出
func HandleLogout(c *gin.Context) {
	log.Printf("[访问] %s %s - 登出", c.Request.Method, c.Request.URL.Path)

	// 删除 Cookie
	c.SetCookie(
		AuthCookieName,
		"",
		-1, // MaxAge -1 表示删除
		"/",
		"",
		false,
		true,
	)

	c.Redirect(http.StatusFound, "/login")
}

// HandleIndexWithAuth 处理主页请求（带认证检查）
func HandleIndexWithAuth(c *gin.Context, version string) {
	log.Printf("[访问] %s %s - 主页", c.Request.Method, c.Request.URL.Path)

	// 检查认证
	if !CheckAuth(c) {
		log.Printf("[授权] 未登录，重定向到 /login")
		c.Redirect(http.StatusFound, "/login")
		return
	}

	// 加载并渲染模板
	tmpl, err := template.ParseFS(web.TemplatesFS, "templates/index.html")
	if err != nil {
		log.Printf("[错误] 加载模板失败: %v", err)
		c.String(http.StatusInternalServerError, "无法加载主页面")
		return
	}

	c.Header("Content-Type", "text/html; charset=utf-8")
	err = tmpl.Execute(c.Writer, gin.H{
		"version": version,
	})
	if err != nil {
		log.Printf("[错误] 渲染模板失败: %v", err)
	}
}
