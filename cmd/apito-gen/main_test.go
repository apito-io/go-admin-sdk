package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestApitoGenWritesConsumerModuleAndPrunes(t *testing.T) {
	dir := t.TempDir()
	schema := []byte(`{"data":{"__schema":{"queryType":{"name":"Query","fields":[{"name":"postList"}]},"types":[]}}}`)
	schemaPath := filepath.Join(dir, "schema.json")
	if err := os.WriteFile(schemaPath, schema, 0o644); err != nil {
		t.Fatal(err)
	}
	out := filepath.Join(dir, "internal", "apito")
	ops := filepath.Join(out, "operations")
	if err := os.MkdirAll(ops, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(ops, "stale.graphql"), []byte("old"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := generate(genOptions{
		schema: schemaPath,
		out:    out,
		module: "example.com/app",
		prune:  true,
		set:    true,
	}); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(filepath.Join(ops, "post.graphql")); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(filepath.Join(ops, "stale.graphql")); !os.IsNotExist(err) {
		t.Fatalf("stale file still present: %v", err)
	}
	doc, err := os.ReadFile(filepath.Join(out, "doc.go"))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(doc), "example.com/app") {
		t.Fatalf("doc.go = %s", doc)
	}
	cfg, err := os.ReadFile(filepath.Join(out, "genqlient.yaml"))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(cfg), "genqlient/genqlient.go") {
		t.Fatalf("genqlient.yaml = %s", cfg)
	}
	if strings.Contains(string(doc), "apt_") || strings.Contains(string(doc), "ak_") {
		t.Fatal("generated Go embedded a credential")
	}
}
