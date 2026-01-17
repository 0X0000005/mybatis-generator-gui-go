#!/bin/bash
# MyBatis Generator GUI 完整构建和部署脚本

VERSION="1.1.0"
APP_NAME="mybatis-generator-gui"

echo "================================================"
echo "MyBatis Generator GUI - 完整构建脚本"
echo "版本: $VERSION"
echo "================================================"
echo ""

# 进入项目根目录
cd "$(dirname "$0")/.."

echo "[1/6] 清理旧文件..."
rm -rf bin
mkdir -p bin

echo ""
echo "[2/6] 运行单元测试..."
go test ./...
if [ $? -ne 0 ]; then
    echo "错误: 单元测试失败"
    exit 1
fi

echo ""
echo "[3/6] 准备依赖包..."
go mod tidy
if [ $? -ne 0 ]; then
    echo "错误: 依赖包下载失败"
    exit 1
fi

echo ""
echo "[4/6] 编译Linux版本..."
GOOS=linux GOARCH=amd64 go build -ldflags "-s -w -X main.version=$VERSION" -o bin/${APP_NAME}-linux-amd64 ./cmd/main.go
if [ $? -ne 0 ]; then
    echo "错误: 编译失败"
    exit 1
fi

echo ""
echo "[5/6] 检查UPX压缩工具..."
if ! command -v upx &> /dev/null; then
    echo "提示: 未找到UPX工具，跳过压缩步骤"
    echo "安装方法: sudo apt-get install upx"
else
    echo "[6/6] 压缩可执行文件..."
    upx --best --lzma bin/${APP_NAME}-linux-amd64
fi

echo ""
echo "================================================"
echo "构建完成!"
echo "================================================"
echo "可执行文件: bin/${APP_NAME}-linux-amd64"
echo "启动命令: ./bin/${APP_NAME}-linux-amd64"
echo "访问地址: http://localhost:8080"
echo "================================================"
