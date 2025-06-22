@echo off
REM 停止脚本 - stop.bat
setlocal enabledelayedexpansion

echo === Restart Life API Docker Stop Script ===

echo Stopping all running Docker containers...
for /f "tokens=1" %%i in ('docker ps -q') do docker stop %%i >nul 2>&1

echo Stopping all services...
REM 检查docker-compose是否存在
docker-compose --version >nul 2>&1
if errorlevel 1 (
    echo Error: docker-compose is not installed or not in PATH
    pause
    exit /b 1
)

docker-compose -f docker\docker-compose.yml down

if errorlevel 1 (
    echo ❌ Failed to stop services
    pause
    exit /b 1
)

echo ✅ All services stopped successfully!

REM 询问是否清理数据
echo.
set /p choice="Do you want to remove volumes (this will delete all data)? (y/N): "
if /i "%choice%"=="y" (
    echo Removing volumes...
    docker-compose -f docker\docker-compose.yml down -v
    echo ✅ Volumes removed successfully!
) else (
    echo Data volumes preserved.
)

echo.
echo === Other Useful Commands ===
echo 🧹 Clean up everything:    docker-compose -f docker\docker-compose.yml down -v --remove-orphans
echo 📊 View stopped containers: docker-compose -f docker\docker-compose.yml ps -a
echo 🔍 Remove unused images:   docker image prune -f

echo.
pause
