FROM golang:1.23.8-alpine3.20 AS builder
# 设置国内代理加速依赖下载
ENV GOPROXY=https://goproxy.cn,direct

WORKDIR /app

COPY  ./go.mod ./go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags '-extldflags "-static"' -o ./cmd/payment/main ./cmd/payment

FROM alpine:3.20.0
WORKDIR /app

# 安装必要依赖（兼容库 + 时区）
RUN apk --no-cache add ca-certificates libc6-compat tzdata && \
    cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime && \
    echo "Asia/Shanghai" > /etc/timezone

# 复制配置文件，保持相对路径结构
COPY --from=builder /app/configs /app/configs
COPY --from=builder /app/.env /app/.env
# 从构建阶段复制二进制文件
COPY --from=builder /app/cmd/payment/main /app/src/payment/main

# 使用非 root 用户增强安全性
RUN addgroup -S appgroup && adduser -S appuser -G appgroup
USER appuser

EXPOSE 8080 9090
CMD ["/app/src/payment/main"]