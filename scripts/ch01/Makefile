
BINARY_NAME=type_conversion

.DEFAULT_GOAL := run
## 產生二進制檔
build: vet
	GOARCH=amd64 GOOS=linux go build -o ${BINARY_NAME}-linux main.go
## 產生二進制檔，但在臨時目錄
run:
	go run ../../pkg/ch01/${BINARY_NAME}/main.go

clean: build
	go clean
	rm ${BINARY_NAME}-linux

fmt: 
	go fmt ../../pkg/ch01/${BINARY_NAME}/main.go

lint: fmt
	golint ../../pkg/ch01/${BINARY_NAME}/main.go

vet: fmt
	go vet ../../pkg/ch01/${BINARY_NAME}/main.go
