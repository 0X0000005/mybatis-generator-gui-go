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

echo [1/5] 清理旧文件...
if exist %APP_NAME%-windows-amd64.exe del %APP_NAME%-windows-amd64.exe
if exist %APP_NAME%-linux-amd64 del %APP_NAME%-linux-amd64

echo.
echo [2/5] 准备依赖包...
go mod tidy
if errorlevel 1 (
    echo 错误: 依赖包下载失败
    pause
    exit /b 1
)

echo.
echo [3/5] 编译Windows版本...
set GOOS=windows
set GOARCH=amd64
go build -ldflags "-s -w -X main.version=%VERSION%" -o %APP_NAME%-windows-amd64.exe .\cmd\main.go
if errorlevel 1 (
    echo 错误: Windows编译失败
    pause
    exit /b 1
)
echo ✓ Windows版本编译完成

echo.
echo [4/5] 编译Linux版本...
set GOOS=linux
set GOARCH=amd64
go build -ldflags "-s -w -X main.version=%VERSION%" -o %APP_NAME%-linux-amd64 .\cmd\main.go
set GOOS=
set GOARCH=
if errorlevel 1 (
    echo 错误: Linux编译失败
    pause
    exit /b 1
)
echo ✓ Linux版本编译完成

echo.
echo [5/5] UPX压缩...
where upx >nul 2>nul
if errorlevel 1 (
    echo 警告: 未找到UPX工具，跳过压缩步骤
    echo 下载地址: https://upx.github.io/
    goto :finish
)

echo 压缩Windows版本...
upx -9 %APP_NAME%-windows-amd64.exe
if errorlevel 1 (
    echo 警告: Windows版本压缩失败
)

echo 压缩Linux版本...
upx -9 %APP_NAME%-linux-amd64
if errorlevel 1 (
    echo 警告: Linux版本压缩失败
)

:finish
echo.
echo ================================================
echo 构建完成!
echo ================================================
echo Windows: %APP_NAME%-windows-amd64.exe
echo Linux:   %APP_NAME%-linux-amd64
echo.
echo 启动命令 (Windows): %APP_NAME%-windows-amd64.exe
echo 访问地址: http://localhost:8080
echo ================================================
echo.
pause
