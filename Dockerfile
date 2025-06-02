# 多阶段构建 - 前端构建阶段
FROM node:18-alpine AS frontend-builder
WORKDIR /app/frontend
COPY frontend/package*.json ./
RUN npm ci --only=production
COPY frontend/ ./
RUN npm run build

# 多阶段构建 - 后端构建阶段
FROM golang:1.23-alpine AS backend-builder
WORKDIR /app/backend
COPY backend/go.mod backend/go.sum ./
RUN go mod download
COPY backend/ ./
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

# 最终运行阶段
FROM ubuntu:22.04

# 安装必要的软件包
RUN apt-get update && apt-get install -y \
    nginx \
    mysql-server \
    supervisor \
    curl \
    && rm -rf /var/lib/apt/lists/*

# 创建工作目录
WORKDIR /app

# 复制构建好的前端文件
COPY --from=frontend-builder /app/frontend/build /var/www/html

# 复制构建好的后端可执行文件
COPY --from=backend-builder /app/backend/main /app/
COPY --from=backend-builder /app/backend/configs /app/configs
COPY --from=backend-builder /app/backend/docs /app/docs

# 创建Nginx配置（只暴露前端）
RUN rm /etc/nginx/sites-enabled/default
COPY docker/nginx-minimal.conf /etc/nginx/sites-enabled/default

# 创建Supervisor配置
COPY docker/supervisord-minimal.conf /etc/supervisor/conf.d/supervisord.conf

# 创建启动脚本
COPY docker/start-minimal.sh /start.sh
RUN chmod +x /start.sh

# 创建必要的目录和设置权限
RUN mkdir -p /var/log/supervisor \
    && mkdir -p /var/run/mysqld \
    && chown mysql:mysql /var/run/mysqld

# 只暴露前端端口
EXPOSE 80

# 启动脚本
CMD ["/start.sh"]