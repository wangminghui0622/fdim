# 验证码配置说明

## 📧 邮箱注册默认验证码

根据配置文件 `chat/config/chat-rpc-chat.yml`，邮箱注册的**默认验证码是：`666666`**

### 配置说明

在 `chat-rpc-chat.yml` 文件中：

```yaml
verifyCode:
  superCode: "666666"  # 超级验证码（测试/开发环境使用）
  
  mail:
    use: "superCode"  # 使用超级验证码模式
```

### 验证码模式

1. **`superCode` 模式**（当前配置）：
   - 使用固定的超级验证码：`666666`
   - 适用于开发/测试环境
   - 不需要真实发送邮件

2. **`mail` 模式**（生产环境）：
   - 需要配置真实的邮件服务器
   - 会发送真实的验证码到邮箱
   - 需要配置以下参数：
     ```yaml
     mail:
       use: "mail"
       title: "验证码邮件标题"
       senderMail: "发送者邮箱"
       senderAuthorizationCode: "邮箱授权码"
       smtpAddr: "smtp服务器地址"
       smtpPort: 587
     ```

### 使用方式

**邮箱注册时，直接输入验证码：`666666`**

### 修改验证码

如果需要修改默认验证码，编辑配置文件：

```bash
# 编辑配置文件
vim /Users/xx/workspace/src/chat/config/chat-rpc-chat.yml

# 修改 superCode 的值
verifyCode:
  superCode: "123456"  # 改为你想要的验证码

# 重启服务使配置生效
docker compose restart openim-chat
```

### 手机号注册

手机号注册也使用相同的超级验证码 `666666`，配置如下：

```yaml
phone:
  use: "superCode"  # 使用超级验证码
```

### 验证码规则

- **长度**：6 位数字（`len: 6`）
- **有效期**：300 秒（5 分钟）
- **有效次数**：5 次
- **每日最大发送次数**：10 次

### 注意事项

⚠️ **生产环境请勿使用 `superCode` 模式！**

生产环境应该：
1. 配置真实的邮件服务器（SMTP）
2. 将 `mail.use` 改为 `"mail"`
3. 配置真实的短信服务（如阿里云短信）
4. 将 `phone.use` 改为 `"ali"` 或其他短信服务商
