# 多阶段构建Dockerfile
FROM golang:1.21-alpine AS builder

# 安装编译依赖
RUN apk add --no-cache git gcc musl-dev sqlite-dev

# 设置工作目录
WORKDIR /build

# 复制go.mod和go.sum
COPY server/go.mod server/go.sum* ./
RUN go mod download || true

# 复制服务器代码
COPY server/ ./

# 编译服务器
RUN CGO_ENABLED=1 GOOS=linux go build -a -installsuffix cgo -ldflags="-w -s" -o license-server main.go

# 运行阶段
FROM alpine:latest

# 安装运行时依赖
RUN apk --no-cache add ca-certificates sqlite-libs

# 创建应用目录
WORKDIR /app

# 从构建阶段复制二进制文件
COPY --from=builder /build/license-server .

# 复制前端文件
COPY frontend/ ./frontend/

# 创建数据目录
RUN mkdir -p /data

# 设置环境变量
ENV DB_PATH=/data/licenses.db
ENV PORT=8080

# 暴露端口
EXPOSE 8080

# 健康检查
HEALTHCHECK --interval=30s --timeout=3s \
  CMD wget --no-verbose --tries=1 --spider http://localhost:8080/ || exit 1

# 运行服务器
CMD ["./license-server"]
