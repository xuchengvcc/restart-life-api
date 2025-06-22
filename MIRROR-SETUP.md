# Docker镜像加速器配置指南

## 问题描述
如果遇到以下错误：
```
failed to solve: golang:1.23.8-alpine: failed to resolve source metadata
```

这是因为网络连接问题导致无法访问Docker Hub。

## 解决方案

### 1. 使用内置腾讯云镜像版本（推荐）

最简单的方法是使用项目内置的腾讯云镜像版本：

```cmd
# Windows用户
scripts\build.bat tencent
scripts\start.bat tencent

# Linux/macOS用户
./scripts/build.sh tencent
./scripts/start.sh tencent
```

### 2. 配置Docker Desktop镜像加速器

1. **打开Docker Desktop**
2. **点击设置图标**（齿轮图标）
3. **选择"Docker Engine"**
4. **在JSON配置中添加以下内容**：

```json
{
  "builder": {
    "gc": {
      "defaultKeepStorage": "20GB",
      "enabled": true
    }
  },
  "experimental": false,
  "features": {
    "buildkit": true
  },
  "registry-mirrors": [
    "https://mirror.ccs.tencentyun.com"
  ]
}
```

5. **点击"Apply & Restart"**
6. **等待Docker重启完成**

### 3. 验证配置

配置完成后，运行以下命令验证：

```cmd
# 测试连接
docker pull hello-world

# 查看配置
docker info | findstr -i mirror
```

### 4. 重新构建

配置完成后，重新运行构建命令：

```cmd
# 使用官方源
scripts\build.bat

# 使用腾讯云镜像
scripts\build.bat tencent
```

### 5. 网络诊断

使用内置的网络诊断工具：

```cmd
# Windows用户
scripts\fix-network.bat

# Linux/macOS用户
./scripts/fix-network.sh
```

该工具会自动检测网络状况并给出建议。

# 网易云镜像
scripts\build.bat 163

# 官方镜像（需要配置加速器）
scripts\build.bat
```

## 常见问题

**Q: 配置后仍然连接失败？**
A: 
1. 尝试使用腾讯云镜像版本：`scripts\build.bat tencent`
2. 重启Docker Desktop
3. 检查网络防火墙设置
4. 运行网络诊断：`scripts\fix-network.bat`

**Q: 如何验证镜像加速器是否生效？**
A: 运行 `docker info` 查看Registry Mirrors部分是否显示配置的镜像源。

**Q: 官方版本和腾讯云版本有什么区别？**
A: 
- **官方版本**: 使用Docker Hub官方镜像，速度可能较慢但最新
- **腾讯云版本**: 使用腾讯云镜像加速，在中国大陆访问更快

## 推荐使用方式

1. **首先尝试**: `scripts\build.bat`（官方版本）
2. **如果失败**: `scripts\build.bat tencent`（腾讯云版本）
3. **网络诊断**: `scripts\fix-network.bat`
4. **配置加速器**: 按照上述步骤配置Docker Desktop
