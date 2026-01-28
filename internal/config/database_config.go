package config

// DatabaseConfig 数据库连接配置
type DatabaseConfig struct {
	ID       int    `json:"id"`       // 主键ID
	Name     string `json:"name"`     // 配置名称
	DbType   string `json:"dbType"`   // 数据库类型: MySQL, PostgreSQL
	Host     string `json:"host"`     // 主机地址
	Port     string `json:"port"`     // 端口号
	Schema   string `json:"schema"`   // 数据库名/Schema
	Username string `json:"username"` // 用户名
	Password string `json:"password"` // 密码
	Encoding string `json:"encoding"` // 编码格式,默认UTF-8
}

// DbType 数据库类型常量
const (
	DbTypeMySQL      = "MySQL"
	DbTypePostgreSQL = "PostgreSQL"
	DbTypeOracle     = "Oracle"
)
