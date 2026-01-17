@echo off
chcp 65001>nul
REM MyBatis Generator GUI Windows构建脚本

SET VERSION=1.0.0
SET APP_NAME=mybatis-generator-gui

echo ================================================
echo MyBatis Generator GUI Go - Windows构建脚本
echo 版本: %VERSION%
echo ================================================
echo.

REM 创建bin目录
if not exist bin mkdir bin

echo [1/4] 准备依赖包...
go mod tidy
if errorlevel 1 (
    echo 错误: 依赖包下载失败
    pause
    exit /b 1
)

echo.
echo [2/4] 编译Windows版本...
go build -ldflags "-s -w -H windowsgui -X main.version=%VERSION%" -o bin/%APP_NAME%-windows-amd64.exe cmd/main.go
if errorlevel 1 (
    echo 错误: 编译失败
    pause
    exit /b 1
)

echo.
echo [3/4] 检查UPX压缩工具...
where upx >nul 2>nul
if errorlevel 1 (
    echo 警告: 未找到UPX工具,跳过压缩步骤
    echo 提示: 可从 https://upx.github.io/ 下载UPX以减小文件大小
    goto :finish
)

echo [4/4] 压缩可执行文件...
upx --best --lzma bin/%APP_NAME%-windows-amd64.exe
if errorlevel 1 (
    echo 警告: 压缩失败,但可执行文件已生成
)

:finish
echo.
echo ================================================
echo 构建完成!
echo 可执行文件位置: bin\%APP_NAME%-windows-amd64.exe
echo ================================================
pause
