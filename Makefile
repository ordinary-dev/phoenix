all: fmt vet

fmt:
	gofmt -s -w .

vet:
	go vet ./...
