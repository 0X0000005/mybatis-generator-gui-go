package generator

import (
	"strings"
	"text/template"
)

func init() {
	// 注册自定义模板函数
	template.FuncMap{
		"title": strings.Title,
	}
}
