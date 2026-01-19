package config

// GeneratorConfig 代码生成配置
type GeneratorConfig struct {
	Name                     string `json:"name"`                     // 配置名称
	ProjectFolder            string `json:"projectFolder"`            // 项目根目录
	ModelPackage             string `json:"modelPackage"`             // Model包名
	ModelPackageTargetFolder string `json:"modelPackageTargetFolder"` // Model目标文件夹
	DaoPackage               string `json:"daoPackage"`               // DAO包名
	DaoTargetFolder          string `json:"daoTargetFolder"`          // DAO目标文件夹
	MapperName               string `json:"mapperName"`               // Mapper名称
	MappingXMLPackage        string `json:"mappingXMLPackage"`        // Mapper XML包名
	MappingXMLTargetFolder   string `json:"mappingXMLTargetFolder"`   // Mapper XML目标文件夹
	TableName                string `json:"tableName"`                // 表名
	DomainObjectName         string `json:"domainObjectName"`         // 实体类名
	GenerateKeys             string `json:"generateKeys"`             // 主键字段名
	Encoding                 string `json:"encoding"`                 // 文件编码

	// 生成选项
	OffsetLimit                bool `json:"offsetLimit"`                // 是否生成分页查询
	Comment                    bool `json:"comment"`                    // 是否生成注释
	OverrideXML                bool `json:"overrideXML"`                // 是否覆盖XML
	NeedToStringHashcodeEquals bool `json:"needToStringHashcodeEquals"` // 是否生成toString/hashCode/equals
	UseLombokPlugin            bool `json:"useLombokPlugin"`            // 是否使用Lombok
	UseTableNameAlias          bool `json:"useTableNameAlias"`          // 是否使用表名别名
	NeedForUpdate              bool `json:"needForUpdate"`              // 是否生成for update
	AnnotationDAO              bool `json:"annotationDAO"`              // 是否使用@Repository注解
	Annotation                 bool `json:"annotation"`                 // 是否使用注解
	UseActualColumnNames       bool `json:"useActualColumnNames"`       // 是否使用实际列名
	UseExample                 bool `json:"useExample"`                 // 是否生成Example
	UseDAOExtendStyle          bool `json:"useDAOExtendStyle"`          // 是否使用DAO扩展风格
	UseSchemaPrefix            bool `json:"useSchemaPrefix"`            // 是否使用Schema前缀
	JSR310Support              bool `json:"jsr310Support"`              // 是否支持JSR310日期类型
	UseJsonProperty            bool `json:"useJsonProperty"`            // 是否使用@JsonProperty注解
	JsonPropertyUpperCase      bool `json:"jsonPropertyUpperCase"`      // @JsonProperty首字母大写
	UseBatchInsert             bool `json:"useBatchInsert"`             // 是否生成批量插入
	UseBatchUpdate             bool `json:"useBatchUpdate"`             // 是否生成批量更新

	// 列定制
	IgnoredColumns  []string         `json:"ignoredColumns"`  // 忽略的列名列表
	ColumnOverrides []ColumnOverride `json:"columnOverrides"` // 列覆盖配置
}

// ColumnOverride 列覆盖配置
type ColumnOverride struct {
	ColumnName   string `json:"columnName"`   // 数据库列名
	PropertyName string `json:"propertyName"` // Java属性名（可选）
	JavaType     string `json:"javaType"`     // Java类型（可选）
}
