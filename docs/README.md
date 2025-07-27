# 重启人生项目文档 (Restart Life Documentation)

## 📖 文档概述

本目录包含重启人生项目的完整技术文档，按照开发阶段分类整理，与项目开发计划紧密对应。

## 📁 文档结构

### 📋 开发计划
- [dev_plan/](./dev_plan/) - 项目开发计划和阶段规划
  - [开发计划.md](./dev_plan/开发计划.md) - 总体开发计划
  - [step1.md](./dev_plan/step1.md) - 第一阶段：基础架构搭建
  - [step2.md](./dev_plan/step2.md) - 第二阶段：核心业务模块
  - [step3.md](./dev_plan/step3.md) - 第三阶段：高级功能模块
  - [step4.md](./dev_plan/step4.md) - 第四阶段：性能优化与部署

### 🏗️ 开发阶段文档

#### [dev_stage1/](./dev_stage1/) - 基础架构搭建 (v0.1.x) ✅
- **状态**: 已完成
- **时间**: 2025-01-26 ~ 2025-02-02
- **文档**: 架构设计、分层实现、Redis集成
- **成果**: Go项目架构、数据库连接、基础服务

#### [dev_stage2/](./dev_stage2/) - 核心业务模块 (v0.2.x) ✅
- **状态**: 已完成
- **时间**: 2025-02-03 ~ 2025-02-16
- **文档**: 认证系统、游戏逻辑、前端集成
- **成果**: 用户系统、角色管理、游戏状态管理

#### [dev_stage3/](./dev_stage3/) - 高级功能模块 (v0.3.x) 🚧
- **状态**: 开发中
- **时间**: 2025-02-17 ~ 2025-03-02
- **文档**: AI系统、事件系统、关系网络
- **当前**: 游戏引擎开发 (`xucheng/feature/v0.2/game-engine`)

#### [dev_stage4/](./dev_stage4/) - 性能优化与部署 (v0.4.x) 📋
- **状态**: 计划中
- **时间**: 2025-03-03 ~ 2025-03-16
- **文档**: 性能优化、监控系统、生产部署
- **准备**: SSL证书管理已完成

## 🗂️ 文档分类

### 🔍 快速导航

| 类别 | 阶段 | 主要文档 | 状态 |
|------|------|----------|------|
| **架构设计** | Stage 1 | [分层架构](./dev_stage1/architecture_layers.md) | ✅ 完成 |
| **认证系统** | Stage 2 | [认证系统设计](./dev_stage2/auth_system.md) | ✅ 完成 |
| **游戏逻辑** | Stage 2 | [GameState架构](./dev_stage2/gamestate_storage_architecture.md) | ✅ 完成 |
| **AI系统** | Stage 3 | [AI优化文档](./dev_stage3/ai_prompt_enhancement.md) | 🚧 开发中 |
| **部署监控** | Stage 4 | [SSL证书监控](./dev_stage4/ssl_cert_monitor.md) | 📋 计划中 |

### 📊 开发进度

```
Stage 1: ████████████████████████ 100% (已完成)
Stage 2: ████████████████████████ 100% (已完成)
Stage 3: ████████████░░░░░░░░░░░░  60% (开发中)
Stage 4: ████░░░░░░░░░░░░░░░░░░░░  20% (计划中)
```

## 🎯 当前状态

### 📈 整体进度
- **已完成**: Stage 1 + Stage 2 (基础架构 + 核心业务)
- **开发中**: Stage 3 (高级功能模块)
- **当前分支**: `xucheng/feature/v0.2/game-engine`
- **下个里程碑**: 关系网络系统、成就系统

### 🔧 技术栈
- **后端**: Go + Gin + MySQL + Redis
- **前端**: React + TypeScript + Vite + Ant Design
- **部署**: Docker + Nginx + SSL
- **监控**: 计划中 (Prometheus + Grafana)

## 📚 其他资源

- [PRD产品需求文档](../prdtd/PRD.md)
- [后端技术设计文档](../prdtd/后端技术设计文档_Backend_TD.md)
- [开发规范](../regulations/regulation.md)

## 📞 文档维护

文档与代码同步更新，每个开发阶段完成后会更新对应的文档状态。如有问题请提交 Issue 或 PR。

---

📅 **最后更新**: 2025-07-27
🏷️ **当前版本**: v0.2.x (Stage 3 开发中)
