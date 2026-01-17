#!/bin/bash
# MyBatis Generator GUI Linux构建脚本

VERSION="1.0.0"
APP_NAME="mybatis-generator-gui"

echo "================================================"
echo "MyBatis Generator GUI Go - Linux构建脚本"
echo "版本: $VERSION"
echo "================================================"
echo ""

# 创建bin目录
mkdir -p bin

echo "[1/4] 准备依赖包..."
go mod tidy
if [ $? -ne 0 ]; then
    echo "错误: 依赖包下载失败"
    exit 1
fi

echo ""
echo "[2/4] 编译Linux版本..."
GOOS=linux GOARCH=amd64 go build -ldflags "-s -w -X main.version=$VERSION" -o bin/${APP_NAME}-linux-amd64 cmd/main.go
if [ $? -ne 0 ]; then
    echo "错误: 编译失败"
    exit 1
fi

echo ""
echo "[3/4] 检查UPX压缩工具..."
if ! command -v upx &> /dev/null; then
    echo "警告: 未找到UPX工具,跳过压缩步骤"
    echo "提示: 可通过包管理器安装UPX: sudo apt-get install upx"
else
    echo "[4/4] 压缩可执行文件..."
    upx --best --lzma bin/${APP_NAME}-linux-amd64
fi

echo ""
echo "================================================"
echo "构建完成!"
echo "可执行文件位置: bin/${APP_NAME}-linux-amd64"
echo "================================================"
