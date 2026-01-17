package generator

import (
	"strings"
	"text/template"
)

// TemplateFuncs 自定义模板函数
var TemplateFuncs = template.FuncMap{
	"title": strings.Title,
}
