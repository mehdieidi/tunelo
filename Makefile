tidy:
	go mod tidy

fmt:
	go fmt ./...
	golines -w -t 4 -m 100 --ignore-generated .

lint:
	golangci-lint run --config "./config/.golangci.yaml"
