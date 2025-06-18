# Step 4: 性能优化与部署 (v0.4.x)

## 📋 阶段概述
- **阶段目标**: 优化系统性能，完善监控体系，准备生产环境部署
- **预计时间**: 2025-03-03 ~ 2025-03-16 (2周)
- **关键交付**: 高性能系统、完整监控、生产就绪的部署方案

## 🎯 详细任务

### 任务10: 性能优化
- **分支**: `xucheng/feature/v0.4/performance-optimization`
- **负责人**: xucheng
- **预计时间**: 4-5天

#### 开发目标
- [ ] Redis缓存策略优化  
- [ ] 数据库查询优化
- [ ] API响应时间优化
- [ ] 内存使用优化

#### 详细任务清单
- [ ] 实现多层缓存架构
- [ ] 优化热点数据缓存策略
- [ ] 数据库连接池优化
- [ ] SQL查询性能优化
- [ ] API接口响应时间优化
- [ ] 内存泄漏检测和修复
- [ ] 并发处理性能提升
- [ ] 资源使用监控和告警

#### 交付物
- [ ] 缓存优化方案 (`internal/cache/`)
- [ ] 数据库优化脚本 (`scripts/db_optimization/`)
- [ ] 性能测试报告 (`docs/performance/`)
- [ ] 优化建议文档 (`docs/optimization_guide.md`)
- [ ] 性能监控仪表板配置
- [ ] 负载测试工具和脚本

#### 缓存优化策略
```go
// 多层缓存架构
type CacheManager struct {
    L1Cache *sync.Map          // 内存缓存 (最热数据)
    L2Cache *redis.Client      // Redis缓存 (热数据)
    L3Cache *sql.DB           // 数据库 (持久化)
}

// 缓存策略配置
type CacheConfig struct {
    UserData    CachePolicy   // 用户数据缓存策略
    CharacterData CachePolicy // 角色数据缓存策略
    GameState   CachePolicy   // 游戏状态缓存策略
    Achievements CachePolicy  // 成就数据缓存策略
}

type CachePolicy struct {
    TTL         time.Duration  // 缓存过期时间
    MaxSize     int           // 最大缓存大小
    EvictPolicy string        // 淘汰策略 (LRU, LFU, FIFO)
    Compress    bool          // 是否压缩
}
```

#### 数据库优化
```sql
-- 索引优化
CREATE INDEX CONCURRENTLY idx_users_username_active ON users(username) WHERE is_active = true;
CREATE INDEX CONCURRENTLY idx_characters_user_active ON characters(user_id) WHERE is_active = true;
CREATE INDEX CONCURRENTLY idx_events_character_age ON events(character_id, age DESC);
CREATE INDEX CONCURRENTLY idx_relationships_character_type ON relationships(character_id, relation_type);

-- 分区表优化（事件表按年份分区）
CREATE TABLE events_2024 PARTITION OF events FOR VALUES FROM ('2024-01-01') TO ('2025-01-01');
CREATE TABLE events_2025 PARTITION OF events FOR VALUES FROM ('2025-01-01') TO ('2026-01-01');

-- 查询优化
EXPLAIN ANALYZE SELECT * FROM characters c 
JOIN character_attributes ca ON c.character_id = ca.character_id 
WHERE c.user_id = $1 AND c.is_active = true;
```

#### 性能指标监控
```go
type PerformanceMetrics struct {
    APIResponseTime map[string]time.Duration  // API响应时间
    DatabaseLatency map[string]time.Duration  // 数据库延迟
    CacheHitRate    map[string]float64        // 缓存命中率
    MemoryUsage     ResourceUsage             // 内存使用情况
    CPUUsage        ResourceUsage             // CPU使用情况
    ConcurrentUsers int                       // 并发用户数
}
```

---

### 任务11: 监控与日志
- **分支**: `xucheng/feature/v0.4/monitoring-logging`
- **负责人**: xucheng
- **预计时间**: 3-4天

#### 开发目标
- [ ] 结构化日志系统
- [ ] 性能监控集成
- [ ] 错误追踪和报警
- [ ] API调用统计

#### 详细任务清单
- [ ] 集成Prometheus监控系统
- [ ] 配置Grafana仪表板
- [ ] 实现结构化日志记录
- [ ] 设置错误追踪和报警
- [ ] 实现API调用链路追踪
- [ ] 配置日志轮转和存储
- [ ] 实现性能指标收集
- [ ] 建立监控告警规则

#### 交付物
- [ ] 日志系统 (`internal/logger/`)
- [ ] 监控仪表板 (`configs/grafana/`)
- [ ] 报警配置 (`configs/alerting/`)
- [ ] 统计分析工具 (`internal/metrics/`)
- [ ] 日志分析脚本 (`scripts/log_analysis/`)
- [ ] 监控运维手册 (`docs/monitoring_guide.md`)

#### 监控架构
```go
// 指标收集器
type MetricsCollector struct {
    prometheus.Collector
    requestCounter   *prometheus.CounterVec
    requestDuration  *prometheus.HistogramVec
    activeUsers      prometheus.Gauge
    errorRate        *prometheus.CounterVec
}

// 日志结构
type LogEntry struct {
    Timestamp   time.Time         `json:"timestamp"`
    Level       string           `json:"level"`
    Message     string           `json:"message"`
    Service     string           `json:"service"`
    UserID      string           `json:"user_id,omitempty"`
    RequestID   string           `json:"request_id,omitempty"`
    Duration    time.Duration    `json:"duration,omitempty"`
    Error       string           `json:"error,omitempty"`
    Fields      map[string]interface{} `json:"fields,omitempty"`
}
```

#### Grafana仪表板配置
```json
{
  "dashboard": {
    "title": "重启人生后端监控",
    "panels": [
      {
        "title": "API响应时间",
        "type": "graph",
        "targets": [
          {
            "expr": "histogram_quantile(0.95, rate(http_request_duration_seconds_bucket[5m]))"
          }
        ]
      },
      {
        "title": "错误率",
        "type": "stat",
        "targets": [
          {
            "expr": "rate(http_requests_total{status=~\"5..\"}[5m]) / rate(http_requests_total[5m])"
          }
        ]
      },
      {
        "title": "活跃用户数",
        "type": "stat",
        "targets": [
          {
            "expr": "active_users"
          }
        ]
      }
    ]
  }
}
```

#### 告警规则
```yaml
groups:
  - name: restart_life_alerts
    rules:
      - alert: HighErrorRate
        expr: rate(http_requests_total{status=~"5.."}[5m]) > 0.1
        for: 5m
        labels:
          severity: critical
        annotations:
          summary: "高错误率告警"
          description: "API错误率超过10%"
          
      - alert: SlowResponse
        expr: histogram_quantile(0.95, rate(http_request_duration_seconds_bucket[5m])) > 1
        for: 5m
        labels:
          severity: warning
        annotations:
          summary: "API响应慢"
          description: "95%的请求响应时间超过1秒"
```

---

### 任务12: 部署与文档
- **分支**: `xucheng/feature/v0.4/deployment-docs`
- **负责人**: xucheng
- **预计时间**: 4-5天

#### 开发目标
- [ ] Docker容器化
- [ ] API文档生成
- [ ] 部署脚本编写
- [ ] 生产环境配置

#### 详细任务清单
- [ ] 编写多阶段Docker构建文件
- [ ] 配置Docker Compose部署
- [ ] 编写Kubernetes部署清单
- [ ] 实现CI/CD管道配置
- [ ] 生成完整API文档
- [ ] 编写部署和运维文档
- [ ] 配置生产环境参数
- [ ] 实现健康检查和就绪探针

#### 交付物
- [ ] Docker镜像 (`docker/`)
- [ ] 完整API文档 (`docs/api/`)
- [ ] 自动化部署脚本 (`scripts/deploy/`)
- [ ] 运维手册 (`docs/operations/`)
- [ ] CI/CD配置 (`.github/workflows/`)
- [ ] Kubernetes清单 (`k8s/`)

#### Docker配置
```dockerfile
# 多阶段构建
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o restart-life-api ./cmd/server

FROM alpine:latest
RUN apk --no-cache add ca-certificates tzdata
WORKDIR /root/
COPY --from=builder /app/restart-life-api .
COPY --from=builder /app/configs ./configs
COPY --from=builder /app/migrations ./migrations

EXPOSE 8080
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
  CMD wget --quiet --tries=1 --spider http://localhost:8080/health || exit 1

CMD ["./restart-life-api"]
```

#### Docker Compose配置
```yaml
version: '3.8'
services:
  api:
    build: .
    ports:
      - "8080:8080"
    environment:
      - DB_HOST=postgres
      - REDIS_HOST=redis
    depends_on:
      - postgres
      - redis
    healthcheck:
      test: ["CMD", "wget", "--quiet", "--tries=1", "--spider", "http://localhost:8080/health"]
      interval: 30s
      timeout: 10s
      retries: 3

  postgres:
    image: postgres:15-alpine
    environment:
      POSTGRES_DB: restart_life
      POSTGRES_USER: postgres
      POSTGRES_PASSWORD: password
    volumes:
      - mysql_data:/var/lib/mysql
      - ./migrations:/docker-entrypoint-initdb.d
    ports:
      - "5432:5432"

  redis:
    image: redis:7.0-alpine
    ports:
      - "6379:6379"
    command: redis-server --appendonly yes
    volumes:
      - redis_data:/data

  prometheus:
    image: prom/prometheus:latest
    ports:
      - "9090:9090"
    volumes:
      - ./configs/prometheus.yml:/etc/prometheus/prometheus.yml

  grafana:
    image: grafana/grafana:latest
    ports:
      - "3000:3000"
    environment:
      - GF_SECURITY_ADMIN_PASSWORD=admin
    volumes:
      - grafana_data:/var/lib/grafana

volumes:
  postgres_data:
  redis_data:
  grafana_data:
```

#### Kubernetes部署清单
```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: restart-life-api
  labels:
    app: restart-life-api
spec:
  replicas: 3
  selector:
    matchLabels:
      app: restart-life-api
  template:
    metadata:
      labels:
        app: restart-life-api
    spec:
      containers:
      - name: api
        image: restart-life-api:latest
        ports:
        - containerPort: 8080
        env:
        - name: DB_HOST
          value: "postgres-service"
        - name: REDIS_HOST
          value: "redis-service"
        livenessProbe:
          httpGet:
            path: /health
            port: 8080
          initialDelaySeconds: 30
          periodSeconds: 10
        readinessProbe:
          httpGet:
            path: /ready
            port: 8080
          initialDelaySeconds: 5
          periodSeconds: 5
        resources:
          requests:
            memory: "128Mi"
            cpu: "100m"
          limits:
            memory: "512Mi"
            cpu: "500m"
---
apiVersion: v1
kind: Service
metadata:
  name: restart-life-api-service
spec:
  selector:
    app: restart-life-api
  ports:
  - protocol: TCP
    port: 80
    targetPort: 8080
  type: LoadBalancer
```

#### CI/CD管道配置
```yaml
name: Build and Deploy

on:
  push:
    branches: [main]
  pull_request:
    branches: [main]

jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - uses: actions/setup-go@v3
        with:
          go-version: 1.21
      - name: Run tests
        run: |
          go test -v -race -coverprofile=coverage.out ./...
          go tool cover -html=coverage.out -o coverage.html
      - name: Upload coverage
        uses: codecov/codecov-action@v3

  build:
    needs: test
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - name: Build Docker image
        run: |
          docker build -t restart-life-api:${{ github.sha }} .
          docker tag restart-life-api:${{ github.sha }} restart-life-api:latest

  deploy:
    needs: [test, build]
    runs-on: ubuntu-latest
    if: github.ref == 'refs/heads/main'
    steps:
      - name: Deploy to production
        run: |
          # 部署脚本
          echo "部署到生产环境"
```

---

## 📊 阶段验收标准

### 性能验收
- [ ] API响应时间P95 < 200ms
- [ ] 数据库查询延迟 < 50ms
- [ ] 缓存命中率 > 80%
- [ ] 并发用户支持 > 1000
- [ ] 内存使用稳定，无泄漏
- [ ] CPU使用率 < 70%

### 监控验收
- [ ] 所有关键指标被监控
- [ ] 告警规则配置正确
- [ ] 日志结构化且可搜索
- [ ] 监控仪表板完整
- [ ] 错误追踪系统正常
- [ ] 性能趋势分析可用

### 部署验收
- [ ] Docker镜像构建成功
- [ ] 容器化部署正常
- [ ] 健康检查机制有效
- [ ] 滚动更新策略正确
- [ ] 负载均衡配置正确
- [ ] 备份恢复流程完整

### 文档验收
- [ ] API文档完整准确
- [ ] 部署文档详细清晰
- [ ] 运维手册实用
- [ ] 故障排除指南完整
- [ ] 性能调优建议实用

---

## 🧪 测试规划

### 性能测试
```bash
# 负载测试
wrk -t12 -c400 -d30s --script=load_test.lua http://localhost:8080/api/v1/characters

# 压力测试
artillery run stress_test.yml

# 内存泄漏测试
go test -memprofile=mem.prof -bench=.
go tool pprof mem.prof
```

### 监控测试
```go
// 监控指标测试
func TestMetricsCollection(t *testing.T) {
    collector := NewMetricsCollector()
    
    // 模拟API调用
    collector.RecordRequest("GET", "/api/v1/characters", 200, time.Millisecond*150)
    
    // 验证指标收集
    assert.Equal(t, 1, collector.RequestCount("GET", "/api/v1/characters"))
}

// 告警测试
func TestAlertRules(t *testing.T) {
    // 测试告警规则触发
}
```

### 部署测试
```bash
# Docker测试
docker build -t restart-life-api:test .
docker run --rm -p 8080:8080 restart-life-api:test

# Kubernetes测试
kubectl apply -k k8s/
kubectl get pods -l app=restart-life-api
kubectl port-forward svc/restart-life-api-service 8080:80
```

---

## 🔧 技术规范

### 性能指标
```go
// 关键性能指标
const (
    MaxAPIResponseTime = 200 * time.Millisecond  // API响应时间上限
    MaxDBQueryTime     = 50 * time.Millisecond   // DB查询时间上限
    MinCacheHitRate    = 0.8                     // 缓存命中率下限
    MaxMemoryUsage     = 512 * 1024 * 1024       // 内存使用上限 512MB
    MaxCPUUsage        = 0.7                     // CPU使用率上限 70%
)
```

### 监控标准
```yaml
# 监控数据保留策略
retention:
  metrics: 30d      # 指标数据保留30天
  logs: 7d          # 日志数据保留7天
  traces: 3d        # 链路追踪数据保留3天

# 告警级别
alert_levels:
  critical: 立即处理    # 系统不可用
  warning: 及时处理     # 性能下降
  info: 关注处理        # 一般信息
```

### 部署标准
```yaml
# 资源配置标准
resources:
  api:
    requests:
      cpu: 100m
      memory: 128Mi
    limits:
      cpu: 500m
      memory: 512Mi
  
  postgres:
    requests:
      cpu: 200m
      memory: 256Mi
    limits:
      cpu: 1000m
      memory: 1Gi
```

---

## 📚 文档结构

```
docs/
├── api/                    # API文档
│   ├── openapi.yaml       # OpenAPI规范
│   ├── postman.json       # Postman集合
│   └── examples/          # API调用示例
├── deployment/            # 部署文档
│   ├── docker.md          # Docker部署
│   ├── kubernetes.md      # K8s部署
│   └── cloud.md           # 云平台部署
├── operations/            # 运维文档
│   ├── monitoring.md      # 监控指南
│   ├── logging.md         # 日志管理
│   ├── backup.md          # 备份恢复
│   └── troubleshooting.md # 故障排除
├── performance/           # 性能文档
│   ├── benchmarks.md      # 性能基准
│   ├── optimization.md    # 优化指南
│   └── load_testing.md    # 负载测试
└── security/              # 安全文档
    ├── security_guide.md  # 安全指南
    ├── auth.md            # 认证授权
    └── vulnerabilities.md # 漏洞管理
```

---

## ⏭️ 项目完成

完成Step 4后，整个《重启人生》后端项目将达到生产就绪状态：

### ✅ 已实现功能
- 完整的用户认证系统
- 强大的角色管理系统
- 复杂的游戏核心逻辑
- 丰富的关系网络系统
- 完善的成就系统
- 可靠的存档统计系统
- 高性能的缓存系统
- 全面的监控体系
- 生产级的部署方案

### 🚀 上线准备
- 性能优化完成
- 监控告警就绪
- 部署流程验证
- 文档完整齐全
- 测试覆盖充分

---

*创建时间: 2025-01-26*
*最后更新: 2025-01-26* 