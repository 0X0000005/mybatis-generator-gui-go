package generator

import (
	"bytes"
	"fmt"
	"sort"
	"strings"
	"text/template"

	"github.com/yourusername/mybatis-generator-gui-go/internal/config"
)

// SnippetResult 片段生成结果
type SnippetResult struct {
	JavaCode string   // Mapper接口方法声明代码（使用简单类名）
	XMLCode  string   // XML SQL片段代码
	Imports  []string // 需要的 import（如 ["java.util.List", "org.apache.ibatis.annotations.Param"]）
}

// GenerateSnippet 生成自定义MyBatis片段
func GenerateSnippet(cfg *config.SnippetConfig, mapperName, modelType, tableName string) (*SnippetResult, error) {
	switch cfg.Operation {
	case config.OperationSelect:
		return generateSelectSnippet(cfg, mapperName, modelType, tableName)
	case config.OperationInsert:
		return generateInsertSnippet(cfg, mapperName, modelType, tableName)
	case config.OperationDelete:
		return generateDeleteSnippet(cfg, mapperName, modelType, tableName)
	case config.OperationUpdate:
		return generateUpdateSnippet(cfg, mapperName, modelType, tableName)
	default:
		return nil, fmt.Errorf("未知操作类型: %s", cfg.Operation)
	}
}

// AppendSnippetToJava 将片段追加到Mapper.java文件内容中（在最后一个 } 前）
func AppendSnippetToJava(javaContent, javaCode string) string {
	lastBrace := strings.LastIndex(javaContent, "}")
	if lastBrace < 0 {
		return javaContent + "\n" + javaCode
	}
	return javaContent[:lastBrace] + "\n" + javaCode + "\n" + javaContent[lastBrace:]
}

// AppendImportsToJava 将缺失的 import 注入到 Mapper.java 的 import 块
func AppendImportsToJava(javaContent string, imports []string) string {
	for _, imp := range imports {
		line := "import " + imp + ";"
		if strings.Contains(javaContent, line) {
			continue
		}
		// 在最后一个 import 语句后插入
		lastImport := strings.LastIndex(javaContent, "import ")
		if lastImport >= 0 {
			newlineAfter := strings.Index(javaContent[lastImport:], "\n")
			if newlineAfter >= 0 {
				insertAt := lastImport + newlineAfter + 1
				javaContent = javaContent[:insertAt] + line + "\n" + javaContent[insertAt:]
				continue
			}
		}
		// 如果没有 import，就在 package 语句后插入
		pkgEnd := strings.Index(javaContent, ";")
		if pkgEnd >= 0 {
			newlineAfter := strings.Index(javaContent[pkgEnd:], "\n")
			if newlineAfter >= 0 {
				insertAt := pkgEnd + newlineAfter + 1
				javaContent = javaContent[:insertAt] + "\n" + line + "\n" + javaContent[insertAt:]
			}
		}
	}
	return javaContent
}

// AppendSnippetToXML 将片段追加到Mapper.xml文件内容中（在 </mapper> 前）
func AppendSnippetToXML(xmlContent, xmlCode string) string {
	closeTag := "</mapper>"
	idx := strings.LastIndex(xmlContent, closeTag)
	if idx < 0 {
		return xmlContent + "\n" + xmlCode
	}
	return xmlContent[:idx] + "\n" + xmlCode + "\n" + xmlContent[idx:]
}

// -----------------------------------------------------------------------
// 内部生成函数
// -----------------------------------------------------------------------

func generateSelectSnippet(cfg *config.SnippetConfig, mapperName, modelType, tableName string) (*SnippetResult, error) {
	methodName := cfg.MethodName
	if methodName == "" {
		methodName = buildSelectMethodName(cfg)
	}
	simpleModel := lastPart(modelType)

	// ---- Java 代码 ----
	var javaBuilder strings.Builder
	javaBuilder.WriteString("    /**\n")
	if cfg.IsBatch {
		javaBuilder.WriteString(fmt.Sprintf("     * 批量查询（IN查询）- %s\n", methodName))
	} else {
		javaBuilder.WriteString(fmt.Sprintf("     * 自定义查询 - %s\n", methodName))
	}
	javaBuilder.WriteString("     */\n")

	if cfg.IsBatch {
		if len(cfg.WhereFields) == 0 {
			return nil, fmt.Errorf("批量查询需要至少一个WHERE字段")
		}
		inField := cfg.WhereFields[0]
		javaBuilder.WriteString(fmt.Sprintf(
			"    List<%s> %s(@Param(\"list\") List<%s> list);\n",
			simpleModel, methodName, inField.JavaType,
		))
	} else {
		params := buildJavaParams(cfg)
		javaBuilder.WriteString(fmt.Sprintf(
			"    List<%s> %s(%s);\n",
			simpleModel, methodName, params,
		))
	}

	// 预计算 SQL 片段
	selectSQL := buildSelectSQL(cfg.SelectFields)
	whereSQL := buildWhereSQL(cfg.WhereFields, cfg.WhereLogic)
	orderBySQL := buildOrderBySQL(cfg.OrderByFields)

	// ---- XML 代码 ----
	xmlCode, err := renderTemplate("selectSnippet", selectSnippetTemplate, map[string]interface{}{
		"MethodName":   methodName,
		"ModelType":    modelType,
		"TableName":    tableName,
		"SelectSQL":    selectSQL,
		"WhereSQL":     whereSQL,
		"OrderBySQL":   orderBySQL,
		"IsBatch":      cfg.IsBatch,
		"InField":      firstWhereField(cfg.WhereFields),
		"HasLimit":     cfg.HasLimit,
		"IsLimitFixed": cfg.IsLimitFixed,
		"LimitValue":   cfg.LimitValue,
	})
	if err != nil {
		return nil, err
	}

	return &SnippetResult{
		JavaCode: javaBuilder.String(),
		XMLCode:  xmlCode,
		Imports:  collectSnippetImports(cfg),
	}, nil
}

func generateInsertSnippet(cfg *config.SnippetConfig, mapperName, modelType, tableName string) (*SnippetResult, error) {
	methodName := cfg.MethodName
	if methodName == "" {
		if cfg.IsBatch {
			methodName = "insertBatchByFields"
		} else {
			methodName = "insertByFields"
		}
	}
	simpleModel := lastPart(modelType)

	// ---- Java 代码 ----
	var javaBuilder strings.Builder
	javaBuilder.WriteString("    /**\n")
	if cfg.IsBatch {
		javaBuilder.WriteString(fmt.Sprintf("     * 批量插入 - %s\n", methodName))
	} else {
		javaBuilder.WriteString(fmt.Sprintf("     * 自定义插入 - %s\n", methodName))
	}
	javaBuilder.WriteString("     */\n")
	if cfg.IsBatch {
		javaBuilder.WriteString(fmt.Sprintf(
			"    int %s(@Param(\"list\") List<%s> list);\n",
			methodName, simpleModel,
		))
	} else {
		javaBuilder.WriteString(fmt.Sprintf("    int %s(%s record);\n", methodName, simpleModel))
	}

	// ---- XML 代码 ----
	xmlCode, err := renderTemplate("insertSnippet", insertSnippetTemplate, map[string]interface{}{
		"MethodName":   methodName,
		"ModelType":    modelType,
		"TableName":    tableName,
		"InsertFields": cfg.InsertFields,
		"IsBatch":      cfg.IsBatch,
	})
	if err != nil {
		return nil, err
	}

	return &SnippetResult{JavaCode: javaBuilder.String(), XMLCode: xmlCode, Imports: collectSnippetImports(cfg)}, nil
}

func generateDeleteSnippet(cfg *config.SnippetConfig, mapperName, modelType, tableName string) (*SnippetResult, error) {
	methodName := cfg.MethodName
	if methodName == "" {
		methodName = buildDeleteMethodName(cfg)
	}

	// ---- Java 代码 ----
	var javaBuilder strings.Builder
	javaBuilder.WriteString("    /**\n")
	if cfg.IsBatch {
		javaBuilder.WriteString(fmt.Sprintf("     * 批量删除（IN删除）- %s\n", methodName))
	} else {
		javaBuilder.WriteString(fmt.Sprintf("     * 自定义删除 - %s\n", methodName))
	}
	javaBuilder.WriteString("     */\n")

	if cfg.IsBatch {
		if len(cfg.WhereFields) == 0 {
			return nil, fmt.Errorf("批量删除需要至少一个WHERE字段")
		}
		inField := cfg.WhereFields[0]
		javaBuilder.WriteString(fmt.Sprintf(
			"    int %s(@Param(\"list\") List<%s> list);\n",
			methodName, inField.JavaType,
		))
	} else {
		params := buildJavaParams(cfg)
		javaBuilder.WriteString(fmt.Sprintf("    int %s(%s);\n", methodName, params))
	}

	whereSQL := buildWhereSQL(cfg.WhereFields, cfg.WhereLogic)

	// ---- XML 代码 ----
	xmlCode, err := renderTemplate("deleteSnippet", deleteSnippetTemplate, map[string]interface{}{
		"MethodName": methodName,
		"TableName":  tableName,
		"WhereSQL":   whereSQL,
		"IsBatch":    cfg.IsBatch,
		"InField":    firstWhereField(cfg.WhereFields),
	})
	if err != nil {
		return nil, err
	}

	return &SnippetResult{JavaCode: javaBuilder.String(), XMLCode: xmlCode, Imports: collectSnippetImports(cfg)}, nil
}

func generateUpdateSnippet(cfg *config.SnippetConfig, mapperName, modelType, tableName string) (*SnippetResult, error) {
	methodName := cfg.MethodName
	if methodName == "" {
		methodName = buildUpdateMethodName(cfg)
	}
	simpleModel := lastPart(modelType)

	// ---- Java 代码 ----
	var javaBuilder strings.Builder
	javaBuilder.WriteString("    /**\n")
	if cfg.IsBatch {
		javaBuilder.WriteString(fmt.Sprintf("     * 批量更新 - %s\n", methodName))
	} else {
		javaBuilder.WriteString(fmt.Sprintf("     * 自定义更新 - %s\n", methodName))
	}
	javaBuilder.WriteString("     */\n")
	if cfg.IsBatch {
		javaBuilder.WriteString(fmt.Sprintf(
			"    int %s(@org.apache.ibatis.annotations.Param(\"list\") java.util.List<%s> list);\n",
			methodName, simpleModel,
		))
	} else {
		javaBuilder.WriteString(fmt.Sprintf("    int %s(%s record);\n", methodName, simpleModel))
	}

	// 批量更新的WHERE固定用=和item前缀；单条用operator
	var whereSQL string
	if cfg.IsBatch {
		whereSQL = buildBatchWhereSQL(cfg.WhereFields)
	} else {
		whereSQL = buildWhereSQL(cfg.WhereFields, cfg.WhereLogic)
	}

	// ---- XML 代码 ----
	xmlCode, err := renderTemplate("updateSnippet", updateSnippetTemplate, map[string]interface{}{
		"MethodName": methodName,
		"ModelType":  modelType,
		"TableName":  tableName,
		"SetFields":  cfg.SetFields,
		"WhereSQL":   whereSQL,
		"IsBatch":    cfg.IsBatch,
	})
	if err != nil {
		return nil, err
	}

	return &SnippetResult{JavaCode: javaBuilder.String(), XMLCode: xmlCode, Imports: collectSnippetImports(cfg)}, nil
}

// -----------------------------------------------------------------------
// WHERE SQL 构建
// -----------------------------------------------------------------------

// buildWhereSQL 根据字段列表和逻辑（AND/OR）生成 WHERE 子句（不含 WHERE 关键字）
func buildWhereSQL(fields []config.SnippetField, logic string) string {
	if len(fields) == 0 {
		return ""
	}
	if logic == "" {
		logic = "AND"
	}
	clauses := make([]string, 0, len(fields))
	for _, f := range fields {
		clauses = append(clauses, buildWhereClause(f, false))
	}
	return strings.Join(clauses, "\n        "+logic+" ")
}

// buildBatchWhereSQL 批量操作的WHERE子句，使用 item. 前缀且运算符固定为 =
func buildBatchWhereSQL(fields []config.SnippetField) string {
	if len(fields) == 0 {
		return ""
	}
	clauses := make([]string, 0, len(fields))
	for _, f := range fields {
		clauses = append(clauses, fmt.Sprintf("%s = #{item.%s,jdbcType=%s}", f.ColumnName, f.FieldName, f.JdbcType))
	}
	return strings.Join(clauses, " AND ")
}

// buildWhereClause 将单个字段转成 WHERE 条件片段
func buildWhereClause(f config.SnippetField, batchMode bool) string {
	op := f.Operator
	if op == "" {
		op = "="
	}
	prefix := ""
	if batchMode {
		prefix = "item."
	}
	switch op {
	case "IS NULL":
		return f.ColumnName + " IS NULL"
	case "IS NOT NULL":
		return f.ColumnName + " IS NOT NULL"
	case "LIKE":
		if f.IsFixed {
			return fmt.Sprintf("%s LIKE '%s'", f.ColumnName, f.FixedValue)
		}
		return fmt.Sprintf("%s LIKE CONCAT('%%', #{%s%s,jdbcType=%s}, '%%')",
			f.ColumnName, prefix, f.FieldName, f.JdbcType)
	case "IN", "NOT IN":
		if f.IsFixed {
			return fmt.Sprintf("%s %s (%s)", f.ColumnName, op, f.FixedValue)
		}
		return fmt.Sprintf("%s %s\n        <foreach collection=\"%s\" item=\"item\" open=\"(\" separator=\",\" close=\")\">\n            #{item,jdbcType=%s}\n        </foreach>",
			f.ColumnName, op, f.FieldName, f.JdbcType)
	default:
		if f.IsFixed {
			// 固定值直接内嵌，字符串类型需加引号
			val := f.FixedValue
			if f.JdbcType == "VARCHAR" || f.JdbcType == "CHAR" || f.JdbcType == "TEXT" {
				val = "'" + val + "'"
			}
			return fmt.Sprintf("%s %s %s", f.ColumnName, op, val)
		}
		return fmt.Sprintf("%s %s #{%s%s,jdbcType=%s}",
			f.ColumnName, op, prefix, f.FieldName, f.JdbcType)
	}
}

// buildSelectSQL 构建 SELECT 字段部分（处理聚合函数和别名）
func buildSelectSQL(fields []config.SnippetField) string {
	if len(fields) == 0 {
		return "*"
	}
	parts := make([]string, len(fields))
	for i, f := range fields {
		fieldSQL := f.ColumnName
		if f.Aggregate != "" {
			fieldSQL = fmt.Sprintf("%s(%s)", strings.ToUpper(f.Aggregate), fieldSQL)
		}
		if f.Alias != "" {
			fieldSQL = fmt.Sprintf("%s AS %s", fieldSQL, f.Alias)
		}
		parts[i] = fieldSQL
	}
	return strings.Join(parts, ", ")
}

// buildOrderBySQL 构建 ORDER BY 语句（含 ORDER BY 关键字）
func buildOrderBySQL(fields []config.OrderByField) string {
	if len(fields) == 0 {
		return ""
	}
	parts := make([]string, len(fields))
	for i, f := range fields {
		dir := f.Direction
		if dir == "" {
			dir = "ASC"
		}
		parts[i] = f.ColumnName + " " + dir
	}
	return "ORDER BY " + strings.Join(parts, ", ")
}

// firstWhereField 安全地取第一个WHERE字段（用于批量IN操作）
func firstWhereField(fields []config.SnippetField) *config.SnippetField {
	if len(fields) > 0 {
		return &fields[0]
	}
	return nil
}

// -----------------------------------------------------------------------
// 方法名自动生成
// -----------------------------------------------------------------------

func buildSelectMethodName(cfg *config.SnippetConfig) string {
	// IS NULL / IS NOT NULL 字段不计入方法名（无参数）
	effective := effectiveWhereFields(cfg.WhereFields)
	if len(effective) == 0 {
		if cfg.IsBatch {
			return "selectAll"
		}
		return "selectByFields"
	}
	parts := make([]string, len(effective))
	for i, f := range effective {
		part := capitalize(f.FieldName)
		if f.Operator == "IN" {
			part += "In"
		} else if f.Operator == "NOT IN" {
			part += "NotIn"
		}
		parts[i] = part
	}
	if cfg.IsBatch {
		return "selectBy" + strings.Join(parts, "And") + "In"
	}
	return "selectBy" + strings.Join(parts, "And")
}

func buildDeleteMethodName(cfg *config.SnippetConfig) string {
	effective := effectiveWhereFields(cfg.WhereFields)
	if len(effective) == 0 {
		return "deleteByFields"
	}
	parts := make([]string, len(effective))
	for i, f := range effective {
		part := capitalize(f.FieldName)
		if f.Operator == "IN" {
			part += "In"
		} else if f.Operator == "NOT IN" {
			part += "NotIn"
		}
		parts[i] = part
	}
	if cfg.IsBatch {
		return "deleteBy" + strings.Join(parts, "And") + "In"
	}
	return "deleteBy" + strings.Join(parts, "And")
}

func buildUpdateMethodName(cfg *config.SnippetConfig) string {
	setParts := make([]string, len(cfg.SetFields))
	for i, f := range cfg.SetFields {
		setParts[i] = capitalize(f.FieldName)
	}
	effective := effectiveWhereFields(cfg.WhereFields)
	whereParts := make([]string, len(effective))
	for i, f := range effective {
		part := capitalize(f.FieldName)
		if f.Operator == "IN" {
			part += "In"
		} else if f.Operator == "NOT IN" {
			part += "NotIn"
		}
		whereParts[i] = part
	}
	name := "update"
	if len(setParts) > 0 {
		name += strings.Join(setParts, "And")
	}
	if len(whereParts) > 0 {
		name += "By" + strings.Join(whereParts, "And")
	}
	if cfg.IsBatch {
		name += "Batch"
	}
	return name
}

// effectiveWhereFields 过滤掉 IS NULL / IS NOT NULL（无需参数）的字段
func effectiveWhereFields(fields []config.SnippetField) []config.SnippetField {
	var result []config.SnippetField
	for _, f := range fields {
		if f.Operator != "IS NULL" && f.Operator != "IS NOT NULL" {
			result = append(result, f)
		}
	}
	return result
}

// -----------------------------------------------------------------------
// 辅助函数
// -----------------------------------------------------------------------

func buildJavaParams(cfg *config.SnippetConfig) string {
	// 过滤掉 IS NULL / IS NOT NULL（无需Java参数）
	effective := effectiveWhereFields(cfg.WhereFields)
	var parts []string

	// 处理 Limit 参数
	if cfg.Operation == "select" && cfg.HasLimit && !cfg.IsLimitFixed {
		limitName := cfg.LimitValue
		if limitName == "" {
			limitName = "limit"
		}
		parts = append(parts, fmt.Sprintf("@Param(\"%s\") Integer %s", limitName, limitName))
	}

	for _, f := range effective {
		javaType := f.JavaType
		if f.Operator == "IN" || f.Operator == "NOT IN" {
			javaType = "List<" + f.JavaType + ">"
		}
		parts = append(parts, fmt.Sprintf("@Param(\"%s\") %s %s",
			f.FieldName, javaType, f.FieldName))
	}

	// 如果只有一个 Where 条件且没有 Limit，不使用 @Param
	if len(effective) == 1 && len(parts) == 1 {
		f := effective[0]
		javaType := f.JavaType
		if f.Operator == "IN" || f.Operator == "NOT IN" {
			javaType = "List<" + f.JavaType + ">"
		}
		return fmt.Sprintf("%s %s", javaType, f.FieldName)
	}

	return strings.Join(parts, ", ")
}

// collectSnippetImports 根据片段配置收集需要的 import
func collectSnippetImports(cfg *config.SnippetConfig) []string {
	importsMap := make(map[string]bool)
	effective := effectiveWhereFields(cfg.WhereFields)

	switch cfg.Operation {
	case config.OperationSelect:
		importsMap["java.util.List"] = true // 返回类型始终是 List
		if cfg.IsBatch {
			importsMap["org.apache.ibatis.annotations.Param"] = true
		} else {
			hasLimitParam := cfg.HasLimit && !cfg.IsLimitFixed
			// @Param: 有 limit 参数 或 effective 字段 > 1
			if hasLimitParam || len(effective) > 1 {
				importsMap["org.apache.ibatis.annotations.Param"] = true
			}
			for _, f := range effective {
				if f.Operator == "IN" || f.Operator == "NOT IN" {
					importsMap["java.util.List"] = true
				}
			}
		}
	case config.OperationInsert:
		if cfg.IsBatch {
			importsMap["java.util.List"] = true
			importsMap["org.apache.ibatis.annotations.Param"] = true
		}
	case config.OperationDelete:
		if cfg.IsBatch {
			importsMap["java.util.List"] = true
			importsMap["org.apache.ibatis.annotations.Param"] = true
		} else {
			if len(effective) > 1 {
				importsMap["org.apache.ibatis.annotations.Param"] = true
			}
			for _, f := range effective {
				if f.Operator == "IN" || f.Operator == "NOT IN" {
					importsMap["java.util.List"] = true
				}
			}
		}
	case config.OperationUpdate:
		if cfg.IsBatch {
			importsMap["java.util.List"] = true
			importsMap["org.apache.ibatis.annotations.Param"] = true
		} else {
			if len(effective) > 1 {
				importsMap["org.apache.ibatis.annotations.Param"] = true
			}
			for _, f := range effective {
				if f.Operator == "IN" || f.Operator == "NOT IN" {
					importsMap["java.util.List"] = true
				}
			}
		}
	}

	result := make([]string, 0, len(importsMap))
	for imp := range importsMap {
		result = append(result, imp)
	}
	sort.Strings(result)
	return result
}

func lastPart(fullType string) string {
	parts := strings.Split(fullType, ".")
	return parts[len(parts)-1]
}

func capitalize(s string) string {
	if s == "" {
		return ""
	}
	return strings.ToUpper(s[:1]) + s[1:]
}

func renderTemplate(name, tmplStr string, data interface{}) (string, error) {
	funcMap := template.FuncMap{
		"last": func(i int, arr interface{}) bool { return false }, // placeholder
	}
	tmpl, err := template.New(name).Funcs(funcMap).Parse(tmplStr)
	if err != nil {
		return "", fmt.Errorf("解析模板失败(%s): %v", name, err)
	}
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", fmt.Errorf("执行模板失败(%s): %v", name, err)
	}
	return buf.String(), nil
}

// -----------------------------------------------------------------------
// XML 模板（使用预计算的 WhereSQL / OrderBySQL 字符串，避免模板内复杂逻辑）
// -----------------------------------------------------------------------

const selectSnippetTemplate = `    <!-- 自定义查询 - {{.MethodName}} -->
    <select id="{{.MethodName}}" resultType="{{.ModelType}}">
        SELECT {{.SelectSQL}}
        FROM {{.TableName}}
{{- if .IsBatch}}
        WHERE {{.InField.ColumnName}} IN
        <foreach collection="list" item="item" open="(" separator="," close=")">
            #{item}
        </foreach>
{{- else if .WhereSQL}}
        WHERE {{.WhereSQL}}
{{- end}}
{{- if .OrderBySQL}}
        {{.OrderBySQL}}
{{- end}}
{{- if .HasLimit}}
    {{- if .IsLimitFixed}}
        LIMIT {{if .LimitValue}}{{.LimitValue}}{{else}}10{{end}}
    {{- else}}
        LIMIT #{{print "{"}}{{if .LimitValue}}{{.LimitValue}}{{else}}limit{{end}}{{print "}"}}
    {{- end}}
{{- end}}
    </select>`

const insertSnippetTemplate = `    <!-- 自定义插入 - {{.MethodName}} -->
{{- if .IsBatch}}
    <insert id="{{.MethodName}}" parameterType="java.util.List">
        INSERT INTO {{.TableName}} (
            {{range $i, $f := .InsertFields}}{{if $i}}, {{end}}{{$f.ColumnName}}{{end}}
        )
        VALUES
        <foreach collection="list" item="item" separator=",">
            ({{range $i, $f := .InsertFields}}{{if $i}}, {{end}}#{{"{"}}item.{{$f.FieldName}},jdbcType={{$f.JdbcType}}{{"}"}}{{end}})
        </foreach>
    </insert>
{{- else}}
    <insert id="{{.MethodName}}" parameterType="{{.ModelType}}">
        INSERT INTO {{.TableName}} (
            {{range $i, $f := .InsertFields}}{{if $i}}, {{end}}{{$f.ColumnName}}{{end}}
        )
        VALUES (
            {{range $i, $f := .InsertFields}}{{if $i}}, {{end}}#{{"{"}}{{$f.FieldName}},jdbcType={{$f.JdbcType}}{{"}"}}{{end}}
        )
    </insert>
{{- end}}`

const deleteSnippetTemplate = `    <!-- 自定义删除 - {{.MethodName}} -->
    <delete id="{{.MethodName}}">
        DELETE FROM {{.TableName}}
{{- if .IsBatch}}
        WHERE {{.InField.ColumnName}} IN
        <foreach collection="list" item="item" open="(" separator="," close=")">
            #{item}
        </foreach>
{{- else if .WhereSQL}}
        WHERE {{.WhereSQL}}
{{- end}}
    </delete>`

const updateSnippetTemplate = `    <!-- 自定义更新 - {{.MethodName}} -->
{{- if .IsBatch}}
    <update id="{{.MethodName}}" parameterType="java.util.List">
        <foreach collection="list" item="item" separator=";">
            UPDATE {{.TableName}}
            <set>
                {{range $i, $f := .SetFields}}{{$f.ColumnName}} = #{{"{"}}item.{{$f.FieldName}},jdbcType={{$f.JdbcType}}{{"}"}},
                {{end}}
            </set>
            WHERE {{.WhereSQL}}
        </foreach>
    </update>
{{- else}}
    <update id="{{.MethodName}}" parameterType="{{.ModelType}}">
        UPDATE {{.TableName}}
        <set>
            {{range $i, $f := .SetFields}}{{$f.ColumnName}} = #{{"{"}}{{$f.FieldName}},jdbcType={{$f.JdbcType}}{{"}"}},
            {{end}}
        </set>
        WHERE {{.WhereSQL}}
    </update>
{{- end}}`
