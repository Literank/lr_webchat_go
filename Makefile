# Binary name
BINARY_NAME=lr_webchat

build:
	@echo "Building $(BINARY_NAME)..."
	go build -o $(BINARY_NAME)

build-linux:
	@echo "Building $(BINARY_NAME) for Linux..."
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o $(BINARY_NAME)