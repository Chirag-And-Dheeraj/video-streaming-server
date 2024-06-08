# Ensure we're in a Unix-like environment
ifeq (, $(shell which uname))
$(error This Makefile requires a Unix-like environment)
endif

# Install dependencies
install:
	@echo "Checking for PostgreSQL Installation..."
	@if command -v psql > /dev/null 2>&1; then \
		echo "PostgreSQL is already installed."; \
	else \
		echo "Installing PostgreSQL..."; \
		(sudo apt-get update && sudo apt-get -y install postgresql postgresql-contrib) || exit 1; \
		echo "PostgreSQL successfully installed."; \
	fi
	@echo "Checking for Go installation..."
	@if command -v go > /dev/null 2>&1; then \
		echo "Go is already installed."; \
	else \
		echo "Installing Go version 1.22.2..."; \
		wget https://golang.org/dl/go1.22.2.linux-amd64.tar.gz -O go1.22.2.tar.gz || exit 1; \
		sudo tar -C /usr/local -xzf go1.22.2.tar.gz || exit 1; \
		echo "Go version 1.22.2 installed."; \
		rm go1.22.2.tar.gz; \
		echo "Updating PATH to include Go binary directory..."; \
		echo "export PATH=\$$PATH:/usr/local/go/bin" >> ~/.bashrc; \
		. ~/.bashrc; \
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
	@echo "Starting PostgreSQL service...";
	@if command -v service > /dev/null 2>&1; then \
		sudo service postgresql start || exit 1; \
	elif command -v systemctl > /dev/null 2>&1; then \
		sudo systemctl start postgresql || exit 1; \
	elif command -v rc-service > /dev/null 2>&1; then \
		sudo rc-service postgresql start || exit 1; \
	fi
	
	@if ps aux | grep -v grep | grep postgres > /dev/null 2>&1; then \
		echo "PostgreSQL service started.";\
	else \
		echo "PostgreSQL service failed to start.";\
	fi

	@echo "Starting the Go application..."
	go run main.go

# Perform a clean start
cleanstart:
	@echo "Performing clean start..."
	make clean || exit 1
	make init || exit 1
	make start || exit 1
