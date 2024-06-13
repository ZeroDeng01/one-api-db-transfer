# 使用官方的Go镜像作为构建阶段的基础镜像
FROM golang:1.20-alpine AS builder

# 设置工作目录
WORKDIR /app

# 复制go.mod和go.sum文件并下载依赖
COPY go.mod go.sum ./
RUN go mod download

# 复制项目的源代码
COPY . .

# 构建Go程序
RUN go build -o /main .

# 使用Alpine作为最终镜像
FROM alpine:latest

# 设置环境变量
ENV ONEAPI_OLD_SQL_DSN=""
ENV ONEAPI_NEW_SQL_DSN=""

# 从构建阶段复制二进制文件
COPY --from=builder /main /main

# 创建一个非root用户
RUN addgroup -S appgroup && adduser -S appuser -G appgroup

# 使用非root用户运行程序
USER appuser

# 运行Go程序
CMD ["/main"]