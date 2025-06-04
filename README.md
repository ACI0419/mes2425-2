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

### 生成&校验密码

```bash
cd release

# 生成密码
> ./gc.exe -g admin123
$2a$10$jSV2ZGn6yDfMWaZTMJEy2uts08jO.Rlu2LoIFKr1ZemsB2Of.dsny

# 校验密码 切记！如果是命令行运行，需要在 $ 美元符前加 ` 反引号
>./gc.exe -c admin123 `$2a`$10`$jSV2ZGn6yDfMWaZTMJEy2uts08jO.Rlu2LoIFKr1ZemsB2Of.dsny
true
```

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