package database

import (
	"database/sql"
	"fmt"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	"github.com/yourusername/mybatis-generator-gui-go/internal/config"
)

// Connector 数据库连接器
type Connector struct {
	config *config.DatabaseConfig
	db     *sql.DB
}

// NewConnector 创建新的数据库连接器
func NewConnector(cfg *config.DatabaseConfig) *Connector {
	return &Connector{
		config: cfg,
	}
}

// Connect 连接到数据库
func (c *Connector) Connect() error {
	var dsn string
	var driverName string

	switch c.config.DbType {
	case config.DbTypeMySQL:
		driverName = "mysql"
		encoding := c.config.Encoding
		if encoding == "" {
			encoding = "utf8mb4"
		}
		dsn = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=True&loc=Local",
			c.config.Username,
			c.config.Password,
			c.config.Host,
			c.config.Port,
			c.config.Schema,
			encoding,
		)

	case config.DbTypePostgreSQL:
		driverName = "postgres"
		dsn = fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
			c.config.Host,
			c.config.Port,
			c.config.Username,
			c.config.Password,
			c.config.Schema,
		)

	default:
		return fmt.Errorf("不支持的数据库类型: %s", c.config.DbType)
	}

	db, err := sql.Open(driverName, dsn)
	if err != nil {
		return fmt.Errorf("打开数据库连接失败: %v", err)
	}

	// 测试连接
	if err := db.Ping(); err != nil {
		db.Close()
		return fmt.Errorf("数据库连接测试失败: %v", err)
	}

	c.db = db
	return nil
}

// Close 关闭数据库连接
func (c *Connector) Close() error {
	if c.db != nil {
		return c.db.Close()
	}
	return nil
}

// TestConnection 测试数据库连接
func TestConnection(cfg *config.DatabaseConfig) error {
	connector := NewConnector(cfg)
	if err := connector.Connect(); err != nil {
		return err
	}
	defer connector.Close()
	return nil
}

// GetTableNames 获取数据库中所有表名
func (c *Connector) GetTableNames(filter string) ([]string, error) {
	if c.db == nil {
		return nil, fmt.Errorf("数据库未连接")
	}

	var query string
	var args []interface{}

	switch c.config.DbType {
	case config.DbTypeMySQL:
		if filter != "" {
			query = "SELECT TABLE_NAME FROM information_schema.TABLES WHERE TABLE_SCHEMA = ? AND TABLE_NAME LIKE ?"
			args = []interface{}{c.config.Schema, "%" + filter + "%"}
		} else {
			query = "SELECT TABLE_NAME FROM information_schema.TABLES WHERE TABLE_SCHEMA = ?"
			args = []interface{}{c.config.Schema}
		}

	case config.DbTypePostgreSQL:
		if filter != "" {
			query = "SELECT tablename FROM pg_tables WHERE schemaname = 'public' AND tablename LIKE $1"
			args = []interface{}{"%" + filter + "%"}
		} else {
			query = "SELECT tablename FROM pg_tables WHERE schemaname = 'public'"
		}

	default:
		return nil, fmt.Errorf("不支持的数据库类型: %s", c.config.DbType)
	}

	rows, err := c.db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("查询表名失败: %v", err)
	}
	defer rows.Close()

	var tables []string
	for rows.Next() {
		var tableName string
		if err := rows.Scan(&tableName); err != nil {
			return nil, fmt.Errorf("读取表名失败: %v", err)
		}
		tables = append(tables, tableName)
	}

	return tables, nil
}

// GetTableColumns 获取表的列信息
func (c *Connector) GetTableColumns(tableName string) ([]*TableColumn, error) {
	if c.db == nil {
		return nil, fmt.Errorf("数据库未连接")
	}

	var query string
	var args []interface{}

	switch c.config.DbType {
	case config.DbTypeMySQL:
		query = `
			SELECT 
				COLUMN_NAME,
				DATA_TYPE,
				IFNULL(COLUMN_COMMENT, '') as COLUMN_COMMENT,
				IS_NULLABLE,
				IFNULL(COLUMN_KEY, '') as COLUMN_KEY,
				IFNULL(EXTRA, '') as EXTRA
			FROM information_schema.COLUMNS
			WHERE TABLE_SCHEMA = ? AND TABLE_NAME = ?
			ORDER BY ORDINAL_POSITION
		`
		args = []interface{}{c.config.Schema, tableName}

	case config.DbTypePostgreSQL:
		query = `
			SELECT 
				a.attname as column_name,
				pg_catalog.format_type(a.atttypid, a.atttypmod) as data_type,
				COALESCE(pg_catalog.col_description(a.attrelid, a.attnum), '') as column_comment,
				NOT a.attnotnull as is_nullable,
				CASE WHEN pk.conname IS NOT NULL THEN 'PRI' ELSE '' END as column_key,
				'' as extra
			FROM pg_catalog.pg_attribute a
			LEFT JOIN pg_catalog.pg_class c ON a.attrelid = c.oid
			LEFT JOIN pg_catalog.pg_namespace n ON c.relnamespace = n.oid
			LEFT JOIN pg_catalog.pg_constraint pk ON pk.conrelid = c.oid AND a.attnum = ANY(pk.conkey) AND pk.contype = 'p'
			WHERE c.relname = $1
				AND n.nspname = 'public'
				AND a.attnum > 0
				AND NOT a.attisdropped
			ORDER BY a.attnum
		`
		args = []interface{}{tableName}

	default:
		return nil, fmt.Errorf("不支持的数据库类型: %s", c.config.DbType)
	}

	rows, err := c.db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("查询列信息失败: %v", err)
	}
	defer rows.Close()

	var columns []*TableColumn
	for rows.Next() {
		var column TableColumn
		var isNullableStr string

		if err := rows.Scan(
			&column.ColumnName,
			&column.DataType,
			&column.ColumnComment,
			&isNullableStr,
			&column.ColumnKey,
			&column.Extra,
		); err != nil {
			return nil, fmt.Errorf("读取列信息失败: %v", err)
		}

		// 处理 IS_NULLABLE
		column.IsNullable = strings.ToUpper(isNullableStr) == "YES" || isNullableStr == "true" || isNullableStr == "t"

		columns = append(columns, &column)
	}

	return columns, nil
}

// GetTableComment 获取表注释
func (c *Connector) GetTableComment(tableName string) (string, error) {
	if c.db == nil {
		return "", fmt.Errorf("数据库未连接")
	}

	var query string
	var args []interface{}
	var comment string

	switch c.config.DbType {
	case config.DbTypeMySQL:
		query = `
			SELECT IFNULL(TABLE_COMMENT, '') as TABLE_COMMENT
			FROM information_schema.TABLES
			WHERE TABLE_SCHEMA = ? AND TABLE_NAME = ?
		`
		args = []interface{}{c.config.Schema, tableName}

	case config.DbTypePostgreSQL:
		query = `
			SELECT COALESCE(obj_description(c.oid), '') as table_comment
			FROM pg_catalog.pg_class c
			LEFT JOIN pg_catalog.pg_namespace n ON n.oid = c.relnamespace
			WHERE c.relname = $1 AND n.nspname = 'public'
		`
		args = []interface{}{tableName}

	default:
		return "", fmt.Errorf("不支持的数据库类型: %s", c.config.DbType)
	}

	err := c.db.QueryRow(query, args...).Scan(&comment)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", nil
		}
		return "", fmt.Errorf("查询表注释失败: %v", err)
	}

	return comment, nil
}
