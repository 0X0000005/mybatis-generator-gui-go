package database

// JavaTypeMapping SQL类型到Java类型的映射
var JavaTypeMapping = map[string]map[string]string{
	"MySQL": {
		"varchar":    "String",
		"char":       "String",
		"text":       "String",
		"mediumtext": "String",
		"longtext":   "String",
		"tinytext":   "String",
		"int":        "Integer",
		"tinyint":    "Integer",
		"smallint":   "Integer",
		"mediumint":  "Integer",
		"bigint":     "Long",
		"decimal":    "BigDecimal",
		"numeric":    "BigDecimal",
		"float":      "Float",
		"double":     "Double",
		"date":       "Date",
		"datetime":   "Date",
		"timestamp":  "Date",
		"time":       "Date",
		"year":       "Date",
		"bit":        "Boolean",
		"boolean":    "Boolean",
		"blob":       "byte[]",
		"mediumblob": "byte[]",
		"longblob":   "byte[]",
		"binary":     "byte[]",
		"varbinary":  "byte[]",
	},
	"PostgreSQL": {
		"varchar":           "String",
		"character varying": "String",
		"char":              "String",
		"character":         "String",
		"text":              "String",
		"integer":           "Integer",
		"int":               "Integer",
		"int4":              "Integer",
		"smallint":          "Integer",
		"int2":              "Integer",
		"bigint":            "Long",
		"int8":              "Long",
		"numeric":           "BigDecimal",
		"decimal":           "BigDecimal",
		"real":              "Float",
		"float4":            "Float",
		"double precision":  "Double",
		"float8":            "Double",
		"date":              "Date",
		"timestamp":         "Date",
		"time":              "Date",
		"boolean":           "Boolean",
		"bool":              "Boolean",
		"bytea":             "byte[]",
		"uuid":              "String",
		"json":              "String",
		"jsonb":             "String",
	},
	"Oracle": {
		"varchar2":      "String",
		"nvarchar2":     "String",
		"char":          "String",
		"nchar":         "String",
		"clob":          "String",
		"nclob":         "String",
		"long":          "String",
		"number":        "BigDecimal",
		"integer":       "Integer",
		"int":           "Integer",
		"smallint":      "Integer",
		"float":         "Double",
		"binary_float":  "Float",
		"binary_double": "Double",
		"date":          "Date",
		"timestamp":     "Date",
		"blob":          "byte[]",
		"raw":           "byte[]",
		"long raw":      "byte[]",
	},
}

// JSR310TypeMapping SQL类型到Java JSR310日期类型的映射
var JSR310TypeMapping = map[string]map[string]string{
	"MySQL": {
		"date":      "LocalDate",
		"datetime":  "LocalDateTime",
		"timestamp": "LocalDateTime",
		"time":      "LocalTime",
	},
	"PostgreSQL": {
		"date":      "LocalDate",
		"timestamp": "LocalDateTime",
		"time":      "LocalTime",
	},
	"Oracle": {
		"date":      "LocalDateTime",
		"timestamp": "LocalDateTime",
	},
}

// MyBatisJdbcTypeMapping SQL类型到MyBatis JDBC类型的映射
var MyBatisJdbcTypeMapping = map[string]map[string]string{
	"MySQL": {
		"varchar":    "VARCHAR",
		"char":       "CHAR",
		"text":       "LONGVARCHAR",
		"mediumtext": "LONGVARCHAR",
		"longtext":   "LONGVARCHAR",
		"tinytext":   "VARCHAR",
		"int":        "INTEGER",
		"tinyint":    "TINYINT",
		"smallint":   "SMALLINT",
		"mediumint":  "INTEGER",
		"bigint":     "BIGINT",
		"decimal":    "DECIMAL",
		"numeric":    "NUMERIC",
		"float":      "FLOAT",
		"double":     "DOUBLE",
		"date":       "DATE",
		"datetime":   "TIMESTAMP",
		"timestamp":  "TIMESTAMP",
		"time":       "TIME",
		"year":       "INTEGER",
		"bit":        "BIT",
		"boolean":    "BOOLEAN",
		"blob":       "BLOB",
		"mediumblob": "BLOB",
		"longblob":   "BLOB",
		"binary":     "BINARY",
		"varbinary":  "VARBINARY",
	},
	"PostgreSQL": {
		"varchar":           "VARCHAR",
		"character varying": "VARCHAR",
		"char":              "CHAR",
		"character":         "CHAR",
		"text":              "LONGVARCHAR",
		"integer":           "INTEGER",
		"int":               "INTEGER",
		"int4":              "INTEGER",
		"smallint":          "SMALLINT",
		"int2":              "SMALLINT",
		"bigint":            "BIGINT",
		"int8":              "BIGINT",
		"numeric":           "NUMERIC",
		"decimal":           "DECIMAL",
		"real":              "REAL",
		"float4":            "REAL",
		"double precision":  "DOUBLE",
		"float8":            "DOUBLE",
		"date":              "DATE",
		"timestamp":         "TIMESTAMP",
		"time":              "TIME",
		"boolean":           "BOOLEAN",
		"bool":              "BOOLEAN",
		"bytea":             "BLOB",
		"uuid":              "VARCHAR",
		"json":              "LONGVARCHAR",
		"jsonb":             "LONGVARCHAR",
	},
	"Oracle": {
		"varchar2":      "VARCHAR",
		"nvarchar2":     "NVARCHAR",
		"char":          "CHAR",
		"nchar":         "NCHAR",
		"clob":          "CLOB",
		"nclob":         "NCLOB",
		"long":          "LONGVARCHAR",
		"number":        "NUMERIC",
		"integer":       "INTEGER",
		"int":           "INTEGER",
		"smallint":      "SMALLINT",
		"float":         "FLOAT",
		"binary_float":  "FLOAT",
		"binary_double": "DOUBLE",
		"date":          "DATE",
		"timestamp":     "TIMESTAMP",
		"blob":          "BLOB",
		"raw":           "VARBINARY",
		"long raw":      "LONGVARBINARY",
	},
}

// GetJavaType 获取Java类型
func GetJavaType(dbType, sqlType string, useJSR310 bool) string {
	sqlType = normalizeType(sqlType)

	// 如果启用JSR310，优先使用JSR310类型
	if useJSR310 {
		if mapping, ok := JSR310TypeMapping[dbType]; ok {
			if javaType, ok := mapping[sqlType]; ok {
				return javaType
			}
		}
	}

	// 使用标准Java类型
	if mapping, ok := JavaTypeMapping[dbType]; ok {
		if javaType, ok := mapping[sqlType]; ok {
			return javaType
		}
	}

	// 默认返回String
	return "String"
}

// GetJdbcType 获取MyBatis JDBC类型
func GetJdbcType(dbType, sqlType string) string {
	sqlType = normalizeType(sqlType)

	if mapping, ok := MyBatisJdbcTypeMapping[dbType]; ok {
		if jdbcType, ok := mapping[sqlType]; ok {
			return jdbcType
		}
	}

	// 默认返回VARCHAR
	return "VARCHAR"
}

// normalizeType 标准化SQL类型（转为小写，去除参数部分）
func normalizeType(sqlType string) string {
	// 转为小写
	sqlType = toLower(sqlType)

	// 去除括号及其内容，例如 varchar(255) -> varchar
	for i, c := range sqlType {
		if c == '(' {
			return sqlType[:i]
		}
	}

	return sqlType
}

func toLower(s string) string {
	b := make([]byte, len(s))
	for i := range s {
		c := s[i]
		if 'A' <= c && c <= 'Z' {
			c += 'a' - 'A'
		}
		b[i] = c
	}
	return string(b)
}
