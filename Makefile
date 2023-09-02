tidy:
	@go mod tidy

build: tidy
	@go build -ldflags "-s -w" -o rsearch.exe main.go

build-linux: tidy
	@CGO_ENABLED=0 GOARCH=amd64 GOOS=linux go build -ldflags "-s -w" -o rsearch main.go
