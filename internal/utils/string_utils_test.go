package utils

import (
	"testing"
)

func TestDBStringToCamelCase(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"user_name", "userName"},
		{"user_id", "userId"},
		{"created_at", "createdAt"},
		{"USER_NAME", "userName"},
		{"username", "username"},
		{"ExposureTime", "exposureTime"}, // 驼峰格式保持
		{"FocalLength", "focalLength"},   // 驼峰格式保持
		{"", ""},
	}

	for _, tt := range tests {
		result := DBStringToCamelCase(tt.input)
		if result != tt.expected {
			t.Errorf("DBStringToCamelCase(%q) = %q, want %q", tt.input, result, tt.expected)
		}
	}
}

func TestDBStringToPascalCase(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"user_name", "UserName"},
		{"user_id", "UserId"},
		{"created_at", "CreatedAt"},
		{"USER_NAME", "UserName"},
		{"username", "Username"},
		{"ExposureTime", "ExposureTime"}, // 驼峰格式保持
		{"FocalLength", "FocalLength"},   // 驼峰格式保持
		{"", ""},
	}

	for _, tt := range tests {
		result := DBStringToPascalCase(tt.input)
		if result != tt.expected {
			t.Errorf("DBStringToPascalCase(%q) = %q, want %q", tt.input, result, tt.expected)
		}
	}
}

func TestCamelCaseToDBString(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"userName", "user_name"},
		{"userId", "user_id"},
		{"createdAt", "created_at"},
		{"UserName", "user_name"},
		{"username", "username"},
		{"", ""},
	}

	for _, tt := range tests {
		result := CamelCaseToDBString(tt.input)
		if result != tt.expected {
			t.Errorf("CamelCaseToDBString(%q) = %q, want %q", tt.input, result, tt.expected)
		}
	}
}
