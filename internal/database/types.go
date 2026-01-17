package database

// TableColumn 表列信息
type TableColumn struct {
	ColumnName    string `json:"columnName"`    // 列名
	DataType      string `json:"dataType"`      // 数据类型
	ColumnComment string `json:"columnComment"` // 列注释
	IsNullable    bool   `json:"isNullable"`    // 是否可为空
	ColumnKey     string `json:"columnKey"`     // 键类型 (PRI, UNI, MUL)
	Extra         string `json:"extra"`         // 额外信息 (auto_increment等)
}

// TableInfo 表信息
type TableInfo struct {
	TableName    string         `json:"tableName"`    // 表名
	TableComment string         `json:"tableComment"` // 表注释
	Columns      []*TableColumn `json:"columns"`      // 列信息
}
