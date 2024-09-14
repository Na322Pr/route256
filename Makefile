APP_NAME=pvz-cli-app
BUILD_DIR=build

GO=go

build: clean
	$(GO) build -o $(BUILD_DIR)/$(APP_NAME) cmd/main.go

run: build
	./$(BUILD_DIR)/$(APP_NAME)

install:
	@echo "Installing dependencies..."
	$(GO) mod tidy
	$(GO) install github.com/uudashr/gocognit/cmd/gocognit@latest
	$(GO) install github.com/fzipp/gocyclo/cmd/gocyclo@latest

update:
	@echo "Updating dependencies..."
	$(GO) get -u ./...

gocognit:
	@echo "Checking gocognit..."
	gocognit -over 5 .

gocyclo:
	@echo "Checking gocyclo..."
	gocyclo -over 5 .

clean: 
	rm -rf $(BUILD_DIR)/$(APP_NAME)