# 故障排查指南

## 问题：点击注册按钮没反应

### 原因分析

1. **前端容器未启动**：`openim-web-front` 容器可能没有运行
2. **后端服务未启动**：`openim-server` 的 `start-config.yml` 配置为空，导致没有服务启动
3. **API 连接失败**：前端无法连接到后端 API

### 解决步骤

#### 1. 检查所有服务状态

```bash
cd /Users/xx/workspace/src/open-im-server
docker compose ps
```

确保以下服务都是 `Up` 状态：
- `openim-web-front`
- `openim-admin-front`
- `openim-server`
- `openim-chat`
- `mongo`, `redis`, `etcd`, `kafka`, `minio`

#### 2. 如果前端容器未启动

```bash
docker compose up -d openim-web-front openim-admin-front
```

#### 3. 如果后端服务未启动（已修复）

已修复 Dockerfile，现在会使用本地的 `start-config.yml` 文件。

重新构建并启动：

```bash
docker compose up -d --build openim-server
```

#### 4. 验证服务

```bash
# 检查 openim-server 是否在监听端口
docker exec openim-server cat /openim-server/start-config.yml

# 应该看到 serviceBinaries 不为空，包含：
# - openim-api
# - openim-msggateway
# - openim-rpc-* 等

# 检查端口是否监听
docker exec openim-server ss -tlnp | grep -E "10001|10002"

# 测试 API
curl http://localhost:10002/healthz
curl http://localhost:10008/api/user/register
```

#### 5. 查看日志

```bash
# 查看前端日志
docker compose logs -f openim-web-front

# 查看后端日志
docker compose logs -f openim-server
docker compose logs -f openim-chat
```

### 浏览器端检查

1. **打开开发者工具**（F12 或 Cmd+Option+I）
2. **查看 Console**：检查是否有 JavaScript 错误
3. **查看 Network**：点击注册按钮，查看是否有 API 请求发出
   - 应该看到请求到 `http://localhost:10008/api/user/register` 或类似地址
   - 检查请求状态码（200 表示成功，4xx/5xx 表示错误）

### 常见错误

#### 错误：`connection refused`
- **原因**：后端服务未启动或端口未监听
- **解决**：检查 `openim-server` 和 `openim-chat` 是否正常运行

#### 错误：`CORS error`
- **原因**：跨域请求被阻止
- **解决**：确保前端和后端在同一个 Docker 网络中，或配置 CORS

#### 错误：`404 Not Found`
- **原因**：API 路径不正确
- **解决**：检查前端配置的 API 地址是否正确

### 完整重启流程

如果问题依然存在，尝试完整重启：

```bash
cd /Users/xx/workspace/src/open-im-server

# 停止所有服务
docker compose down

# 重新构建并启动
docker compose up -d --build

# 等待服务启动（约 1-2 分钟）
sleep 60

# 检查状态
docker compose ps

# 查看日志
docker compose logs -f
```

### 联系支持

如果问题仍未解决，请提供：
1. `docker compose ps` 的输出
2. `docker compose logs openim-server` 的最后 50 行
3. `docker compose logs openim-chat` 的最后 50 行
4. 浏览器 Console 和 Network 标签的截图
