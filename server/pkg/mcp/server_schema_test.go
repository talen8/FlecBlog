package mcp

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http/httptest"
	"testing"

	sdkmcp "github.com/modelcontextprotocol/go-sdk/mcp"
)

var standardJSONSchemaTypes = map[string]struct{}{
	"array":   {},
	"boolean": {},
	"integer": {},
	"null":    {},
	"number":  {},
	"object":  {},
	"string":  {},
}

func TestPublicToolSchemasUseStandardJSONSchemaTypes(t *testing.T) {
	handler := NewPublicHandler(nil, nil, nil, nil, nil, nil, nil, nil, nil)
	server := httptest.NewServer(handler)
	defer server.Close()

	ctx := context.Background()
	client := sdkmcp.NewClient(&sdkmcp.Implementation{Name: "flecblog-schema-test", Version: "1.0"}, nil)
	session, err := client.Connect(ctx, &sdkmcp.StreamableClientTransport{Endpoint: server.URL}, nil)
	if err != nil {
		t.Fatalf("connect MCP client: %v", err)
	}
	defer session.Close()

	result, err := session.ListTools(ctx, nil)
	if err != nil {
		t.Fatalf("list tools: %v", err)
	}
	for _, tool := range result.Tools {
		checkStandardJSONSchemaTypes(t, tool.Name+".input", tool.InputSchema)
		checkStandardJSONSchemaTypes(t, tool.Name+".output", tool.OutputSchema)
	}
}

func checkStandardJSONSchemaTypes(t *testing.T, name string, schema any) {
	t.Helper()
	if schema == nil {
		return
	}
	data, err := json.Marshal(schema)
	if err != nil {
		t.Fatalf("%s marshal schema: %v", name, err)
	}
	var value any
	if err := json.Unmarshal(data, &value); err != nil {
		t.Fatalf("%s unmarshal schema: %v", name, err)
	}
	walkJSONSchemaTypes(t, name, value)
}

func walkJSONSchemaTypes(t *testing.T, path string, value any) {
	t.Helper()
	walkJSONSchemaNode(t, path, value, true)
}

func walkJSONSchemaNode(t *testing.T, path string, value any, isSchema bool) {
	t.Helper()
	switch node := value.(type) {
	case map[string]any:
		if isSchema {
			if rawType, ok := node["type"]; ok {
				switch typed := rawType.(type) {
				case string:
					assertStandardJSONSchemaType(t, path+".type", typed)
				case []any:
					for i, item := range typed {
						typeName, ok := item.(string)
						if !ok {
							t.Fatalf("%s.type[%d] = %T, want string", path, i, item)
						}
						assertStandardJSONSchemaType(t, fmt.Sprintf("%s.type[%d]", path, i), typeName)
					}
				default:
					t.Fatalf("%s.type = %T, want string or array", path, rawType)
				}
			}
		}

		for key, child := range node {
			if isSchema && key == "type" {
				continue
			}
			childPath := path + "." + key
			switch key {
			case "properties", "patternProperties", "$defs", "definitions", "dependentSchemas":
				if schemaMap, ok := child.(map[string]any); ok {
					for name, schema := range schemaMap {
						walkJSONSchemaNode(t, childPath+"."+name, schema, true)
					}
					continue
				}
			case "items", "additionalProperties", "contains", "not", "if", "then", "else", "propertyNames", "unevaluatedProperties", "unevaluatedItems":
				walkJSONSchemaNode(t, childPath, child, true)
				continue
			case "allOf", "anyOf", "oneOf", "prefixItems":
				if schemas, ok := child.([]any); ok {
					for i, schema := range schemas {
						walkJSONSchemaNode(t, fmt.Sprintf("%s[%d]", childPath, i), schema, true)
					}
					continue
				}
			}
			walkJSONSchemaNode(t, childPath, child, false)
		}
	case []any:
		for i, child := range node {
			walkJSONSchemaNode(t, fmt.Sprintf("%s[%d]", path, i), child, false)
		}
	}
}

func assertStandardJSONSchemaType(t *testing.T, path, typeName string) {
	t.Helper()
	if _, ok := standardJSONSchemaTypes[typeName]; !ok {
		t.Fatalf("%s = %q, not a standard JSON Schema type", path, typeName)
	}
}
