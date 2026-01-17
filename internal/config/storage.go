package config

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3" // SQLite驱动
)

const (
	configDir  = "config"
	dbFileName = "sqlite3.db"
)

var db *sql.DB

// InitDatabase 初始化数据库
func InitDatabase() error {
	// 创建config目录
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return fmt.Errorf("创建配置目录失败: %v", err)
	}

	dbPath := filepath.Join(configDir, dbFileName)

	// 打开SQLite数据库
	var err error
	db, err = sql.Open("sqlite3", dbPath)
	if err != nil {
		return fmt.Errorf("打开数据库失败: %v", err)
	}

	// 创建表结构
	if err := createTables(); err != nil {
		return fmt.Errorf("创建表失败: %v", err)
	}

	return nil
}

// createTables 创建数据库表
func createTables() error {
	// 创建数据库配置表
	dbsTable := `
	CREATE TABLE IF NOT EXISTS dbs (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL UNIQUE,
		value TEXT NOT NULL
	);`

	// 创建代码生成配置表
	generatorTable := `
	CREATE TABLE IF NOT EXISTS generator_config (
		name TEXT PRIMARY KEY,
		value TEXT NOT NULL
	);`

	_, err := db.Exec(dbsTable)
	if err != nil {
		return fmt.Errorf("创建dbs表失败: %v", err)
	}

	_, err = db.Exec(generatorTable)
	if err != nil {
		return fmt.Errorf("创建generator_config表失败: %v", err)
	}

	return nil
}

// CloseDatabase 关闭数据库连接
func CloseDatabase() error {
	if db != nil {
		return db.Close()
	}
	return nil
}

// SaveDatabaseConfig 保存数据库配置
func SaveDatabaseConfig(config *DatabaseConfig, isUpdate bool) error {
	// 加密密码
	encryptedPassword, err := Encrypt(config.Password)
	if err != nil {
		return fmt.Errorf("加密密码失败: %v", err)
	}

	// 创建用于存储的副本，使用加密后的密码
	configToSave := *config
	configToSave.Password = encryptedPassword

	jsonData, err := json.Marshal(configToSave)
	if err != nil {
		return fmt.Errorf("序列化配置失败: %v", err)
	}

	if isUpdate {
		_, err = db.Exec("UPDATE dbs SET name = ?, value = ? WHERE id = ?",
			configToSave.Name, string(jsonData), configToSave.ID)
	} else {
		// 检查名称是否已存在
		var count int
		err = db.QueryRow("SELECT COUNT(*) FROM dbs WHERE name = ?", configToSave.Name).Scan(&count)
		if err != nil {
			return fmt.Errorf("检查配置名称失败: %v", err)
		}
		if count > 0 {
			return fmt.Errorf("配置名称已存在: %s", configToSave.Name)
		}

		_, err = db.Exec("INSERT INTO dbs (name, value) VALUES (?, ?)",
			configToSave.Name, string(jsonData))
	}

	return err
}

// LoadDatabaseConfigs 加载所有数据库配置
func LoadDatabaseConfigs() ([]*DatabaseConfig, error) {
	rows, err := db.Query("SELECT id, value FROM dbs")
	if err != nil {
		return nil, fmt.Errorf("查询数据库配置失败: %v", err)
	}
	defer rows.Close()

	var configs []*DatabaseConfig
	for rows.Next() {
		var id int
		var value string
		if err := rows.Scan(&id, &value); err != nil {
			return nil, fmt.Errorf("读取配置数据失败: %v", err)
		}

		var config DatabaseConfig
		if err := json.Unmarshal([]byte(value), &config); err != nil {
			return nil, fmt.Errorf("反序列化配置失败: %v", err)
		}
		config.ID = id

		// 解密密码
		decryptedPassword, err := Decrypt(config.Password)
		if err != nil {
			return nil, fmt.Errorf("解密密码失败: %v", err)
		}
		config.Password = decryptedPassword

		configs = append(configs, &config)
	}

	return configs, nil
}

// DeleteDatabaseConfig 删除数据库配置
func DeleteDatabaseConfig(id int) error {
	_, err := db.Exec("DELETE FROM dbs WHERE id = ?", id)
	return err
}

// SaveGeneratorConfig 保存代码生成配置
func SaveGeneratorConfig(config *GeneratorConfig) error {
	jsonData, err := json.Marshal(config)
	if err != nil {
		return fmt.Errorf("序列化配置失败: %v", err)
	}

	// 先删除同名配置再插入
	_, _ = db.Exec("DELETE FROM generator_config WHERE name = ?", config.Name)
	_, err = db.Exec("INSERT INTO generator_config (name, value) VALUES (?, ?)",
		config.Name, string(jsonData))

	return err
}

// LoadGeneratorConfigs 加载所有代码生成配置
func LoadGeneratorConfigs() ([]*GeneratorConfig, error) {
	rows, err := db.Query("SELECT value FROM generator_config")
	if err != nil {
		return nil, fmt.Errorf("查询生成配置失败: %v", err)
	}
	defer rows.Close()

	var configs []*GeneratorConfig
	for rows.Next() {
		var value string
		if err := rows.Scan(&value); err != nil {
			return nil, fmt.Errorf("读取配置数据失败: %v", err)
		}

		var config GeneratorConfig
		if err := json.Unmarshal([]byte(value), &config); err != nil {
			return nil, fmt.Errorf("反序列化配置失败: %v", err)
		}
		configs = append(configs, &config)
	}

	return configs, nil
}

// LoadGeneratorConfigByName 根据名称加载代码生成配置
func LoadGeneratorConfigByName(name string) (*GeneratorConfig, error) {
	var value string
	err := db.QueryRow("SELECT value FROM generator_config WHERE name = ?", name).Scan(&value)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("查询配置失败: %v", err)
	}

	var config GeneratorConfig
	if err := json.Unmarshal([]byte(value), &config); err != nil {
		return nil, fmt.Errorf("反序列化配置失败: %v", err)
	}

	return &config, nil
}

// DeleteGeneratorConfig 删除代码生成配置
func DeleteGeneratorConfig(name string) error {
	_, err := db.Exec("DELETE FROM generator_config WHERE name = ?", name)
	return err
}
