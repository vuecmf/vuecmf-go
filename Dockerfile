# 使用官方的Go镜像作为构建环境
FROM golang:1.19-alpine AS builder

# 设置 GOPROXY 环境变量
ENV CGO_ENABLED 0
ENV GOPROXY=https://goproxy.cn,direct
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories
RUN apk update --no-cache && apk add --no-cache tzdata

# 设置工作目录
WORKDIR /app

# 将当前目录的内容复制到容器的工作目录中
COPY . .

# 下载所有依赖
RUN go mod download

# 构建项目
RUN go build -o main .

# 使用官方的Alpine镜像作为运行环境
FROM alpine:latest

# 设置工作目录
WORKDIR /app

# 安装时区数据
RUN apk add --no-cache tzdata && \
    cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime && \
    echo "Asia/Shanghai" > /etc/timezone && \
    apk del tzdata

# 从构建阶段复制二进制文件
COPY --from=builder /app/main .

# 确保二进制文件具有可执行权限
RUN chmod +x /app/main

# 暴露端口
EXPOSE 8080

# 运行应用
CMD ["./main"]
