# 编译阶段
FROM golang:1.25-alpine AS builder
WORKDIR /app

RUN apk add --no-cache git

# 利用 Docker 缓存层优化依赖下载
COPY go.mod go.sum ./
RUN go mod download

COPY . .
# -ldflags="-s -w": 去除调试信息和符号表，减小二进制体积
ARG APP_VERSION=dev
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w -X 'flec_blog/internal/service.AppVersion=${APP_VERSION}'" -o flec_blog ./cmd/main.go

# 运行阶段
FROM alpine:latest
RUN apk add --no-cache ca-certificates tzdata su-exec

ENV TZ=Asia/Shanghai

RUN addgroup -g 1000 appuser && \
    adduser -D -u 1000 -G appuser appuser

WORKDIR /app

COPY --from=builder /app/flec_blog .
# 只读资源打包在镜像内
COPY --from=builder /app/templates ./templates

# 统一挂载点：所有需要持久化的数据都在 /app/data/uploads 下
RUN mkdir -p /app/data/uploads /app/data/logs && \
    chown -R appuser:appuser /app

EXPOSE 8080

# 启动脚本：以 root 确保目录存在并设置权限，然后降权运行应用
ENTRYPOINT ["sh", "-c", "mkdir -p /app/data/uploads /app/data/logs && chown -R appuser:appuser /app/data && ls -la /app/data && exec su-exec appuser \"$@\"", "--"]
CMD ["./flec_blog"]
