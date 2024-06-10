# 使用官方的Go镜像作为基础镜像
FROM golang:1.20-alpine

# 设置工作目录
WORKDIR /app

# 复制当前目录下的所有文件到工作目录
COPY . .

# 安装必要的包
RUN apk add --no-cache git

# 下载依赖
RUN go mod tidy

# 构建Go程序
RUN go build -o main .

# 设置环境变量
ENV ONEAPI_OLD_SQL_DSN=""
ENV ONEAPI_NEW_SQL_DSN=""


# 运行Go程序
CMD ["./main"]