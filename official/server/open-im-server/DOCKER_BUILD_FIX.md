# Docker Compose 构建问题解决方案

## 问题 1: 警告信息已解决 ✅

已在 `.env` 文件中添加了空值，消除以下警告：
- `KAFKA_USERNAME` / `KAFKA_PASSWORD`
- `ETCD_USERNAME` / `ETCD_PASSWORD`

这些变量设置为空字符串表示不使用认证（默认模式）。

## 问题 2: 关于 "找官方镜像" 的警告 ⚠️

**这是正常行为！** Docker Compose 在构建前会先检查镜像是否存在：
- `openim-server:local` - 如果不存在，会从本地代码构建
- `openim-chat:local` - 如果不存在，会从本地代码构建

这些警告可以忽略，Docker Compose 会自动切换到构建模式。

## 问题 3: golang:1.22-alpine 拉取失败 🔧

### 解决方案 A: 配置 Docker 镜像加速器（推荐）

在 Docker Desktop 中配置镜像加速器：

1. 打开 Docker Desktop
2. 进入 Settings → Docker Engine
3. 添加以下配置：

```json
{
  "registry-mirrors": [
    "https://docker.mirrors.ustc.edu.cn",
    "https://hub-mirror.c.163.com",
    "https://mirror.baidubce.com"
  ]
}
```

4. 点击 "Apply & Restart"

### 解决方案 B: 手动拉取镜像

在构建前先手动拉取 golang 镜像：

```bash
docker pull golang:1.22-alpine
```

### 解决方案 C: 使用国内镜像源（临时方案）

如果上述方案都不行，可以修改 Dockerfile 使用国内镜像：

在 `open-im-server/build/images/monolith/openim-server/Dockerfile` 和 
`chat/build/images/monolith/openim-chat/Dockerfile` 中：

将：
```dockerfile
FROM golang:1.22-alpine AS builder
```

改为：
```dockerfile
FROM registry.cn-hangzhou.aliyuncs.com/google_containers/golang:1.22-alpine AS builder
```

（注意：需要确认该镜像是否存在）

## 重新构建

配置完成后，重新运行：

```bash
docker-compose up -d --build
```

`--build` 参数会强制重新构建本地镜像。
