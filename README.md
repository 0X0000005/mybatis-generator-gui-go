# MyBatis Generator GUI  

<p align="center">
  <img src="https://img.shields.io/badge/Go-1.20+-00ADD8?style=flat&logo=go"/>
  <img src="https://img.shields.io/badge/Gin-Web-00ACD7?style=flat&logo=go"/>
  <img src="https://img.shields.io/badge/License-Apache%202.0-blue.svg"/>
</p>

基于Go语言和Gin框架开发的MyBatis代码生成器Web应用，用于快速生成MyBatis的Java实体类、Mapper接口和XML映射文件。

## 📸 界面预览

### 登录界面
![登录界面](docs/screenshots/login.png)

### 主界面
![主界面](docs/screenshots/main.png)

## ✨ 核心特性

- 🌐 **Web界面** - 现代化Web技术，浏览器访问，无需安装
- 🗄️ **多数据库支持** - MySQL 和 PostgreSQL
- 💾 **配置持久化** - SQLite本地存储连接和生成配置
- 🔄 **自动命名转换** - 下划线命名自动转换为驼峰命名
- 📦 **完整代码生成** - 实体类、Mapper接口、XML映射文件一键生成
- 📥 **ZIP自动打包** - 生成后自动打包下载
- 🔐 **登录认证** - 基于Cookie的简单认证保护

## 🚀 快速开始

### 方式一：从源码运行

```bash
git clone https://github.com/yourusername/mybatis-generator-gui-go.git
cd mybatis-generator-gui-go
go mod tidy
go run cmd/main.go
```

### 方式二：编译后运行

```bash
# Windows
.\build.bat
.\mgg.exe

# Linux
./build.sh
./mgg
```

启动后访问：**http://localhost:8080**

### 命令行参数

```bash
./mgg -p 9090  # 指定端口
./mgg -v       # 显示版本
./mgg -h       # 显示帮助
```

## 📖 使用说明

1. **创建数据库连接** - 点击"新建连接"，填写连接信息，点击"测试连接"验证
2. **选择表** - 选择连接后，加载表列表，选择要生成代码的表
3. **配置生成选项** - 设置包名、输出目录、生成选项等
4. **生成代码** - 点击"生成代码"按钮，自动下载ZIP文件

### 生成选项说明

| 选项 | 说明 |
|------|------|
| 注释生成 | 从数据库注释生成Java注释 |
| Lombok | 使用@Data注解简化代码 |
| 分页查询 | 生成分页查询方法 |
| JSR310 | 使用LocalDate/LocalDateTime |
| 覆盖XML | 重新生成时覆盖已存在的XML |

## 🛠️ 技术栈

| 组件 | 技术 |
|------|------|
| 语言 | Go 1.20+ |
| Web框架 | Gin |
| 前端 | HTML5 + CSS3 + JavaScript |
| 数据库驱动 | go-sql-driver/mysql, lib/pq |
| 本地存储 | SQLite |

## 📂 项目结构

```
mybatis-generator-gui-go/
├── cmd/main.go          # 入口文件
├── internal/
│   ├── api/             # REST API
│   ├── config/          # 配置管理
│   ├── database/        # 数据库操作
│   ├── generator/       # 代码生成器
│   ├── utils/           # 工具函数
│   └── web/             # 前端资源
├── docs/screenshots/    # 截图
├── build.bat/.sh        # 构建脚本
└── README.md
```

## 🧪 运行测试

```bash
go test ./...            # 运行测试
go test ./... -cover     # 显示覆盖率
```

## 📄 许可证

Apache 2.0 - 查看 [LICENSE](LICENSE) 了解详情

## 🙏 致谢

本项目参考了 [mybatis-generator-gui](https://github.com/zouzg/mybatis-generator-gui)

---

⭐ 如果这个项目对你有帮助，请给个Star支持！
