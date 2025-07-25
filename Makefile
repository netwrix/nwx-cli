# Makefile for nwx CLI

# Build variables
BINARY_NAME=nwx
BINARY_PATH=./$(BINARY_NAME)
INSTALL_PATH=/usr/local/bin/$(BINARY_NAME)

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod

# Build the binary
build:
	$(GOBUILD) -o $(BINARY_PATH) -v .

# Clean build artifacts
clean:
	$(GOCLEAN)
	rm -f $(BINARY_PATH)

# Test the application
test:
	$(GOTEST) -v ./...

# Install dependencies
deps:
	$(GOMOD) tidy
	$(GOMOD) download

# Install globally (requires sudo)
install: build
	sudo cp $(BINARY_PATH) $(INSTALL_PATH)
	sudo chmod +x $(INSTALL_PATH)
	@echo "✅ nwx installed globally to $(INSTALL_PATH)"
	@echo "You can now run 'nwx' from anywhere!"

# Install to user bin directory (no sudo required)
install-user: build
	mkdir -p ~/bin
	cp $(BINARY_PATH) ~/bin/$(BINARY_NAME)
	chmod +x ~/bin/$(BINARY_NAME)
	@echo "✅ nwx installed to ~/bin/$(BINARY_NAME)"
	@echo "Make sure ~/bin is in your PATH:"
	@echo "  echo 'export PATH=\"\$$HOME/bin:\$$PATH\"' >> ~/.zshrc"
	@echo "  source ~/.zshrc"

# Uninstall from system
uninstall:
	sudo rm -f $(INSTALL_PATH)
	@echo "✅ nwx removed from system"

# Uninstall from user bin
uninstall-user:
	rm -f ~/bin/$(BINARY_NAME)
	@echo "✅ nwx removed from ~/bin"

# Run the application
run:
	$(GOBUILD) -o $(BINARY_PATH) -v .
	$(BINARY_PATH)

# Help
help:
	@echo "Available commands:"
	@echo "  make build       - Build the binary"
	@echo "  make install     - Install globally (requires sudo)"
	@echo "  make install-user - Install to ~/bin (no sudo)"
	@echo "  make uninstall   - Remove from system"
	@echo "  make uninstall-user - Remove from ~/bin"
	@echo "  make test        - Run tests"
	@echo "  make clean       - Clean build artifacts"
	@echo "  make deps        - Install dependencies"
	@echo "  make run         - Build and run"

.PHONY: build clean test deps install install-user uninstall uninstall-user run help