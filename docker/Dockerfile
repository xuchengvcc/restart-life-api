# 多阶段构建
# 第一阶段：构建应用
FROM golang:1.23.8-alpine AS builder

# 设置工作目录
WORKDIR /app

# 设置Go模块代理和环境变量
ENV GOPROXY=https://proxy.golang.org,direct
ENV GO111MODULE=on
ENV CGO_ENABLED=0
ENV GOOS=linux

# 复制go mod文件
COPY go.mod go.sum ./

# 下载依赖
RUN go mod download

# 复制源代码
COPY . .

# 构建应用
RUN go build -a -installsuffix cgo -o main cmd/server/main.go

# 第二阶段：运行时镜像
FROM alpine:latest

# 使用国内Alpine镜像源
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories

# 安装必要的包
RUN apk --no-cache add ca-certificates tzdata

# 设置时区
ENV TZ=Asia/Shanghai

# 创建非root用户
RUN addgroup -g 1001 appgroup && \
    adduser -D -s /bin/sh -u 1001 -G appgroup appuser

# 设置工作目录
WORKDIR /app

# 从构建阶段复制二进制文件
COPY --from=builder /app/main .

# 复制配置文件
COPY --from=builder /app/configs ./configs

# 创建日志目录
RUN mkdir -p logs && chown -R appuser:appgroup /app

# 切换到非root用户
USER appuser

# 暴露端口
EXPOSE 8080

# 设置健康检查
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
  CMD wget --no-verbose --tries=1 --spider http://localhost:8080/health || exit 1

# 启动应用
CMD ["./main"]
