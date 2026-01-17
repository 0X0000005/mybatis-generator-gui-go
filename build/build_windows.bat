@echo off
chcp 65001>nul
REM MyBatis Generator GUI 完整构建和部署脚本

SET VERSION=1.1.0
SET APP_NAME=mybatis-generator-gui

echo ================================================
echo MyBatis Generator GUI - 完整构建脚本
echo 版本: %VERSION%
echo ================================================
echo.

REM 返回项目根目录
cd /d %~dp0\..

echo [1/6] 清理旧文件...
if exist bin rmdir /s /q bin
mkdir bin

echo.
echo [2/6] 运行单元测试...
go test ./...
if errorlevel 1 (
    echo 错误: 单元测试失败
    pause
    exit /b 1
)

echo.
echo [3/6] 准备依赖包...
go mod tidy
if errorlevel 1 (
    echo 错误: 依赖包下载失败
    pause
    exit /b 1
)

echo.
echo [4/6] 编译Windows版本...
go build -ldflags "-s -w -H windowsgui -X main.version=%VERSION%" -o bin\%APP_NAME%-windows-amd64.exe .\cmd\main.go
if errorlevel 1 (
    echo 错误: 编译失败
    pause
    exit /b 1
)

echo.
echo [5/6] 检查UPX压缩工具...
where upx >nul 2>nul
if errorlevel 1 (
    echo 提示: 未找到UPX工具，跳过压缩步骤
    echo 下载地址: https://upx.github.io/
    goto :skip_upx
)

echo [6/6] 压缩可执行文件...
upx --best --lzma bin\%APP_NAME%-windows-amd64.exe
if errorlevel 1 (
    echo 警告: 压缩失败，但可执行文件已生成
)

:skip_upx
echo.
echo ================================================
echo 构建完成!
echo ================================================
echo 可执行文件: bin\%APP_NAME%-windows-amd64.exe
echo 启动命令: bin\%APP_NAME%-windows-amd64.exe
echo 访问地址: http://localhost:8080
echo ================================================
echo.
pause
