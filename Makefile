tidy:
	@go mod tidy

build: tidy
	go build -ldflags "-s -w" -o rsearch.exe main.go
