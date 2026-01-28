package utils

import (
	"strings"
	"unicode"
)

// DBStringToCamelCase 将数据库下划线命名转换为Java驼峰命名
// 例如: user_name -> userName, USER_ID -> userId, ExposureTime -> exposureTime
func DBStringToCamelCase(s string) string {
	if s == "" {
		return ""
	}

	parts := strings.Split(s, "_")
	if len(parts) == 1 {
		// 没有下划线,只做首字母小写处理,保持可能的驼峰格式
		return FirstLower(s)
	}

	// 有下划线时,先转小写处理大写列名
	s = strings.ToLower(s)
	parts = strings.Split(s, "_")

	// 第一个部分保持小写,其余部分首字母大写
	result := parts[0]
	for i := 1; i < len(parts); i++ {
		if parts[i] != "" {
			result += FirstUpper(parts[i])
		}
	}

	return result
}

// DBStringToPascalCase 将数据库下划线命名转换为Java Pascal命名(首字母大写)
// 例如: user_name -> UserName, USER_ID -> UserId, ExposureTime -> ExposureTime
func DBStringToPascalCase(s string) string {
	if s == "" {
		return ""
	}

	parts := strings.Split(s, "_")
	if len(parts) == 1 {
		// 没有下划线,只做首字母大写处理,保持可能的驼峰格式
		return FirstUpper(s)
	}

	// 有下划线时,先转小写处理大写列名
	s = strings.ToLower(s)
	parts = strings.Split(s, "_")

	var result string
	for _, part := range parts {
		if part != "" {
			result += FirstUpper(part)
		}
	}

	return result
}

// FirstUpper 首字母大写
func FirstUpper(s string) string {
	if s == "" {
		return ""
	}
	r := []rune(s)
	r[0] = unicode.ToUpper(r[0])
	return string(r)
}

// FirstLower 首字母小写
func FirstLower(s string) string {
	if s == "" {
		return ""
	}
	r := []rune(s)
	r[0] = unicode.ToLower(r[0])
	return string(r)
}

// CamelCaseToDBString 将驼峰命名转换为下划线命名
// 例如: userName -> user_name, userId -> user_id
func CamelCaseToDBString(s string) string {
	if s == "" {
		return ""
	}

	var result []rune
	for i, r := range s {
		if unicode.IsUpper(r) {
			if i > 0 {
				result = append(result, '_')
			}
			result = append(result, unicode.ToLower(r))
		} else {
			result = append(result, r)
		}
	}

	return string(result)
}
