# MinIO Docker 安装指南

## 1. Docker 安装 MinIO

```bash
docker run -d \
  --name minio \
  --restart always \
  -p 9000:9000 \
  -p 9090:9090 \
  -e MINIO_ROOT_USER=minioadmin \
  -e MINIO_ROOT_PASSWORD=minioadmin \
  -v /data/minio:/data \
  minio/minio server /data --console-address ":9090"
```

## 2. 参数说明

- **9000** �?S3 API 端口（应用上�?下载用）
- **9090** �?Web 管理控制台端口（浏览器访问）
- **MINIO_ROOT_USER** �?管理员用户名（对�?third.yaml 中的 AccessKeyID�?- **MINIO_ROOT_PASSWORD** �?管理员密码（对应 third.yaml 中的 SecretAccessKey�?- **/data/minio** �?宿主机数据持久化目录

## 3. 验证是否启动成功

```bash
# 检查容器状�?docker ps | grep minio

# 健康检�?curl http://127.0.0.1:9000/minio/health/live
```

## 4. 访问管理控制�?
浏览器打开：http://服务器IP:9090

用户名：minioadmin
密码：minioadmin

## 5. third.yaml 配置

确保 `server/etc/third.yaml` 中的配置与上面一致：

```yaml
ObjectStorage:
  Type: minio
  Endpoint: http://127.0.0.1:9000
  Bucket: fdim
  AccessKeyID: minioadmin
  SecretAccessKey: minioadmin
```

Bucket `fdim` 会在首次使用时自动创建�?
## 6. 安装后重�?third 服务

```bash
# 重新编译
go build -o bin/third rpc/third/third.go

# 重启 third RPC 服务
./bin/third -f etc/third.yaml
```
