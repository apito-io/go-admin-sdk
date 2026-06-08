.PHONY: gen-operations gen-types gen test

gen-operations:
	go run ./cmd/apito-gen

gen-types:
	go run github.com/Khan/genqlient

gen: gen-operations gen-types

test:
	go test ./...
