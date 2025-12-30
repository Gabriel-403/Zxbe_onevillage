@echo off
echo 正在启动乡村振兴后端服务...
echo.

REM 检查Go是否安装
go version >nul 2>&1
if %errorlevel% neq 0 (
    echo 错误: 未找到Go环境，请先安装Go语言
    pause
    exit /b 1
)

REM 安装依赖
echo 正在安装依赖...
go mod tidy

REM 启动服务器
echo 正在启动服务器...
echo 服务器将在 http://localhost:8080 启动
echo 按 Ctrl+C 停止服务器
echo.
go run main.go

pause
