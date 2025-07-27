# SSL证书监控系统使用文档

## 概述

SSL证书监控系统是集成在重启人生API服务中的定时任务功能，能够自动监控SSL证书的有效期，并在证书即将到期时发送邮件提醒。

## 功能特性

- ✅ **自动监控**: 定时检查SSL证书有效期
- ✅ **智能提醒**: 在14天、7天、3天到期前分别发送邮件
- ✅ **防重复发送**: 相同阶段的提醒只发送一次
- ✅ **HTML邮件**: 美观的邮件模板，包含详细信息
- ✅ **配置灵活**: 支持自定义域名、检查间隔、提醒阈值
- ✅ **日志记录**: 完整的监控活动日志

## 配置说明

### 1. 邮件配置

在 `configs/live.yaml` 中配置SMTP邮件服务：

\`\`\`yaml
email:
  smtp:
    server: smtp.163.com          # SMTP服务器
    port: "465"                   # SMTP端口（SSL）
    user: A Second Chance         # 发件人名称
    password: your-auth-password  # 邮箱授权密码（重要：需要配置真实密码）
    from: your-email@163.com      # 发件人邮箱
  emails:
    pic_email: 985751277@qq.com   # 管理员邮箱（接收提醒）
\`\`\`

### 2. SSL监控配置

\`\`\`yaml
ssl_monitor:
  domain: asecondchance.cn        # 监控的域名
  check_interval: 24h             # 检查间隔（24h=每天，12h=每12小时，1h=每小时）
  alert_days: [14, 7, 3]         # 提醒阈值（天数）
\`\`\`

## 使用方法

### 启动监控服务

SSL证书监控功能已集成到主API服务中，启动API服务即自动开始监控：

\`\`\`bash
# 方法1: 直接运行
./build/restart-life-api

# 方法2: 使用VS Code任务
# 运行 "运行 API 服务" 任务
\`\`\`

### 查看监控状态

监控活动会记录在服务日志中：

\`\`\`bash
# 查看日志
tail -f logs/app.log

# 或者查看控制台输出
# 监控相关的日志会包含 "SSL certificate" 关键字
\`\`\`

### 测试功能

运行测试脚本验证配置：

\`\`\`bash
./scripts/test_ssl_monitor.sh
\`\`\`

## 邮件模板

系统使用 `template/ssl_cert_expiry.html` 作为邮件模板，包含：

- 🚨 **警告级别图标**: 根据剩余天数显示不同图标
- 📊 **证书详情**: 域名、颁发机构、有效期等
- 📋 **建议操作**: 续期步骤指导
- ℹ️ **系统信息**: 服务器和配置信息

## 工作流程

1. **服务启动**: API服务启动时自动启动SSL证书监控任务
2. **定时检查**: 按配置的间隔（默认24小时）检查证书
3. **状态判断**: 比较剩余天数与提醒阈值
4. **发送提醒**: 达到阈值时发送邮件（防重复发送）
5. **日志记录**: 记录所有监控活动

## 提醒阶段

| 剩余天数 | 提醒级别 | 图标 | 邮件主题前缀 |
|---------|---------|------|-------------|
| <= 3天  | 🚨 紧急 | 🚨 | 【重要】 |
| <= 7天  | ⚠️ 警告 | ⚠️ | 【重要】 |
| <= 14天 | 📋 提醒 | 📋 | 【重要】 |

## 故障排除

### 1. 邮件发送失败

**问题**: 日志显示邮件发送失败

**解决方案**:
- 检查SMTP配置是否正确
- 确认邮箱授权密码是否有效
- 验证网络连接和防火墙设置

### 2. 证书检查失败

**问题**: 无法获取证书信息

**解决方案**:
- 确认域名可以正常访问
- 检查SSL证书是否正确安装
- 验证网络连接

### 3. 任务未启动

**问题**: 监控任务没有运行

**解决方案**:
- 检查服务日志中是否有错误信息
- 确认配置文件格式正确
- 重启API服务

## 自定义配置

### 修改检查间隔

\`\`\`yaml
ssl_monitor:
  check_interval: 12h  # 每12小时检查一次
  # 或
  check_interval: 1h   # 每小时检查一次（测试用）
\`\`\`

### 修改提醒阈值

\`\`\`yaml
ssl_monitor:
  alert_days: [30, 14, 7, 3, 1]  # 增加30天和1天的提醒
\`\`\`

### 修改监控域名

\`\`\`yaml
ssl_monitor:
  domain: your-domain.com  # 监控其他域名
\`\`\`

## 日志示例

\`\`\`
[2025-07-26 20:53:01] INFO SSL certificate monitor job starting domain=asecondchance.cn interval=24h alert_days=[14 7 3]
[2025-07-26 20:53:01] INFO checking SSL certificate domain=asecondchance.cn
[2025-07-26 20:53:01] INFO certificate info retrieved domain=asecondchance.cn days_left=63 issuer="TrustAsia DV TLS RSA CA 2025"
[2025-07-26 20:53:01] INFO no alert needed, certificate status is normal or alert already sent
\`\`\`

## 安全注意事项

1. **邮箱密码**: 使用邮箱服务商提供的授权密码，不是登录密码
2. **权限控制**: 确保配置文件的读取权限正确设置
3. **网络安全**: 监控过程中的网络通信使用SSL加密

## 技术架构

- **定时任务管理器**: `internal/job/manager.go`
- **SSL证书监控任务**: `internal/job/ssl_cert_monitor.go`
- **邮件服务**: 复用现有的 `EmailService`
- **配置管理**: 集成到现有配置系统
- **日志系统**: 使用统一的日志框架

---

## 常见问题

**Q: 如何测试邮件发送功能？**
A: 可以临时修改 `alert_days` 为更大的值（如 [100, 80, 65]），这样会触发邮件发送。

**Q: 可以监控多个域名吗？**
A: 当前版本支持单域名监控，如需监控多域名可以创建多个监控任务实例。

**Q: 邮件模板可以自定义吗？**
A: 可以，修改 `template/ssl_cert_expiry.html` 文件即可自定义邮件样式和内容。

**Q: 如何停止监控？**
A: 停止API服务即可停止所有定时任务，包括SSL证书监控。
