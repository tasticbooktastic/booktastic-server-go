build:
    GOOS=linux GOARCH=amd64 go mod download github.com/Azure/go-ansiterm
    GOOS=linux GOARCH=amd64 go build -o functions/online ./functions-src/online.go