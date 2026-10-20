package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	goapitosdk "github.com/apito-io/go-admin-sdk"
)

func main() {
	if err := run(os.Args[1:]); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}

func run(args []string) error {
	opts, err := parseGenArgs(args)
	if err != nil {
		return err
	}
	return generate(opts)
}

type genOptions struct {
	schema string
	out    string
	models string
	module string
	prune  bool
	set    bool
}

func parseGenArgs(args []string) (genOptions, error) {
	fs := flag.NewFlagSet("apito-gen", flag.ContinueOnError)
	fs.SetOutput(os.Stderr)
	opts := genOptions{}
	fs.StringVar(&opts.schema, "schema", "", "introspection JSON")
	fs.StringVar(&opts.out, "out", "", "output directory")
	fs.StringVar(&opts.models, "models", "", "comma-separated model names")
	fs.StringVar(&opts.module, "module", "", "consumer Go module path")
	fs.BoolVar(&opts.prune, "prune", false, "delete graphql files this run did not write")
	if err := fs.Parse(args); err != nil {
		return genOptions{}, err
	}
	opts.set = opts.schema != "" || opts.out != "" || opts.models != "" || opts.module != "" || opts.prune
	return opts, nil
}

func generate(opts genOptions) error {
	root, _ := os.Getwd()
	if !opts.set {
		for !fileExists(filepath.Join(root, "go.mod")) && root != "/" {
			root = filepath.Dir(root)
		}
	}
	schemaPath := opts.schema
	if schemaPath == "" {
		schemaPath = envOr("APITO_SCHEMA_FILE", filepath.Join(root, "schema", "apito_introspection.json"))
	}
	outDir := opts.out
	if outDir == "" {
		outDir = filepath.Join(root, "codegen")
	}
	opsDir := filepath.Join(outDir, "operations")

	var filter []string
	modelList := opts.models
	if modelList == "" {
		modelList = os.Getenv("APITO_MODELS")
	}
	for _, part := range strings.Split(modelList, ",") {
		if s := strings.TrimSpace(part); s != "" {
			filter = append(filter, s)
		}
	}

	raw, err := os.ReadFile(schemaPath)
	if err != nil {
		return fmt.Errorf("read schema: %w", err)
	}
	var intro map[string]interface{}
	if err := json.Unmarshal(raw, &intro); err != nil {
		return fmt.Errorf("parse schema: %w", err)
	}

	schema, err := goapitosdk.ParseIntrospection(intro, filter)
	if err != nil {
		return fmt.Errorf("parse introspection: %w", err)
	}

	if err := os.MkdirAll(opsDir, 0o755); err != nil {
		return err
	}
	if err := os.WriteFile(
		filepath.Join(outDir, "schema.graphql"),
		[]byte(goapitosdk.IntrospectionToSDL(intro)),
		0o644,
	); err != nil {
		return err
	}
	if opts.module != "" {
		note := fmt.Sprintf("package apito\n\n// Code generated for module %s. DO NOT EDIT.\n", opts.module)
		if err := os.WriteFile(filepath.Join(outDir, "doc.go"), []byte(note), 0o644); err != nil {
			return err
		}
		genqlient := "schema: schema.graphql\noperations:\n  - operations/*.graphql\ngenerated: genqlient/genqlient.go\npackage: genqlient\nbindings:\n  DateTime:\n    type: string\n  JSON:\n    type: map[string]interface{}\n"
		if err := os.WriteFile(filepath.Join(outDir, "genqlient.yaml"), []byte(genqlient), 0o644); err != nil {
			return err
		}
	}

	written := map[string]bool{}
	for _, model := range schema.Models {
		fields := make([]string, 0, len(model.Fields))
		for _, f := range model.Fields {
			fields = append(fields, f.Name)
		}
		doc := goapitosdk.NewDocumentBuilder(model.Name).GenerateGraphqlFile(fields)
		name := model.Name + ".graphql"
		path := filepath.Join(opsDir, name)
		if err := os.WriteFile(path, []byte(doc), 0o644); err != nil {
			return err
		}
		written[name] = true
		fmt.Println("wrote operations/" + name)
	}
	if opts.prune {
		entries, err := os.ReadDir(opsDir)
		if err != nil {
			return err
		}
		for _, entry := range entries {
			if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".graphql") || written[entry.Name()] {
				continue
			}
			if err := os.Remove(filepath.Join(opsDir, entry.Name())); err != nil {
				return err
			}
		}
	}
	fmt.Printf("wrote schema.graphql (%d models)\n", len(schema.Models))
	return nil
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
