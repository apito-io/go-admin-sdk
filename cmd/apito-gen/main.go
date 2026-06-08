package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	goapitosdk "github.com/apito-io/go-admin-sdk"
)

func main() {
	root, _ := os.Getwd()
	for !fileExists(filepath.Join(root, "go.mod")) && root != "/" {
		root = filepath.Dir(root)
	}
	schemaPath := envOr("APITO_SCHEMA_FILE", filepath.Join(root, "schema", "apito_introspection.json"))
	outDir := filepath.Join(root, "codegen")
	opsDir := filepath.Join(outDir, "operations")

	var filter []string
	if m := os.Getenv("APITO_MODELS"); m != "" {
		for _, part := range strings.Split(m, ",") {
			if s := strings.TrimSpace(part); s != "" {
				filter = append(filter, s)
			}
		}
	}

	raw, err := os.ReadFile(schemaPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "read schema: %v\n", err)
		os.Exit(1)
	}
	var intro map[string]interface{}
	if err := json.Unmarshal(raw, &intro); err != nil {
		fmt.Fprintf(os.Stderr, "parse schema: %v\n", err)
		os.Exit(1)
	}

	schema, err := goapitosdk.ParseIntrospection(intro, filter)
	if err != nil {
		fmt.Fprintf(os.Stderr, "parse introspection: %v\n", err)
		os.Exit(1)
	}

	if err := os.MkdirAll(opsDir, 0o755); err != nil {
		panic(err)
	}
	if err := os.WriteFile(
		filepath.Join(outDir, "schema.graphql"),
		[]byte(goapitosdk.IntrospectionToSDL(intro)),
		0o644,
	); err != nil {
		panic(err)
	}

	for _, model := range schema.Models {
		fields := make([]string, 0, len(model.Fields))
		for _, f := range model.Fields {
			fields = append(fields, f.Name)
		}
		doc := goapitosdk.NewDocumentBuilder(model.Name).GenerateGraphqlFile(fields)
		path := filepath.Join(opsDir, model.Name+".graphql")
		if err := os.WriteFile(path, []byte(doc), 0o644); err != nil {
			panic(err)
		}
		fmt.Println("wrote operations/" + model.Name + ".graphql")
	}
	fmt.Printf("wrote schema.graphql (%d models)\n", len(schema.Models))
}

func envOr(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}
