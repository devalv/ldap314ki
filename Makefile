setup:
	go install golang.org/x/tools/gopls@latest
	go install golang.org/x/tools/cmd/goimports@latest
	go install github.com/securego/gosec/v2/cmd/gosec@latest
	go install mvdan.cc/gofumpt@latest
	go install github.com/sqs/goreturns@latest
	go install -v github.com/go-critic/go-critic/cmd/gocritic@latest
	go install golang.org/x/tools/cmd/goimports@latest
	go install github.com/daixiang0/gci@latest

fmt:
	gofmt -w -s ./cmd
	gofumpt -w ./cmd
	goimports -w ./cmd
	gci write ./cmd
	go mod tidy
	golangci-lint run

test:
	go test ./... -race

# cover:
# 	go test ./... -race -cover

build:
	$(MAKE) fmt
	go build -o application ./cmd/app

run:
	go run ./cmd/app --config ./config-local.yml
