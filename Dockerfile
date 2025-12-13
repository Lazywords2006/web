FROM golang:1.21-alpine AS builder

# 设置工作目录
WORKDIR /app

# 安装编译依赖
RUN apk add --no-cache git gcc musl-dev

# 复制 go.mod 和 go.sum（如果存在）
COPY server/go.mod server/go.sum* ./server/
WORKDIR /app/server
RUN go mod download

# 复制服务器代码
COPY server/ ./

# 编译服务器
RUN CGO_ENABLED=1 GOOS=linux go build -ldflags="-s -w" -o license-server

# 运行阶段
FROM alpine:latest

# 安装运行时依赖
RUN apk --no-cache add ca-certificates sqlite-libs

# 创建非 root 用户
RUN addgroup -g 1000 appuser && \
    adduser -D -u 1000 -G appuser appuser

# 设置工作目录
WORKDIR /app

# 从构建阶段复制二进制文件
COPY --from=builder /app/server/license-server .

# 创建数据目录
RUN mkdir -p /app/data && \
    chown -R appuser:appuser /app

# 切换到非 root 用户
USER appuser

# 暴露端口
EXPOSE 8080

# 设置环境变量
ENV PORT=8080
ENV DB_PATH=/app/data/licenses.db

# 启动服务
CMD ["./license-server"]
