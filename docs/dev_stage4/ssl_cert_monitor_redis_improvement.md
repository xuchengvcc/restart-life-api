# SSL证书监控防重复发送机制改进

## 问题分析

原始实现使用内存中的 `map[string]bool` 来记录已发送的告警，存在以下问题：

1. **服务重启丢失状态**: 内存中的状态在服务重启后会丢失，导致重复发送告警
2. **无过期机制**: 状态永久保存，可能导致证书续期后仍然无法发送新的告警
3. **不支持集群**: 多实例部署时，每个实例都有自己的状态，可能导致重复发送

## 解决方案

使用Redis存储告警状态，并设置合理的过期时间：

### 1. Redis键设计

```go
func (j *SSLCertMonitorJob) getAlertKey(domain string, alertDay int) string {
    return fmt.Sprintf("ssl_cert_alert:%s:%d_days", domain, alertDay)
}
```

**键命名规则**: `ssl_cert_alert:{domain}:{alertDay}_days`

**示例**:
- `ssl_cert_alert:asecondchance.cn:14_days`
- `ssl_cert_alert:asecondchance.cn:7_days`
- `ssl_cert_alert:asecondchance.cn:3_days`

### 2. 过期时间策略

```go
// 计算过期时间：比告警阈值多一天，确保不会重复发送
expiration := time.Duration(alertDay+1) * 24 * time.Hour

// 如果剩余天数很少，设置较短的过期时间
if daysLeft <= 1 {
    expiration = 2 * 24 * time.Hour // 2天
} else if daysLeft <= 3 {
    expiration = 4 * 24 * time.Hour // 4天
}
```

**过期时间逻辑**:
- 14天告警 → 15天过期
- 7天告警 → 8天过期
- 3天告警 → 4天过期
- 1天及以下 → 2天过期

### 3. 核心方法实现

#### 检查告警状态
```go
func (j *SSLCertMonitorJob) isAlertSent(domain string, alertDay int) bool {
    ctx := context.Background()
    key := j.getAlertKey(domain, alertDay)

    result, err := j.redisClient.Exists(ctx, key).Result()
    if err != nil {
        j.logger.WithError(err).Error("failed to check alert status in redis")
        return false // 出错时认为未发送，允许重新尝试
    }

    return result > 0
}
```

#### 设置告警状态
```go
func (j *SSLCertMonitorJob) setAlertSent(domain string, alertDay int, daysLeft int) error {
    ctx := context.Background()
    key := j.getAlertKey(domain, alertDay)

    // 计算过期时间...
    value := fmt.Sprintf("sent_at:%d", time.Now().Unix())
    err := j.redisClient.Set(ctx, key, value, expiration).Err()
    // ...
}
```

#### 清理过期状态
```go
func (j *SSLCertMonitorJob) cleanupExpiredAlerts(domain string, currentDaysLeft int) {
    ctx := context.Background()

    // 清理比当前剩余天数更大的告警状态
    for _, alertDay := range j.alertDays {
        if currentDaysLeft > alertDay {
            key := j.getAlertKey(domain, alertDay)
            j.redisClient.Del(ctx, key)
        }
    }
}
```

## 改进后的执行流程

### 1. 检查证书状态
```go
certInfo, err := j.getCertInfo(j.domain)
```

### 2. 判断是否需要告警
```go
for _, day := range j.alertDays {
    if certInfo.DaysLeft <= day {
        if !j.isAlertSent(j.domain, day) {  // 从Redis检查
            shouldAlert = true
            alertDay = day
            break
        }
    }
}
```

### 3. 发送告警并记录状态
```go
if shouldAlert {
    if err := j.sendAlertEmail(certInfo); err != nil {
        return err
    }

    // 记录到Redis，设置过期时间
    j.setAlertSent(j.domain, alertDay, certInfo.DaysLeft)
}
```

### 4. 清理过期状态
```go
j.cleanupExpiredAlerts(j.domain, certInfo.DaysLeft)
```

## 优势对比

| 特性 | 原实现（内存） | 新实现（Redis） |
|------|---------------|----------------|
| 持久性 | ❌ 服务重启丢失 | ✅ 持久化存储 |
| 过期机制 | ❌ 无自动清理 | ✅ 自动过期 |
| 集群支持 | ❌ 各实例独立 | ✅ 共享状态 |
| 故障恢复 | ❌ 状态丢失 | ✅ 状态保留 |
| 内存消耗 | 较少 | 几乎无 |
| 网络依赖 | 无 | 需要Redis |

## 实际场景示例

### 场景1：正常告警流程
1. 证书剩余15天 → 发送14天告警 → Redis记录15天过期
2. 证书剩余8天 → 发送7天告警 → Redis记录8天过期
3. 证书剩余4天 → 发送3天告警 → Redis记录4天过期

### 场景2：服务重启
1. 发送14天告警后服务重启
2. 重启后检查Redis，发现已发送14天告警
3. 不会重复发送14天告警

### 场景3：证书续期
1. 证书续期后，剩余天数增加到90天
2. 自动清理所有过期的告警状态
3. 下次到期前会重新发送告警

### 场景4：Redis故障处理
1. Redis连接失败时，`isAlertSent`返回false
2. 允许发送告警，避免遗漏重要通知
3. Redis恢复后正常记录状态

## 配置建议

在Redis配置中启用键过期通知（可选）：
```yaml
redis:
  # 其他配置...
  notify-keyspace-events: "Ex"  # 启用过期事件通知
```

## 监控建议

可以通过Redis监控告警状态：
```bash
# 查看所有SSL告警键
redis-cli KEYS "ssl_cert_alert:*"

# 查看特定域名的告警状态
redis-cli KEYS "ssl_cert_alert:asecondchance.cn:*"

# 查看键的TTL
redis-cli TTL "ssl_cert_alert:asecondchance.cn:14_days"
```

这个改进确保了SSL证书监控系统的可靠性和准确性，避免了重复告警的问题，同时支持服务重启和集群部署。
