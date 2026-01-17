#!/bin/bash
# MyBatis Generator GUI - 跨平台构建脚本

VERSION="1.1.0"
APP_NAME="mybatis-generator-gui"

echo "================================================"
echo "MyBatis Generator GUI - 构建脚本"
echo "版本: $VERSION"
echo "================================================"
echo ""

echo "[1/4] 清理旧文件..."
rm -f ${APP_NAME}-windows-amd64.exe ${APP_NAME}-linux-amd64

echo ""
echo "[2/4] 准备依赖包..."
go mod tidy
if [ $? -ne 0 ]; then
    echo "错误: 依赖包下载失败"
    exit 1
fi

echo ""
echo "[3/4] 编译Linux版本..."
GOOS=linux GOARCH=amd64 go build -ldflags "-s -w -X main.version=$VERSION" -o ${APP_NAME}-linux-amd64 ./cmd/main.go
if [ $? -ne 0 ]; then
    echo "错误: Linux编译失败"
    exit 1
fi
echo "✓ Linux版本编译完成"

echo ""
echo "[4/4] 编译Windows版本..."
GOOS=windows GOARCH=amd64 go build -ldflags "-s -w -X main.version=$VERSION" -o ${APP_NAME}-windows-amd64.exe ./cmd/main.go
if [ $? -ne 0 ]; then
    echo "错误: Windows编译失败"
    exit 1
fi
echo "✓ Windows版本编译完成"

echo ""
echo "================================================"
echo "构建完成!"
echo "================================================"
echo "Windows: ${APP_NAME}-windows-amd64.exe"
echo "Linux:   ${APP_NAME}-linux-amd64"
echo ""
echo "启动命令 (Linux): ./${APP_NAME}-linux-amd64"
echo "访问地址: http://localhost:8080"
echo "================================================"
