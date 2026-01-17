package generator

import (
	"strings"
)

// TemplateFuncs 模板函数
var TemplateFuncs = map[string]interface{}{
	"title":   strings.Title,
	"toLower": strings.ToLower,
	"toUpper": strings.ToUpper,
}
