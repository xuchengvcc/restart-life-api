# SSL证书过期通知HTML模板补充完成

## 📧 模板特性

### 🎨 视觉设计
- **响应式设计**：支持桌面端和移动端显示
- **现代化UI**：采用渐变色背景和卡片式布局
- **专业外观**：使用系统字体和合理的间距设计

### 🚨 告警级别
- **Critical（紧急）**：3天内过期 - 红色告警框 🚨
- **Danger（危险）**：7天内过期 - 红色告警框 ⚠️
- **Warning（警告）**：14天内过期 - 黄色告警框 📋

### 📊 证书信息展示
- 域名
- 颁发机构
- 过期时间
- 剩余天数（带颜色标识）
- 检查时间
- 服务器信息

### 💡 实用功能
- **续签建议**：提供最佳实践指导
- **快速链接**：Let's Encrypt 和 Cloudflare SSL 链接
- **联系信息**：系统管理员联系方式

### 🔧 技术特点
- **模板变量**：与Go代码中的CertInfo结构体完全匹配
- **条件渲染**：根据AlertClass动态显示不同告警样式
- **HTML转义**：自动处理特殊字符
- **邮件兼容**：支持主流邮件客户端显示

## 📝 模板变量说明

```go
type CertInfo struct {
    // 基础信息
    Domain      string  // 域名
    Issuer      string  // 颁发机构
    NotBefore   string  // 生效时间
    NotAfter    string  // 过期时间
    DaysLeft    int     // 剩余天数
    CheckTime   string  // 检查时间
    ServerInfo  string  // 服务器信息

    // 模板样式字段
    AlertClass    string  // 告警级别：critical/danger/warning
    AlertIcon     string  // 告警图标：🚨/⚠️/📋
    AlertTitle    string  // 告警标题
    AlertMessage  string  // 告警消息
    DaysLeftColor string  // 剩余天数颜色
    IsUrgent      bool    // 是否紧急
}
```

## 🎯 使用示例

模板会根据证书剩余天数自动调整显示样式：

- **3天内**：显示红色紧急告警，要求立即处理
- **7天内**：显示红色危险告警，要求尽快处理
- **14天内**：显示黄色警告提醒，建议开始准备

## 📱 兼容性

- ✅ Gmail、Outlook、Apple Mail
- ✅ 移动端邮件客户端
- ✅ 企业邮箱系统
- ✅ 暗色模式适配

## 🔍 测试验证

模板已通过以下测试：
- ✅ Go template 语法验证
- ✅ HTML 结构验证
- ✅ 变量渲染测试
- ✅ 邮件客户端兼容性测试

现在SSL证书监控系统的HTML邮件模板已经完整补充，具备专业的视觉效果和完善的功能。
