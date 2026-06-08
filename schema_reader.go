package goapitosdk

import (
	"fmt"
	"strings"
)

// ApitoSchemaField is a parsed field from introspection.
type ApitoSchemaField struct {
	Name string
}

// ApitoSchemaModel is a parsed model from introspection.
type ApitoSchemaModel struct {
	Name   string
	Fields []ApitoSchemaField
}

// ApitoSchema is the full parsed schema.
type ApitoSchema struct {
	Models []ApitoSchemaModel
}

// ParseIntrospection extracts list models from a GraphQL introspection JSON map.
func ParseIntrospection(intro map[string]interface{}, modelFilter []string) (*ApitoSchema, error) {
	data, _ := intro["data"].(map[string]interface{})
	schema, _ := data["__schema"].(map[string]interface{})
	if schema == nil {
		return nil, fmt.Errorf("invalid introspection JSON")
	}
	types, _ := schema["types"].([]interface{})
	queryType, _ := schema["queryType"].(map[string]interface{})
	queryFields, _ := queryType["fields"].([]interface{})

	var models []ApitoSchemaModel
	for _, qf := range queryFields {
		field, _ := qf.(map[string]interface{})
		name, _ := field["name"].(string)
		if !strings.HasSuffix(name, "List") || strings.HasSuffix(name, "ListCount") {
			continue
		}
		modelName := listFieldToSnakeModel(name)
		if len(modelFilter) > 0 && !contains(modelFilter, modelName) {
			continue
		}
		createPayload := findCreatePayload(types, modelName)
		fields := defaultFields()
		if createPayload != "" {
			fields = fieldsFromInput(types, createPayload)
		}
		models = append(models, ApitoSchemaModel{Name: modelName, Fields: fields})
	}
	return &ApitoSchema{Models: models}, nil
}

func listFieldToSnakeModel(listField string) string {
	camel := strings.TrimSuffix(listField, "List")
	var b strings.Builder
	for i, r := range camel {
		if i > 0 && r >= 'A' && r <= 'Z' {
			b.WriteByte('_')
		}
		b.WriteRune(r)
	}
	return strings.ToLower(b.String())
}

func findCreatePayload(types []interface{}, modelName string) string {
	expected := ApitoGraphQLComposedTypeName(modelName, "Create_Payload")
	for _, t := range types {
		m, _ := t.(map[string]interface{})
		if m["name"] == expected {
			return expected
		}
	}
	return ""
}

func fieldsFromInput(types []interface{}, typeName string) []ApitoSchemaField {
	for _, t := range types {
		m, _ := t.(map[string]interface{})
		if m["name"] != typeName {
			continue
		}
		inputFields, _ := m["inputFields"].([]interface{})
		var out []ApitoSchemaField
		for _, inf := range inputFields {
			f, _ := inf.(map[string]interface{})
			n, _ := f["name"].(string)
			if strings.HasPrefix(n, "_") {
				continue
			}
			out = append(out, ApitoSchemaField{Name: n})
		}
		if len(out) > 0 {
			return out
		}
	}
	return defaultFields()
}

func defaultFields() []ApitoSchemaField {
	return []ApitoSchemaField{{Name: "id"}}
}

func contains(slice []string, v string) bool {
	for _, s := range slice {
		if s == v {
			return true
		}
	}
	return false
}

// IntrospectionToSDL converts introspection JSON to minimal GraphQL SDL.
func IntrospectionToSDL(intro map[string]interface{}) string {
	data, _ := intro["data"].(map[string]interface{})
	schema, _ := data["__schema"].(map[string]interface{})
	types, _ := schema["types"].([]interface{})
	var b strings.Builder
	b.WriteString("# AUTO-GENERATED — DO NOT EDIT\n\n")
	for _, t := range types {
		m, _ := t.(map[string]interface{})
		kind, _ := m["kind"].(string)
		name, _ := m["name"].(string)
		if name == "" || strings.HasPrefix(name, "__") {
			continue
		}
		switch kind {
		case "OBJECT":
			b.WriteString("type " + name + " {\n")
			if fields, ok := m["fields"].([]interface{}); ok {
				for _, f := range fields {
					fm, _ := f.(map[string]interface{})
					fn, _ := fm["name"].(string)
					b.WriteString(fmt.Sprintf("  %s: %s\n", fn, unwrapTypeName(fm["type"])))
				}
			}
			b.WriteString("}\n\n")
		case "INPUT_OBJECT":
			b.WriteString("input " + name + " {\n")
			if inputFields, ok := m["inputFields"].([]interface{}); ok {
				for _, f := range inputFields {
					fm, _ := f.(map[string]interface{})
					fn, _ := fm["name"].(string)
					b.WriteString(fmt.Sprintf("  %s: %s\n", fn, unwrapTypeName(fm["type"])))
				}
			}
			b.WriteString("}\n\n")
		case "ENUM":
			b.WriteString("enum " + name + " {\n")
			if enumValues, ok := m["enumValues"].([]interface{}); ok {
				for _, v := range enumValues {
					vm, _ := v.(map[string]interface{})
					b.WriteString("  " + vm["name"].(string) + "\n")
				}
			}
			b.WriteString("}\n\n")
		case "SCALAR":
			b.WriteString("scalar " + name + "\n\n")
		}
	}
	return b.String()
}

func unwrapTypeName(t interface{}) string {
	m, _ := t.(map[string]interface{})
	kind, _ := m["kind"].(string)
	switch kind {
	case "NON_NULL":
		return unwrapTypeName(m["ofType"]) + "!"
	case "LIST":
		return "[" + unwrapTypeName(m["ofType"]) + "]"
	default:
		n, _ := m["name"].(string)
		if n == "" {
			return "String"
		}
		return n
	}
}
