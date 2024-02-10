APP_NAME=hashnode

build:
	@echo "Building..."
	@go build -o bin/$(APP_NAME) -ldflags "-s -w" main.go
	@echo "Build complete"

