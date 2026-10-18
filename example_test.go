package goapitosdk

import "fmt"

func ExampleNewClient() {
	client := NewClient(Config{
		BaseURL: "http://localhost:5055/system/graphql",
		APIKey:  "ak_example",
	})
	fmt.Println(client != nil)
	// Output:
	// true
}

func ExampleClient_SearchResources() {
	client := NewClient(Config{
		BaseURL: "http://localhost:5055/system/graphql",
		APIKey:  "ak_example",
	})
	_ = client
	model := "post"
	filter := map[string]interface{}{"limit": 10}
	fmt.Println(model, filter["limit"])
	// Output:
	// post 10
}
