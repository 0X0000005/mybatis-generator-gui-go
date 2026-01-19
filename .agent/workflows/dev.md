---
description: 开发工作流程
---

# 开发工作流程

## 1. 需求分析

- 明确功能需求和边界
- 分析对现有代码的影响
- 创建需求分析文档（如需要）
- 确认实现方案后再开始编码

## 2. 编码实现

- 修改配置结构体 (`internal/config/`)
- 修改模板文件 (`internal/generator/*_template.go`)
- 修改生成逻辑 (`internal/generator/generator.go`)
- 修改API处理 (`internal/api/`)
- 修改前端界面 (`internal/web/templates/`, `internal/web/static/js/`)

## 3. 功能测试

### 3.1 单元测试
// turbo
```powershell
go test ./... -v
```

### 3.2 HTTP接口测试
每个新增/修改的HTTP接口必须在 `internal/api/api_test.go` 中添加测试用例。

测试模板：
```go
func TestXxxAPI(t *testing.T) {
    router := setupTestRouter()
    w := httptest.NewRecorder()
    
    // 准备请求
    body := gin.H{"key": "value"}
    jsonBody, _ := json.Marshal(body)
    req, _ := http.NewRequest("POST", "/api/xxx", bytes.NewBuffer(jsonBody))
    req.Header.Set("Content-Type", "application/json")
    
    // 执行请求
    router.ServeHTTP(w, req)
    
    // 验证结果
    assert.Equal(t, http.StatusOK, w.Code)
}
```

## 4. 脚本编译

// turbo
```powershell
.\build.bat
```

编译产出：
- `mgg.exe` (Windows)
- `mgg` (Linux)

## 5. 更新版本

### 版本号规范
- **大版本号** (如 1.x.x → 2.x.x): 由用户决定
- **小版本号** (如 1.4.0 → 1.4.1): AI可自行更新

### 需更新的文件
1. `cmd/main.go` - `version` 常量
2. `internal/api/auth_api.go` - `AuthCookieValue`

## 6. Git推送

// turbo
```powershell
git add -A
git commit -m "v版本号: 更新说明"
git tag v版本号
git push
git push --tags
```
