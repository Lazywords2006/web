.PHONY: all build run test clean server help

# 默认目标
all: build

# 编译主程序
build:
	@echo "Building secure-client..."
	go build -o secure-client main.go
	@echo "Build complete: ./secure-client"

# 编译测试服务器
build-server:
	@echo "Building test server..."
	go build -o test-server test-server.go
	@echo "Build complete: ./test-server"

# 编译所有
build-all: build build-server

# 运行主程序
run:
	@echo "Running secure-client..."
	./secure-client

# 运行测试服务器
server:
	@echo "Starting mock license server..."
	go run test-server.go

# 运行测试
test:
	@echo "Running tests..."
	go test ./...

# 整理依赖
tidy:
	@echo "Tidying dependencies..."
	go mod tidy

# 清理编译产物
clean:
	@echo "Cleaning build artifacts..."
	rm -f secure-client test-server
	@echo "Clean complete"

# 跨平台编译
build-windows:
	@echo "Building for Windows..."
	GOOS=windows GOARCH=amd64 go build -o secure-client.exe main.go

build-linux:
	@echo "Building for Linux..."
	GOOS=linux GOARCH=amd64 go build -o secure-client-linux main.go

build-mac:
	@echo "Building for macOS..."
	GOOS=darwin GOARCH=amd64 go build -o secure-client-mac main.go

# 编译所有平台
build-cross: build-windows build-linux build-mac
	@echo "Cross-platform build complete"

# 优化编译（减小体积）
build-release:
	@echo "Building release version..."
	go build -ldflags="-s -w" -o secure-client main.go
	@echo "Release build complete: ./secure-client"

# 创建配置文件
config:
	@if [ ! -f config.json ]; then \
		cp config.json.example config.json; \
		echo "Created config.json from example"; \
	else \
		echo "config.json already exists"; \
	fi

# 帮助信息
help:
	@echo "Available targets:"
	@echo "  make build          - Build main client"
	@echo "  make build-server   - Build test server"
	@echo "  make build-all      - Build both client and server"
	@echo "  make run            - Run the client"
	@echo "  make server         - Run mock license server"
	@echo "  make test           - Run tests"
	@echo "  make tidy           - Tidy dependencies"
	@echo "  make clean          - Remove build artifacts"
	@echo "  make build-cross    - Build for all platforms"
	@echo "  make build-release  - Build optimized release"
	@echo "  make config         - Create config.json from example"
	@echo "  make help           - Show this help message"
