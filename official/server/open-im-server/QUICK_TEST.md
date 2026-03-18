# 快速测试指南

## ✅ 系统状态检查

所有服务现在应该都正常运行了！

### 1. 检查服务状态

```bash
docker compose ps
```

应该看到所有服务都是 `Up (healthy)` 状态。

### 2. 测试后端 API

```bash
# 测试 OpenIM Server API
curl http://localhost:10002/

# 测试 Chat Service API  
curl http://localhost:10008/
```

### 3. 访问前端

1. **打开浏览器**，访问：http://localhost:11001
2. **点击注册按钮**
3. **打开开发者工具**（F12 或 Cmd+Option+I）
4. **查看 Network 标签**，点击注册后应该看到：
   - 有 API 请求发出（到 `localhost:10008` 或类似地址）
   - 检查请求状态码和响应

### 4. 如果注册按钮还是没反应

#### 检查浏览器 Console

1. 打开开发者工具（F12）
2. 切换到 **Console** 标签
3. 点击注册按钮
4. 查看是否有红色错误信息

常见错误：
- `CORS error`：跨域问题
- `Network error`：网络连接问题
- `404 Not Found`：API 路径错误
- `500 Internal Server Error`：服务器错误

#### 检查 Network 请求

1. 打开开发者工具（F12）
2. 切换到 **Network** 标签
3. 点击注册按钮
4. 查看是否有请求发出：
   - 如果有请求：查看状态码和响应内容
   - 如果没有请求：可能是前端 JavaScript 错误

### 5. 常见问题解决

#### 问题：前端页面空白或加载失败

```bash
# 检查前端容器
docker logs openim-web-front

# 重启前端
docker compose restart openim-web-front
```

#### 问题：API 请求失败

```bash
# 检查后端服务
docker logs openim-chat
docker logs openim-server

# 检查服务是否在监听
docker exec openim-chat netstat -tlnp | grep 10008
docker exec openim-server netstat -tlnp | grep 10002
```

#### 问题：CORS 错误

前端和后端需要在同一个 Docker 网络中，确保 `docker-compose.yml` 中所有服务都在 `openim` 网络下。

### 6. 完整重启（如果还有问题）

```bash
cd /Users/xx/workspace/src/open-im-server

# 停止所有服务
docker compose down

# 重新启动
docker compose up -d

# 等待服务启动（约 1-2 分钟）
sleep 60

# 检查状态
docker compose ps
```

## 🎯 现在应该可以正常使用了！

访问 http://localhost:11001，点击注册按钮应该能正常响应了。

如果还有问题，请提供：
1. 浏览器 Console 的错误信息
2. Network 标签中的请求详情
3. `docker compose logs openim-chat | tail -30` 的输出
