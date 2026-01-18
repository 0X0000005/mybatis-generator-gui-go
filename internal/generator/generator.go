package generator

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/yourusername/mybatis-generator-gui-go/internal/config"
	"github.com/yourusername/mybatis-generator-gui-go/internal/database"
	"github.com/yourusername/mybatis-generator-gui-go/internal/utils"
)

// Generator 代码生成器
type Generator struct {
	config    *config.GeneratorConfig
	dbConfig  *config.DatabaseConfig
	connector *database.Connector
}

// NewGenerator 创建新的代码生成器
func NewGenerator(cfg *config.GeneratorConfig, dbCfg *config.DatabaseConfig) *Generator {
	return &Generator{
		config:    cfg,
		dbConfig:  dbCfg,
		connector: database.NewConnector(dbCfg),
	}
}

// Generate 生成代码
func (g *Generator) Generate() ([]string, error) {
	var generatedFiles []string

	// 连接数据库
	if err := g.connector.Connect(); err != nil {
		return nil, fmt.Errorf("连接数据库失败: %v", err)
	}
	defer g.connector.Close()

	// 获取表列信息
	columns, err := g.connector.GetTableColumns(g.config.TableName)
	if err != nil {
		return nil, fmt.Errorf("获取表列信息失败: %v", err)
	}

	// 获取表注释
	tableComment, err := g.connector.GetTableComment(g.config.TableName)
	if err != nil {
		return nil, fmt.Errorf("获取表注释失败: %v", err)
	}

	// 生成Model类
	modelFile, err := g.generateModel(columns, tableComment)
	if err != nil {
		return nil, fmt.Errorf("生成Model失败: %v", err)
	}
	generatedFiles = append(generatedFiles, modelFile)

	// 生成Mapper接口
	mapperFile, err := g.generateMapper(columns)
	if err != nil {
		return nil, fmt.Errorf("生成Mapper接口失败: %v", err)
	}
	generatedFiles = append(generatedFiles, mapperFile)

	// 生成Mapper XML
	xmlFile, err := g.generateMapperXML(columns)
	if err != nil {
		return nil, fmt.Errorf("生成Mapper XML失败: %v", err)
	}
	if xmlFile != "" {
		generatedFiles = append(generatedFiles, xmlFile)
	}

	return generatedFiles, nil
}

// generateModel 生成Java Model类
func (g *Generator) generateModel(columns []*database.TableColumn, tableComment string) (string, error) {
	log.Printf("[Generator] 开始生成Model - 配置: Package=%s, UseLombok=%v, UseJsonProperty=%v, JsonPropertyUpperCase=%v",
		g.config.ModelPackage, g.config.UseLombokPlugin, g.config.UseJsonProperty, g.config.JsonPropertyUpperCase)

	// 准备模板数据
	data := g.prepareModelData(columns, tableComment)
	log.Printf("[Generator] 模板数据准备完成 - ClassName=%s, Package=%s, UseJsonProperty=%v, JsonPropertyUpperCase=%v, Fields=%d",
		data.ClassName, data.Package, data.UseJsonProperty, data.JsonPropertyUpperCase, len(data.Fields))

	// 选择模板
	var tmpl *template.Template
	var err error
	if g.config.UseLombokPlugin {
		log.Printf("[Generator] 使用Lombok模板")
		tmpl, err = template.New("model").Funcs(TemplateFuncs).Parse(modelLombokTemplate)
	} else {
		log.Printf("[Generator] 使用标准模板")
		tmpl, err = template.New("model").Funcs(TemplateFuncs).Parse(modelTemplate)
	}
	if err != nil {
		return "", fmt.Errorf("解析模板失败: %v", err)
	}

	// 生成文件路径
	filePath := g.getModelFilePath()
	log.Printf("[Generator] 生成文件路径: %s", filePath)
	if err := os.MkdirAll(filepath.Dir(filePath), 0755); err != nil {
		return "", fmt.Errorf("创建目录失败: %v", err)
	}

	// 创建文件
	file, err := os.Create(filePath)
	if err != nil {
		return "", fmt.Errorf("创建文件失败: %v", err)
	}
	defer file.Close()

	// 执行模板
	if err := tmpl.Execute(file, data); err != nil {
		return "", fmt.Errorf("执行模板失败: %v", err)
	}

	log.Printf("[Generator] Model生成成功: %s", filePath)
	return filePath, nil
}

// generateMapper 生成Java Mapper接口
func (g *Generator) generateMapper(columns []*database.TableColumn) (string, error) {
	// 准备模板数据
	data := g.prepareMapperData(columns)

	// 解析模板
	tmpl, err := template.New("mapper").Parse(mapperTemplate)
	if err != nil {
		return "", fmt.Errorf("解析模板失败: %v", err)
	}

	// 生成文件路径
	filePath := g.getMapperFilePath()
	if err := os.MkdirAll(filepath.Dir(filePath), 0755); err != nil {
		return "", fmt.Errorf("创建目录失败: %v", err)
	}

	// 创建文件
	file, err := os.Create(filePath)
	if err != nil {
		return "", fmt.Errorf("创建文件失败: %v", err)
	}
	defer file.Close()

	// 执行模板
	if err := tmpl.Execute(file, data); err != nil {
		return "", fmt.Errorf("执行模板失败: %v", err)
	}

	return filePath, nil
}

// generateMapperXML 生成MyBatis Mapper XML
func (g *Generator) generateMapperXML(columns []*database.TableColumn) (string, error) {
	// 准备模板数据
	data := g.prepareMapperXMLData(columns)

	// 解析模板
	tmpl, err := template.New("mapperXML").Parse(mapperXMLTemplate)
	if err != nil {
		return "", fmt.Errorf("解析模板失败: %v", err)
	}

	// 生成文件路径
	filePath := g.getMapperXMLFilePath()
	if err := os.MkdirAll(filepath.Dir(filePath), 0755); err != nil {
		return "", fmt.Errorf("创建目录失败: %v", err)
	}

	// 如果不覆盖且文件已存在，则跳过
	if !g.config.OverrideXML {
		if _, err := os.Stat(filePath); err == nil {
			return "", nil // 文件存在，跳过
		}
	}

	// 创建文件
	file, err := os.Create(filePath)
	if err != nil {
		return "", fmt.Errorf("创建文件失败: %v", err)
	}
	defer file.Close()

	// 执行模板
	if err := tmpl.Execute(file, data); err != nil {
		return "", fmt.Errorf("执行模板失败: %v", err)
	}

	return filePath, nil
}

// ModelField Model字段
type ModelField struct {
	FieldName    string // Java字段名
	FieldType    string // Java类型
	ColumnName   string // 数据库列名
	Comment      string // 注释
	IsPrimaryKey bool   // 是否主键
}

// ModelData Model模板数据
type ModelData struct {
	Package               string
	ClassName             string
	TableComment          string
	Fields                []*ModelField
	Imports               []string
	UseJsonProperty       bool
	JsonPropertyUpperCase bool
}

// prepareModelData 准备Model模板数据
func (g *Generator) prepareModelData(columns []*database.TableColumn, tableComment string) *ModelData {
	data := &ModelData{
		Package:               g.config.ModelPackage,
		ClassName:             g.config.DomainObjectName,
		TableComment:          tableComment,
		Fields:                make([]*ModelField, 0),
		UseJsonProperty:       g.config.UseJsonProperty,
		JsonPropertyUpperCase: g.config.JsonPropertyUpperCase,
	}

	imports := make(map[string]bool)

	for _, col := range columns {
		// 获取字段名
		fieldName := g.getFieldName(col.ColumnName)

		// 获取Java类型
		javaType := database.GetJavaType(g.dbConfig.DbType, col.DataType, g.config.JSR310Support)

		// 添加导入
		g.addImport(imports, javaType, g.config.JSR310Support)

		field := &ModelField{
			FieldName:    fieldName,
			FieldType:    javaType,
			ColumnName:   col.ColumnName,
			Comment:      col.ColumnComment,
			IsPrimaryKey: col.ColumnKey == "PRI",
		}
		data.Fields = append(data.Fields, field)
	}

	// 转换imports为切片
	for imp := range imports {
		data.Imports = append(data.Imports, imp)
	}

	return data
}

// getFieldName 获取字段名
func (g *Generator) getFieldName(columnName string) string {
	if g.config.UseActualColumnNames {
		return columnName
	}
	return utils.DBStringToCamelCase(columnName)
}

// addImport 添加导入语句
func (g *Generator) addImport(imports map[string]bool, javaType string, useJSR310 bool) {
	switch javaType {
	case "Date":
		imports["java.util.Date"] = true
	case "BigDecimal":
		imports["java.math.BigDecimal"] = true
	case "LocalDate", "LocalDateTime", "LocalTime":
		if useJSR310 {
			imports["java.time."+javaType] = true
		}
	}
}

// getModelFilePath 获取Model文件路径
func (g *Generator) getModelFilePath() string {
	packagePath := strings.ReplaceAll(g.config.ModelPackage, ".", string(filepath.Separator))
	return filepath.Join(
		g.config.ProjectFolder,
		g.config.ModelPackageTargetFolder,
		packagePath,
		g.config.DomainObjectName+".java",
	)
}

// getMapperFilePath 获取Mapper文件路径
func (g *Generator) getMapperFilePath() string {
	packagePath := strings.ReplaceAll(g.config.DaoPackage, ".", string(filepath.Separator))
	mapperName := g.config.MapperName
	if mapperName == "" {
		mapperName = g.config.DomainObjectName + "Mapper"
	}
	return filepath.Join(
		g.config.ProjectFolder,
		g.config.DaoTargetFolder,
		packagePath,
		mapperName+".java",
	)
}

// getMapperXMLFilePath 获取Mapper XML文件路径
func (g *Generator) getMapperXMLFilePath() string {
	packagePath := strings.ReplaceAll(g.config.MappingXMLPackage, ".", string(filepath.Separator))
	mapperName := g.config.MapperName
	if mapperName == "" {
		mapperName = g.config.DomainObjectName + "Mapper"
	}
	return filepath.Join(
		g.config.ProjectFolder,
		g.config.MappingXMLTargetFolder,
		packagePath,
		mapperName+".xml",
	)
}

// MapperData Mapper模板数据
type MapperData struct {
	Package        string
	MapperName     string
	ModelPackage   string
	ModelName      string
	PrimaryKey     *ModelField
	UseExample     bool
	OffsetLimit    bool
	UseBatchInsert bool
	UseBatchUpdate bool
}

// prepareMapperData 准备Mapper模板数据
func (g *Generator) prepareMapperData(columns []*database.TableColumn) *MapperData {
	data := &MapperData{
		Package:        g.config.DaoPackage,
		MapperName:     g.config.MapperName,
		ModelPackage:   g.config.ModelPackage,
		ModelName:      g.config.DomainObjectName,
		UseExample:     g.config.UseExample,
		OffsetLimit:    g.config.OffsetLimit,
		UseBatchInsert: g.config.UseBatchInsert,
		UseBatchUpdate: g.config.UseBatchUpdate,
	}

	if data.MapperName == "" {
		data.MapperName = g.config.DomainObjectName + "Mapper"
	}

	// 查找主键
	for _, col := range columns {
		if col.ColumnKey == "PRI" {
			fieldName := g.getFieldName(col.ColumnName)
			javaType := database.GetJavaType(g.dbConfig.DbType, col.DataType, g.config.JSR310Support)
			data.PrimaryKey = &ModelField{
				FieldName:  fieldName,
				FieldType:  javaType,
				ColumnName: col.ColumnName,
			}
			break
		}
	}

	return data
}

// MapperXMLData Mapper XML模板数据
type MapperXMLData struct {
	Namespace        string
	ModelType        string
	TableName        string
	Columns          []*ColumnMapping
	PrimaryKey       *ColumnMapping
	OffsetLimit      bool
	UseGeneratedKeys bool
	GenerateKeys     string
	UseBatchInsert   bool
	UseBatchUpdate   bool
}

// ColumnMapping 列映射
type ColumnMapping struct {
	ColumnName string
	FieldName  string
	JdbcType   string
	JavaType   string
}

// prepareMapperXMLData 准备Mapper XML模板数据
func (g *Generator) prepareMapperXMLData(columns []*database.TableColumn) *MapperXMLData {
	mapperName := g.config.MapperName
	if mapperName == "" {
		mapperName = g.config.DomainObjectName + "Mapper"
	}

	data := &MapperXMLData{
		Namespace:        g.config.DaoPackage + "." + mapperName,
		ModelType:        g.config.ModelPackage + "." + g.config.DomainObjectName,
		TableName:        g.config.TableName,
		Columns:          make([]*ColumnMapping, 0),
		OffsetLimit:      g.config.OffsetLimit,
		UseGeneratedKeys: g.config.GenerateKeys != "",
		GenerateKeys:     g.config.GenerateKeys,
		UseBatchInsert:   g.config.UseBatchInsert,
		UseBatchUpdate:   g.config.UseBatchUpdate,
	}

	for _, col := range columns {
		fieldName := g.getFieldName(col.ColumnName)
		javaType := database.GetJavaType(g.dbConfig.DbType, col.DataType, g.config.JSR310Support)
		jdbcType := database.GetJdbcType(g.dbConfig.DbType, col.DataType)

		mapping := &ColumnMapping{
			ColumnName: col.ColumnName,
			FieldName:  fieldName,
			JdbcType:   jdbcType,
			JavaType:   javaType,
		}

		data.Columns = append(data.Columns, mapping)

		// 记录主键
		if col.ColumnKey == "PRI" {
			data.PrimaryKey = mapping
		}
	}

	return data
}
