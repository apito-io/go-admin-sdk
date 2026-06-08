package goapitosdk

import (
	"fmt"
	"strings"
)

// DocumentBuilder builds GraphQL operation strings (aligned with flutter_admin_sdk).
type DocumentBuilder struct {
	model string
}

// NewDocumentBuilder creates a builder for the given model id.
func NewDocumentBuilder(model string) *DocumentBuilder {
	return &DocumentBuilder{model: model}
}

func (d *DocumentBuilder) listField() string     { return ApitoMultipleResourceName(d.model) }
func (d *DocumentBuilder) countField() string    { return d.listField() + "Count" }
func (d *DocumentBuilder) singularField() string { return ApitoSingularResourceName(d.model) }
func (d *DocumentBuilder) listPascal() string    { return ApitoListGraphQLTypeName(d.model) }
func (d *DocumentBuilder) singularPascal() string {
	return ApitoSingularGraphQLTypeName(d.model)
}

// BuildListQuery returns the list + count query document.
func (d *DocumentBuilder) BuildListQuery(fields []string) string {
	vars := fmt.Sprintf(`$connection: %s
    $where: %s
    $whereCount: %s
    $sort: %s
    $page: Int
    $limit: Int`,
		ApitoConnectionFilterConditionType(d.model),
		ApitoWhereInputType(d.model),
		ApitoListCountWhereInputType(d.model),
		ApitoSortInputType(d.model),
	)
	fieldLines := strings.Join(fields, "\n      ")
	return fmt.Sprintf(`query Get%s(
    %s
) {
  %s(connection: $connection, where: $where, sort: $sort, page: $page, limit: $limit) {
    id
    data {
      %s
    }
    
    meta {
      created_at
      status
      updated_at
    }
  }    %s(connection: $connection, where: $whereCount, page: $page, limit: $limit) {
      total
    }
}`, d.listPascal(), vars, d.listField(), fieldLines, d.countField())
}

// BuildGetQuery returns single-record query document.
func (d *DocumentBuilder) BuildGetQuery(fields []string) string {
	fieldLines := strings.Join(fields, "\n      ")
	return fmt.Sprintf(`query Get%s($id: String!) {
  %s(_id: $id) {
    id
    data {
      %s
    }
    
    meta {
      created_at
      status
      updated_at
    }
  }
}`, d.singularPascal(), d.singularField(), fieldLines)
}

// BuildCreateMutation returns create mutation document.
func (d *DocumentBuilder) BuildCreateMutation(fields []string) string {
	payload := ApitoGraphQLComposedTypeName(d.model, "Create_Payload")
	connect := ApitoGraphQLComposedTypeName(d.model, "Relation_Connect_Payload")
	fieldLines := strings.Join(fields, "\n      ")
	return fmt.Sprintf(`mutation Create%s($payload: %s!, $connect: %s) {
  create%s(payload: $payload, connect: $connect, status: published) {
    id
    data {
      %s
    }
    meta {
      created_at
      status
      updated_at
    }
  }
}`, d.singularPascal(), payload, connect, d.singularPascal(), fieldLines)
}

// BuildUpdateMutation returns update mutation document.
func (d *DocumentBuilder) BuildUpdateMutation(fields []string) string {
	payload := ApitoGraphQLComposedTypeName(d.model, "Update_Payload")
	connect := ApitoGraphQLComposedTypeName(d.model, "Relation_Connect_Payload")
	disconnect := ApitoGraphQLComposedTypeName(d.model, "Relation_Disconnect_Payload")
	fieldLines := strings.Join(fields, "\n      ")
	return fmt.Sprintf(`mutation Update%s(
    $id: String!,
    $deltaUpdate: Boolean,
    $payload: %s!,
    $connect: %s,
    $disconnect: %s
) {
  update%s(_id: $id, deltaUpdate: $deltaUpdate, payload: $payload, connect: $connect, disconnect: $disconnect, status: published) {
    id
    data {
      %s
    }
    meta {
      created_at
      status
      updated_at
    }
  }
}`, d.singularPascal(), payload, connect, disconnect, d.singularPascal(), fieldLines)
}

// BuildDeleteMutation returns delete mutation document.
func (d *DocumentBuilder) BuildDeleteMutation() string {
	return fmt.Sprintf(`mutation Delete%s($ids: [String]!) {
  delete%s(_ids: $ids) {
    response
  }
}`, d.singularPascal(), d.singularPascal())
}

// GenerateGraphqlFile returns the full .graphql file for a model.
func (d *DocumentBuilder) GenerateGraphqlFile(fieldNames []string) string {
	fields := make([]string, 0, len(fieldNames))
	for _, f := range fieldNames {
		if f != "id" {
			fields = append(fields, f)
		}
	}
	if len(fields) == 0 {
		fields = []string{"id"}
	}
	var b strings.Builder
	b.WriteString("# AUTO-GENERATED — DO NOT EDIT\n")
	b.WriteString("# Model: " + d.model + "\n\n")
	b.WriteString(d.BuildListQuery(fields))
	b.WriteString("\n\n")
	b.WriteString(d.BuildGetQuery(fields))
	b.WriteString("\n\n")
	b.WriteString(d.BuildCreateMutation(fields))
	b.WriteString("\n\n")
	b.WriteString(d.BuildUpdateMutation(fields))
	b.WriteString("\n\n")
	b.WriteString(d.BuildDeleteMutation())
	b.WriteString("\n")
	return b.String()
}
