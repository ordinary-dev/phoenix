all: fmt vet

fmt:
	go fmt *.go
	go fmt backend/*.go
	go fmt views/*.go

vet:
	go vet *.go
	go vet backend/*.go
	go vet views/*.go
