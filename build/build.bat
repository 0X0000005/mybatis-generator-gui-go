@echo off
chcp 65001>nul
REM MyBatis Generator GUI - 跨平台构建脚本

SET VERSION=1.1.0
SET APP_NAME=mybatis-generator-gui

echo ================================================
echo MyBatis Generator GUI - 构建脚本
echo 版本: %VERSION%
echo ================================================
echo.

REM 返回项目根目录
cd /d %~dp0\..

echo [1/4] 清理旧文件...
if exist bin rmdir /s /q bin
mkdir bin

echo.
echo [2/4] 准备依赖包...
go mod tidy
if errorlevel 1 (
    echo 错误: 依赖包下载失败
    pause
    exit /b 1
)

echo.
echo [3/4] 编译Windows版本...
go build -ldflags "-s -w -X main.version=%VERSION%" -o bin\%APP_NAME%-windows-amd64.exe .\cmd\main.go
if errorlevel 1 (
    echo 错误: Windows编译失败
    pause
    exit /b 1
)
echo ✓ Windows版本编译完成

echo.
echo [4/4] 编译Linux版本...
set GOOS=linux
set GOARCH=amd64
go build -ldflags "-s -w -X main.version=%VERSION%" -o bin\%APP_NAME%-linux-amd64 .\cmd\main.go
if errorlevel 1 (
    echo 错误: Linux编译失败
    pause
    exit /b 1
)
echo ✓ Linux版本编译完成

echo.
echo ================================================
echo 构建完成!
echo ================================================
echo Windows: bin\%APP_NAME%-windows-amd64.exe
echo Linux:   bin\%APP_NAME%-linux-amd64
echo.
echo 启动命令 (Windows): bin\%APP_NAME%-windows-amd64.exe
echo 访问地址: http://localhost:8080
echo ================================================
echo.
pause
