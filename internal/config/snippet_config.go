package config

// SnippetOperation 操作类型
type SnippetOperation string

const (
	OperationSelect SnippetOperation = "select"
	OperationInsert SnippetOperation = "insert"
	OperationDelete SnippetOperation = "delete"
	OperationUpdate SnippetOperation = "update"
)

// SnippetField 片段字段配置
type SnippetField struct {
	ColumnName string `json:"columnName"` // 数据库列名
	FieldName  string `json:"fieldName"`  // Java字段名
	JdbcType   string `json:"jdbcType"`   // JDBC类型
	JavaType   string `json:"javaType"`   // Java类型
	// WHERE条件运算符，仅用于 WhereFields：
	// =, !=, >, <, >=, <=, LIKE, IS NULL, IS NOT NULL
	// 空值默认为 "="
	Operator string `json:"operator"`
}

// OrderByField 排序字段配置
type OrderByField struct {
	ColumnName string `json:"columnName"` // 数据库列名
	FieldName  string `json:"fieldName"`  // Java字段名
	JdbcType   string `json:"jdbcType"`   // JDBC类型
	Direction  string `json:"direction"`  // ASC / DESC
}

// SnippetConfig 单个自定义片段配置
type SnippetConfig struct {
	MethodName    string           `json:"methodName"`    // 方法名（用户可自定义，空则自动生成）
	Operation     SnippetOperation `json:"operation"`     // 操作类型
	IsBatch       bool             `json:"isBatch"`       // 是否批量
	WhereLogic    string           `json:"whereLogic"`    // WHERE条件之间的逻辑：AND / OR（默认AND）
	SelectFields  []SnippetField   `json:"selectFields"`  // 查询：SELECT列（顺序有效）
	WhereFields   []SnippetField   `json:"whereFields"`   // 查询/删除/更新：WHERE条件列（含运算符）
	OrderByFields []OrderByField   `json:"orderByFields"` // 查询：ORDER BY列（顺序有效）
	InsertFields  []SnippetField   `json:"insertFields"`  // 新增：INSERT列（顺序有效）
	SetFields     []SnippetField   `json:"setFields"`     // 更新：SET列（顺序有效）
}

