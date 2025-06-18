# Docker 开发脚本

这个目录包含了用于Docker开发环境的便捷脚本。

## Windows 用户

### 基本用法
```cmd
scripts\build.bat       # 构建Docker镜像（官方源）
scripts\start.bat       # 启动所有服务（官方源）
scripts\stop.bat        # 停止所有服务
```

### 腾讯云镜像版本
```cmd
scripts\build.bat tencent    # 使用腾讯云镜像构建
scripts\start.bat tencent    # 使用腾讯云镜像启动
```

### 网络诊断
```cmd
scripts\fix-network.bat     # 诊断网络问题并给出建议
```

## Linux/macOS 用户

首先给脚本添加执行权限：
```bash
chmod +x scripts/*.sh
```

### 基本用法
```bash
./scripts/build.sh       # 构建Docker镜像（官方源）
./scripts/start.sh       # 启动所有服务（官方源）
./scripts/stop.sh        # 停止所有服务
```

### 腾讯云镜像版本
```bash
./scripts/build.sh tencent    # 使用腾讯云镜像构建
./scripts/start.sh tencent    # 使用腾讯云镜像启动
```

### 网络诊断
```bash
./scripts/fix-network.sh     # 诊断网络问题并给出建议
```

## 镜像源说明

- **官方源**: 使用Docker Hub官方镜像（默认）
- **腾讯云**: 使用腾讯云容器镜像服务

## 使用建议

1. **网络正常**：使用官方源即可
2. **网络较慢或无法连接**：使用腾讯云镜像
3. **不确定网络状况**：运行网络诊断脚本获取建议

## 故障排除

1. **构建失败**：
   ```cmd
   scripts\fix-network.bat  # 运行网络诊断
   scripts\build.bat tencent # 尝试腾讯云镜像
   ```

2. **启动失败**：
   ```cmd
   docker-compose logs      # 查看详细日志
   scripts\start.bat tencent # 尝试腾讯云镜像
   ```

详细说明请参考项目根目录的 `DOCKER.md` 文件。
