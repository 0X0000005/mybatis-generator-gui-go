---
description: 版本发布流程
---

# 版本发布流程

## 版本号规范

- **大版本号** (如 1.x.x → 2.x.x): 由用户决定，AI不得自行修改
- **小版本号** (如 1.4.0 → 1.4.1): AI可以自行更新

## 发布前检查

// turbo
1. 运行测试确保通过
```powershell
go test ./... -v
```

## 发布步骤

// turbo-all

1. 更新版本号
   - 修改 `cmd/main.go` 中的 `version` 常量
   - 修改 `internal/api/auth_api.go` 中的 `AuthCookieValue`

2. 编译项目
```powershell
.\build.bat
```

3. Git提交
```powershell
git add -A
git commit -m "v版本号: 更新说明"
```

4. 创建标签
```powershell
git tag v版本号
```

5. 推送
```powershell
git push
git push --tags
```
