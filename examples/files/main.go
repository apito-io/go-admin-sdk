// Files REST example. Run from repo root: go run ./examples/files/
package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	goapitosdk "github.com/apito-io/go-admin-sdk"
)

func main() {
	baseURL := os.Getenv("APITO_BASE_URL")
	if baseURL == "" {
		baseURL = "http://localhost:5050/system/graphql"
	}
	apiKey := os.Getenv("APITO_API_KEY")
	if apiKey == "" {
		log.Fatal("APITO_API_KEY is required")
	}

	client := goapitosdk.NewClient(goapitosdk.Config{
		BaseURL: baseURL,
		APIKey:  apiKey,
		Timeout: 30 * time.Second,
	})
	ctx := context.Background()

	list, err := client.ListFiles(ctx, "", 20, 0)
	if err != nil {
		log.Fatalf("ListFiles: %v", err)
	}
	fmt.Printf("Files (total=%d):\n", list.Total)
	for _, f := range list.Files {
		fmt.Printf("  - %s (%s) %s\n", f.FileName, f.ID, f.URL)
	}

	path := os.Getenv("APITO_UPLOAD_FILE")
	if path == "" {
		return
	}
	data, err := os.ReadFile(path)
	if err != nil {
		log.Fatalf("read file: %v", err)
	}
	uploaded, err := client.UploadFile(ctx, goapitosdk.UploadFileParams{
		FileName: path,
		Content:  data,
	})
	if err != nil {
		log.Fatalf("UploadFile: %v", err)
	}
	fmt.Printf("Uploaded: %s -> %s\n", uploaded.ID, uploaded.URL)
}
