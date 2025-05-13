# 使用官方的 Go 镜像作为基础镜像
FROM golang:1.24-alpine

# 设置工作目录
WORKDIR /app

# 复制 go.mod 和 go.sum 文件
COPY go.mod go.sum ./

# 下载依赖
RUN go mod download

# 复制项目文件
COPY . .

# 构建应用
RUN go build -o go-reptile .

# 暴露端口（如果需要）
EXPOSE 8080

# 运行应用
CMD ["./go-reptile"]