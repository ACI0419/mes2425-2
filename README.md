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
