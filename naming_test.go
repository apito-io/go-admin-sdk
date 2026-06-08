package goapitosdk

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

type namingVector struct {
	Input                  string `json:"input"`
	SingularResourceName   string `json:"singularResourceName"`
	MultipleResourceName   string `json:"multipleResourceName"`
	GraphQLTypeName        string `json:"graphqlTypeName"`
	GraphQLTypeNamePlural  string `json:"graphqlTypeNamePlural"`
}

func TestNamingVectorsParity(t *testing.T) {
	path := filepath.Join("test", "fixtures", "naming_vectors.json")
	raw, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read fixtures: %v", err)
	}
	var vectors []namingVector
	if err := json.Unmarshal(raw, &vectors); err != nil {
		t.Fatalf("parse fixtures: %v", err)
	}
	for _, row := range vectors {
		t.Run(row.Input, func(t *testing.T) {
			if got := ApitoSingularResourceName(row.Input); got != row.SingularResourceName {
				t.Errorf("ApitoSingularResourceName: got %q want %q", got, row.SingularResourceName)
			}
			if got := ApitoMultipleResourceName(row.Input); got != row.MultipleResourceName {
				t.Errorf("ApitoMultipleResourceName: got %q want %q", got, row.MultipleResourceName)
			}
			if got := ApitoSingularGraphQLTypeName(row.Input); got != row.GraphQLTypeName {
				t.Errorf("ApitoSingularGraphQLTypeName: got %q want %q", got, row.GraphQLTypeName)
			}
			if got := ApitoListGraphQLTypeName(row.Input); got != row.GraphQLTypeNamePlural {
				t.Errorf("ApitoListGraphQLTypeName: got %q want %q", got, row.GraphQLTypeNamePlural)
			}
		})
	}
}

func TestDocumentBuilderLoan(t *testing.T) {
	doc := NewDocumentBuilder("loan").BuildListQuery([]string{"loan_id"})
	if !containsAll(doc, "loanList(", "loanListCount(", "LOANLIST_INPUT_WHERE_PAYLOAD") {
		t.Errorf("unexpected list query: %s", doc)
	}
}

func containsAll(s string, parts ...string) bool {
	for _, p := range parts {
		if !stringsContains(s, p) {
			return false
		}
	}
	return true
}

func stringsContains(s, sub string) bool {
	return len(s) >= len(sub) && (s == sub || len(sub) == 0 || indexOf(s, sub) >= 0)
}

func indexOf(s, sub string) int {
	for i := 0; i+len(sub) <= len(s); i++ {
		if s[i:i+len(sub)] == sub {
			return i
		}
	}
	return -1
}
