# MyBatis Generator GUI  

<p align="center">
  <img src="https://img.shields.io/badge/Go-1.20+-00ADD8?style=flat&logo=go"/>
  <img src="https://img.shields.io/badge/Gin-Web-00ACD7?style=flat&logo=go"/>
  <img src="https://img.shields.io/badge/License-Apache%202.0-blue.svg"/>
</p>

基于Go语言和Gin框架开发的MyBatis代码生成器Web应用，用于快速生成MyBatis的Java实体类、Mapper接口和XML映射文件。

## ✨ 核心特性

- 🌐 **Web界面** - 使用现代化Web技术，浏览器访问，无需安装
- 🗄️ **数据库支持** - 支持MySQL和PostgreSQL数据库
- 💾 **配置管理** - SQLite本地存储，保存数据库连接和生成配置
- 🔄 **自动命名转换** - 数据库下划线命名自动转换为Java驼峰命名
- 📝 **注释生成** - 从数据库列注释自动生成Java代码注释
- 🎯 **灵活配置** - 支持Lombok、JSR310日期类型、分页查询等多种选项
- 📦 **完整代码** - 一键生成Java实体类、Mapper接口和MyBatis XML文件
- 🚀 **RESTful API** - 提供完整的REST API接口

## 🎯 功能列表

### 数据库连接管理
- ✅ MySQL数据库连接
- ✅ PostgreSQL数据库连接  
- ✅ 连接配置保存和管理
- ✅ 数据库连接测试
- ✅ 表列表查看和过滤

### 代码生成
- ✅ Java实体类生成（支持标准Bean和Lombok两种风格）
- ✅ Mapper接口生成
- ✅ MyBatis XML映射文件生成
- ✅ 驼峰命名自动转换
- ✅ 数据库注释转Java注释
- ✅ 主键自动识别
- ✅ 分页查询支持
- ✅ JSR310日期类型支持

### 生成选项
- **注释生成**: 使用数据库表和列的注释生成Java注释
- **Lombok支持**: 使用@Data注解简化实体类代码
- **分页查询**: 生成分页查询方法
- **JSR310**: 使用LocalDate、LocalDateTime等现代日期类型
- **覆盖XML**: 重新生成时覆盖已存在的XML文件

## 🚀 快速开始

### 系统要求

- Go 1.20或更高版本
- MySQL 5.7+或PostgreSQL 9.0+数据库（用于连接测试）
- 现代浏览器（Chrome、Firefox、Edge等）

###  安装步骤

#### 方式一：从源码运行

```bash
# 1. 克隆仓库
git clone https://github.com/yourusername/mybatis-generator-gui-go.git
cd mybatis-generator-gui-go

# 2. 下载依赖
go mod tidy

# 3. 运行程序
go run cmd/main.go

# 4. 浏览器访问
# 打开浏览器访问: http://localhost:8080
```

#### 方式二：编译后运行

```bash
# 编译（同时生成Windows和Linux版本）
.\build.bat          # Windows
./build.sh           # Linux

# 运行可执行文件
.\mybatis-generator-gui-windows-amd64.exe    # Windows
./mybatis-generator-gui-linux-amd64          # Linux

# 浏览器访问
# 打开浏览器访问: http://localhost:8080
```

## 📖 使用说明

### 1. 启动应用

```bash
go run cmd/main.go
```

应用启动后，在浏览器中访问：**http://localhost:8080**

### 2. 创建数据库连接

1. 点击左侧"+ 新建连接"按钮
2. 填写数据库连接信息：
   - 连接名称：自定义名称，用于标识连接
   - 数据库类型：选择MySQL或PostgreSQL
   - 主机：数据库服务器地址（例如：localhost）
   - 端口：数据库端口（MySQL默认3306，PostgreSQL默认5432）
   - 数据库名：要连接的数据库名称
   - 用户名和密码：数据库登录凭证
3. 点击"测试连接"按钮验证连接
4. 点击"保存"保存连接配置

### 3. 选择数据库表

1. 在左侧连接列表中点击已保存的连接
2. 表列表会自动加载
3. 点击表名，右侧配置面板会自动填充表信息

### 4. 配置代码生成选项

1. **项目目录**：输入Java项目的根目录（例如：`D:\project\mybatis-demo`）
2. **包名配置**：
   - Model包名：实体类的包名（例如：`com.example.model`）
   - DAO包名：Mapper接口的包名（例如：`com.example.mapper`）
   - Mapper包名：XML文件的包名（例如：`mapper`）
3. **目标文件夹**：
   - Model目标文件夹：通常为`src/main/java`
   - DAO目标文件夹：通常为`src/main/java`
   - Mapper目标文件夹：通常为`src/main/resources`
4. **类名配置**：
   - 实体类名：自动根据表名生成（可修改）
   - Mapper名：Mapper接口名称（可修改）
5. **选择生成选项**：勾选需要的选项（注释、Lombok、分页等）
6. **编码格式**：选择生成文件的编码（推荐UTF-8）

### 5. 生成代码

1. 确认所有配置无误
2. 点击"生成代码"按钮
3. 等待生成完成提示
4. 到项目目录查看生成的文件

### 6. 保存配置

点击"保存配置"按钮可以保存当前的代码生成配置，下次使用时可以快速加载。

## 🎨 生成代码示例

### 数据库表

```sql
CREATE TABLE user (
    id BIGINT PRIMARY KEY AUTO_INCREMENT COMMENT '用户ID',
    user_name VARCHAR(50) NOT NULL COMMENT '用户名',
    email VARCHAR(100) COMMENT '邮箱',
    created_at DATETIME COMMENT '创建时间'
);
```

### 生成的Java实体类（Lombok风格）

```java
package com.example.model;

import lombok.Data;
import java.io.Serializable;
import java.util.Date;

/**
 * 用户表
 */
@Data
public class User implements Serializable {
    private static final long serialVersionUID = 1L;

    /** 用户ID */
    private Long id;

    /** 用户名 */
    private String userName;

    /** 邮箱 */
    private String email;

    /** 创建时间 */
    private Date createdAt;
}
```

### 生成的Mapper接口

```java
package com.example.mapper;

import com.example.model.User;
import java.util.List;
import org.apache.ibatis.annotations.Param;

/**
 * UserMapper接口
 */
public interface UserMapper {
    int deleteByPrimaryKey(Long id);
    int insert(User record);
    int insertSelective(User record);
    User selectByPrimaryKey(Long id);
    int updateByPrimaryKeySelective(User record);
    int updateByPrimaryKey(User record);
    List<User> selectByPage(@Param("offset") int offset, @Param("limit") int limit);
}
```

### 生成的Mapper XML（部分）

```xml
<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE mapper PUBLIC "-//mybatis.org//DTD Mapper 3.0//EN" 
"http://mybatis.org/dtd/mybatis-3-mapper.dtd">
<mapper namespace="com.example.mapper.UserMapper">
    <resultMap id="BaseResultMap" type="com.example.model.User">
        <id column="id" jdbcType="BIGINT" property="id" />
        <result column="user_name" jdbcType="VARCHAR" property="userName" />
        <result column="email" jdbcType="VARCHAR" property="email" />
        <result column="created_at" jdbcType="TIMESTAMP" property="createdAt" />
    </resultMap>

    <sql id="Base_Column_List">
        id, user_name, email, created_at
    </sql>

    <select id="selectByPrimaryKey" parameterType="java.lang.Long" resultMap="BaseResultMap">
        SELECT <include refid="Base_Column_List" />
        FROM user
        WHERE id = #{id,jdbcType=BIGINT}
    </select>
    
    <!-- 更多SQL语句... -->
</mapper>
```

## 🛠️ 技术栈

- **语言**: Go 1.20+
- **Web框架**: Gin v1.9+
- **前端**: HTML5 + CSS3 + JavaScript (原生)
- **数据库驱动**: 
  - MySQL: go-sql-driver/mysql
  - PostgreSQL: lib/pq  
- **本地存储**: mattn/go-sqlite3
- **模板引擎**: text/template (Go标准库)

## 📂 项目结构

```
mybatis-generator-gui-go/
├── build.bat                  # Windows构建脚本
├── build.sh                   # Linux构建脚本
├── workflow.bat               # Windows完整工作流
├── workflow.sh                # Linux完整工作流
├── cmd/                       # 主程序入口
│   └── main.go               # Web服务器
├── internal/                  # 内部包
│   ├── config/               # 配置管理
│   │   ├── database_config.go  # 数据库配置模型
│   │   ├── generator_config.go # 生成配置模型
│   │   └── storage.go          # SQLite存储
│   ├── database/             # 数据库操作
│   │   ├── connector.go        # 数据库连接
│   │   ├── types.go            # 表结构类型
│   │   └── type_mapping.go     # 类型映射
│   ├── generator/            # 代码生成器
│   │   ├── generator.go        # 生成器主逻辑
│   │   ├── model_template.go   # Model模板
│   │   ├── mapper_template.go  # Mapper模板
│   │   └── mapper_xml_template.go # XML模板
│   ├── api/                  # REST API
│   │   ├── database_api.go     # 数据库API
│   │   └── generator_api.go    # 代码生成API
│   ├── web/                  # Web资源
│   │   ├── embed.go           # 资源嵌入
│   │   ├── templates/         # HTML模板
│   │   │   └── index.html
│   │   └── static/            # 静态资源
│   │       ├── css/style.css
│   │       └── js/app.js
│   └── utils/                # 工具函数
│       └── string_utils.go     # 字符串处理
├── resources/                # 资源文件
├── go.mod                    # Go模块定义
└── README.md                 # 本文件
```

## 🧪 运行测试

```bash
# 运行所有测试
go test ./...

# 运行特定包的测试
go test ./internal/utils

# 运行测试并显示覆盖率
go test ./... -cover
```

## 🤝 贡献指南

欢迎贡献代码！请遵循以下步骤：

1. Fork本仓库
2. 创建特性分支 (`git checkout -b feature/AmazingFeature`)
3. 提交更改 (`git commit -m 'Add some AmazingFeature'`)
4. 推送到分支 (`git push origin feature/AmazingFeature`)
5. 开启Pull Request

## 📄 许可证

本项目采用 Apache 2.0 许可证 - 查看 [LICENSE](LICENSE) 文件了解详情

## 📞 联系方式

- 项目主页: https://github.com/yourusername/mybatis-generator-gui-go
- 问题反馈: https://github.com/yourusername/mybatis-generator-gui-go/issues

## 🙏 致谢

本项目参考了原Java版本的 [mybatis-generator-gui](https://github.com/zouzg/mybatis-generator-gui) 项目。

---

⭐ 如果这个项目对你有帮助，请给个Star支持一下！
