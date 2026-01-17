#!/bin/bash
# MyBatis Generator GUI - 跨平台构建脚本

VERSION="1.1.0"
APP_NAME="mgg"

echo "================================================"
echo "MyBatis Generator GUI - 构建脚本"
echo "版本: $VERSION"
echo "================================================"
echo ""

echo "[1/5] 清理旧文件..."
rm -f ${APP_NAME}-windows-amd64.exe ${APP_NAME}-linux-amd64

echo ""
echo "[2/5] 准备依赖包..."
go mod tidy
if [ $? -ne 0 ]; then
    echo "错误: 依赖包下载失败"
    exit 1
fi

echo ""
echo "[3/5] 编译Linux版本..."
GOOS=linux GOARCH=amd64 go build -ldflags "-s -w -X main.version=$VERSION" -o ${APP_NAME}-linux-amd64 ./cmd/main.go
if [ $? -ne 0 ]; then
    echo "错误: Linux编译失败"
    exit 1
fi
echo "✓ Linux版本编译完成"

echo ""
echo "[4/5] 编译Windows版本..."
GOOS=windows GOARCH=amd64 go build -ldflags "-s -w -X main.version=$VERSION" -o ${APP_NAME}-windows-amd64.exe ./cmd/main.go
if [ $? -ne 0 ]; then
    echo "错误: Windows编译失败"
    exit 1
fi
echo "✓ Windows版本编译完成"

echo ""
echo "[5/5] UPX压缩..."
if ! command -v upx &> /dev/null; then
    echo "警告: 未找到UPX工具，跳过压缩步骤"
    echo "安装方法: sudo apt-get install upx"
else
    echo "压缩Linux版本..."
    upx -9 ${APP_NAME}-linux-amd64
    if [ $? -ne 0 ]; then
        echo "警告: Linux版本压缩失败"
    fi
    
    echo "压缩Windows版本..."
    upx -9 ${APP_NAME}-windows-amd64.exe
    if [ $? -ne 0 ]; then
        echo "警告: Windows版本压缩失败"
    fi
fi

echo ""
echo "================================================"
echo "构建完成!"
echo "================================================"
echo "Windows: ${APP_NAME}.exe"
echo "Linux:   ${APP_NAME}"
echo ""
echo "启动命令 (Linux): ./${APP_NAME}"
echo "访问地址: http://localhost:8080"
echo "================================================"
