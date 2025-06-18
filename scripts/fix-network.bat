@echo off
REM 网络问题解决脚本 - fix-network.bat

echo === Docker网络问题修复工具 ===

REM 检查Docker是否运行
docker info >nul 2>&1
if errorlevel 1 (
    echo ❌ Docker未运行，请先启动Docker
    pause
    exit /b 1
)

echo 🔍 正在诊断网络问题...

REM 测试Docker Hub连接
echo 测试Docker Hub连接...
docker pull hello-world >nul 2>&1
if errorlevel 1 (
    echo ❌ Docker Hub连接失败
    set DOCKER_HUB_OK=false
) else (
    echo ✅ Docker Hub连接正常
    set DOCKER_HUB_OK=true
)

REM 测试阿里云镜像
echo 测试腾讯云镜像连接...
docker pull mirror.ccs.tencentyun.com/library/hello-world >nul 2>&1
if errorlevel 1 (
    echo ❌ 腾讯云镜像连接失败
    set TENCENT_OK=false
) else (
    echo ✅ 腾讯云镜像连接正常
    set TENCENT_OK=true
)

echo.
echo === 诊断结果 ===

if "%DOCKER_HUB_OK%"=="true" (
    echo 推荐使用：官方版本
    echo 命令：scripts\build.bat
) else if "%TENCENT_OK%"=="true" (
    echo 推荐使用：腾讯云版本
    echo 命令：scripts\build.bat tencent
) else (
    echo ⚠️  所有镜像源都无法连接
    echo 建议：
    echo 1. 检查网络连接
    echo 2. 配置Docker镜像加速器（参考MIRROR-SETUP.md）
    echo 3. 使用VPN或代理
    echo 4. 联系网络管理员
)

echo.
echo === 镜像加速器配置建议 ===
echo 在Docker Desktop的设置中添加以下镜像加速器：
echo.
echo {
echo   "registry-mirrors": [
echo     "https://mirror.ccs.tencentyun.com",
echo     "https://docker.mirrors.ustc.edu.cn",
echo     "https://hub-mirror.c.163.com",
echo     "https://registry.cn-hangzhou.aliyuncs.com"
echo   ]
echo }

echo.
echo === 手动配置步骤 ===
echo 1. 打开Docker Desktop
echo 2. 点击Settings（设置）
echo 3. 选择Docker Engine
echo 4. 在JSON配置中添加上述镜像加速器
echo 5. 点击Apply ^& Restart

pause
