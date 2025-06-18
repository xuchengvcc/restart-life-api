@echo off
REM 启动脚本 - start.bat
setlocal enabledelayedexpansion

set PROJECT_NAME=restart-life-api

REM 检查镜像源选择
set COMPOSE_FILE=docker-compose.yml
set MIRROR_TYPE=official

if "%1"=="tencent" (
    set MIRROR_TYPE=tencent
    set COMPOSE_FILE=docker-compose.tencent.yml
    echo Using Tencent Cloud mirror for faster startup...
)

echo === Restart Life API Docker Start Script ===
echo Using mirror: %MIRROR_TYPE%

REM 检查Docker是否运行
docker info >nul 2>&1
if errorlevel 1 (
    echo Error: Docker is not running. Please start Docker and try again.
    pause
    exit /b 1
)

REM 检查docker-compose是否存在
docker-compose --version >nul 2>&1
if errorlevel 1 (
    echo Error: docker-compose is not installed or not in PATH
    pause
    exit /b 1
)

REM 检查compose文件是否存在
if not exist "%COMPOSE_FILE%" (    echo Error: %COMPOSE_FILE% not found
    echo Available options: start.bat [tencent]
    pause
    exit /b 1
)

REM 创建必要的目录
if not exist logs mkdir logs

REM 启动服务
echo Starting all services...
echo Using compose file: %COMPOSE_FILE%
docker-compose -f %COMPOSE_FILE% up -d

if errorlevel 1 (
    echo ❌ Failed to start services
    echo Checking logs for errors:
    docker-compose -f %COMPOSE_FILE% logs
    echo.
    echo === 故障排除建议 ===    if "%MIRROR_TYPE%"=="official" (
        echo 1. 尝试腾讯镜像: scripts\start.bat tencent
        echo 2. 运行网络诊断: scripts\fix-network.bat
        echo 3. 配置Docker镜像加速器（参考MIRROR-SETUP.md）
    ) else (
        echo 1. 检查网络连接
        echo 2. 运行网络诊断: scripts\fix-network.bat
        echo 3. 尝试官方镜像: scripts\start.bat
    )
    pause
    exit /b 1
)

echo ✅ Services started successfully!

REM 等待服务启动
echo Waiting for services to start...
timeout /t 15 /nobreak >nul

REM 显示服务状态
echo === Service Status ===
docker-compose -f %COMPOSE_FILE% ps

echo.
echo === Service URLs ===
echo 🚀 API Server:          http://localhost:8080
echo 🔍 Health Check:        http://localhost:8080/health
echo 🗄️  Database Admin:      http://localhost:8081
echo 🔴 Redis Commander:     http://localhost:8082

echo.
echo === Useful Commands ===
echo 📊 View logs:           docker-compose -f %COMPOSE_FILE% logs -f
echo 📊 View app logs:       docker-compose -f %COMPOSE_FILE% logs -f app
echo 🔄 Restart services:    docker-compose -f %COMPOSE_FILE% restart
echo 🛑 Stop services:       docker-compose -f %COMPOSE_FILE% down
echo 🧹 Clean up:            docker-compose -f %COMPOSE_FILE% down -v

echo.
echo === Recent Logs ===
docker-compose -f %COMPOSE_FILE% logs --tail=20

echo.
echo Press any key to continue...
pause >nul
