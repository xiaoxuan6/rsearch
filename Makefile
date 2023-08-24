build: tidy
	@go build -ldflags "-s -w" -o rsearch.exe

tidy:
	@go mod tidy