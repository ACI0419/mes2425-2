# mes2425-2

## 构建和运行

**构建镜像**:

```bash
docker build -t mes-minimal:latest .
```

**运行容器**:

```bash
docker run -d \
    --name mes-minimal \
    -p 80:80 \
    -v mes-data:/var/lib/mysql \
    mes-minimal:latest
```

## 容器化启动

直接运行 docker-compose.yml 文件即可

## 可执行文件启动

### 数据库

启动 ./backend/docker-compose.yml 中的 mysql 服务

### 后端

启动 ./release/be.exe 即可

### 前端

不会打包，从代码启动吧

## 从代码启动

### 后端

```bash
cd backend
go run main.go
```

### 前端

```bash
cd frontend
npm install
npm start
```