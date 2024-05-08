# Ensure we're in a Unix-like environment
ifeq (, $(shell which uname))
$(error This Makefile requires a Unix-like environment)
endif

# Install dependencies
install:
	@echo "Checking for Go installation..."
	@if command -v go > /dev/null 2>&1; then \
		echo "Go is already installed."; \
	else \
		echo "Installing Go version 1.22.2..."; \
		wget https://golang.org/dl/go1.22.2.linux-amd64.tar.gz -O go1.22.2.tar.gz || exit 1; \
		tar -C /usr/local -xzf go1.22.2.tar.gz || exit 1; \
		echo "Go version 1.22.2 installed."; \
		echo "Updating PATH to include Go binary directory..."; \
		echo "export PATH=\$$PATH:/usr/local/go/bin" >> ~/.bashrc; \
		source ~/.bashrc; \
	fi
	@echo "Checking for FFMpeg installation..."
	@if command -v ffmpeg > /dev/null 2>&1; then \
		echo "FFMpeg is already installed."; \
	else \
		echo "Installing FFMpeg..."; \
		sudo apt-get update && sudo apt-get install -y ffmpeg || exit 1; \
		echo "FFMpeg installed."; \
	fi
	@echo "Installing Go dependencies..."
	go mod download || exit 1
	@echo "Dependencies installed."

# Clean up directories and files
clean:
	@echo "Cleaning up directories and files..."
	@if [ -d "video" ]; then rm -r video; fi
	@if [ -d "segments" ]; then rm -r segments; fi
	@if [ -f "database.db" ]; then rm database.db; fi
	@echo "Clean up complete."

# Initialize directories
init:
	@echo "Initializing directories..."
	@mkdir -p video
	@mkdir -p segments
	@echo "Initialization complete."

# Start the Go application
start:
	@echo "Starting the Go application..."
	go run main.go

# Perform a clean start
cleanstart:
	@echo "Performing clean start..."
	make clean || exit 1
	make init || exit 1
	make start || exit 1
