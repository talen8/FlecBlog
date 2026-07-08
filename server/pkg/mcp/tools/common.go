package tools

import (
	"reflect"

	"github.com/google/jsonschema-go/jsonschema"
)

// ============ 分页函数 ============

// NormalizePage 归一化分页参数
func NormalizePage(page, pageSize int) (int, int) {
	if page < 1 {
		page = 1
	}
	// pageSize <= 0 表示返回全部结果
	return page, pageSize
}

// ExplicitListFields keeps list-action contracts explicit when items are empty and total is zero.
// Pointer fields remain nil for non-list actions, so unrelated outputs do not gain list metadata.
func ExplicitListFields[T any](items []T, total int64) (*[]T, *int64) {
	if items == nil {
		items = make([]T, 0)
	}
	return &items, &total
}

// ============ Schema 构建函数 ============

// BuildActionSchema 构建 Action Schema
func BuildActionSchema(action, description string, payload *jsonschema.Schema) *jsonschema.Schema {
	return &jsonschema.Schema{
		Type: "object",
		Properties: map[string]*jsonschema.Schema{
			"action": {
				Type:        "string",
				Enum:        []any{action},
				Description: description,
			},
			"payload": payload,
		},
		Required: []string{"action", "payload"},
	}
}

// BuildPayloadSchema 构建 Payload Schema
func BuildPayloadSchema(properties map[string]*jsonschema.Schema, required ...string) *jsonschema.Schema {
	return &jsonschema.Schema{
		Type:       "object",
		Properties: properties,
		Required:   required,
	}
}

// FalseSchema 返回禁止额外属性的 Schema
func FalseSchema() *jsonschema.Schema {
	return &jsonschema.Schema{Not: &jsonschema.Schema{}}
}

// PageSizeSchema 返回 page_size 字段的 schema
func PageSizeSchema() *jsonschema.Schema {
	return &jsonschema.Schema{
		Type:        "integer",
		Description: "每页数量，<= 0 表示返回全部结果",
	}
}

// ============ 转换函数 ============

// ToTimeStringPtr 将带 String() 方法的类型转换为 *string
func ToTimeStringPtr(t interface{ String() string }) *string {
	if t == nil {
		return nil
	}
	v := reflect.ValueOf(t)
	if v.Kind() == reflect.Pointer && v.IsNil() {
		return nil
	}
	s := t.String()
	if s == "" {
		return nil
	}
	return &s
}

// ToStringPtr 将 string 转换为 *string（空字符串返回 nil）
func ToStringPtr(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}
