@echo off
REM 构建脚本 - build.bat
setlocal enabledelayedexpansion

REM 项目配置
set PROJECT_NAME=restart-life-api
set VERSION=%1
if "%VERSION%"=="" set VERSION=latest
set DOCKER_IMAGE=%PROJECT_NAME%:%VERSION%

REM 检查镜像源选择
set DOCKERFILE=Dockerfile
set MIRROR_TYPE=official

if "%1"=="tencent" (
    set MIRROR_TYPE=tencent
    set DOCKERFILE=Dockerfile.tencent
    set VERSION=latest
    set DOCKER_IMAGE=%PROJECT_NAME%:latest
    echo Using Tencent Cloud mirror for faster build...
)

echo === Restart Life API Docker Build Script ===
echo Building Docker image: %DOCKER_IMAGE%
echo Using mirror: %MIRROR_TYPE%

REM 检查Docker是否运行
docker info >nul 2>&1
if errorlevel 1 (
    echo Error: Docker is not running. Please start Docker and try again.
    pause
    exit /b 1
)

REM 检查Dockerfile是否存在
if not exist "%DOCKERFILE%" (    echo Error: %DOCKERFILE% not found
    echo Available options: build.bat [tencent]
    pause
    exit /b 1
)

REM 构建Docker镜像
echo Building Docker image...
echo Using Dockerfile: %DOCKERFILE%
docker build -f %DOCKERFILE% -t "%DOCKER_IMAGE%" .

if errorlevel 1 (
    echo ❌ Docker build failed
    echo.
    echo === 故障排除建议 ===    if "%MIRROR_TYPE%"=="official" (
        echo 1. 尝试腾讯镜像: scripts\build.bat tencent
        echo 2. 运行网络诊断: scripts\fix-network.bat
        echo 3. 配置Docker镜像加速器（参考MIRROR-SETUP.md）
    ) else (
        echo 1. 检查网络连接
        echo 2. 运行网络诊断: scripts\fix-network.bat
        echo 3. 尝试官方镜像: scripts\build.bat
    )
    pause
    exit /b 1
)

echo ✅ Docker image built successfully: %DOCKER_IMAGE%

REM 显示镜像信息
echo Image details:
docker images %PROJECT_NAME% --format "table {{.Repository}}\t{{.Tag}}\t{{.ID}}\t{{.CreatedAt}}\t{{.Size}}"

echo.
echo ✅ Build completed successfully!
echo Next steps:
if "%MIRROR_TYPE%"=="tencent" (
    echo   1. Run: scripts\start.bat tencent to start with Tencent mirrors
) else (
    echo   1. Run: scripts\start.bat to start the application
)
echo   2. Or run: docker-compose up -d to start all services
pause
