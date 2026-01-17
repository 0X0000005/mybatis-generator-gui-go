#!/bin/bash
# MyBatis Generator GUI - 完整工作流脚本
# 包含: 测试 -> 编译 -> Git提交 -> 推送

VERSION="1.1.0"

echo "================================================"
echo "MyBatis Generator GUI - 完整工作流"
echo "版本: $VERSION"
echo "================================================"
echo ""

# 进入项目根目录
cd "$(dirname "$0")/.."

# ========== 1. 运行单元测试 ==========
echo "[1/5] 运行单元测试..."
go test ./...
if [ $? -ne 0 ]; then
    echo "错误: 单元测试失败"
    exit 1
fi
echo "✓ 单元测试通过"

echo ""
# ========== 2. 编译 ==========
echo "[2/5] 执行编译..."
bash build/build.sh
if [ $? -ne 0 ]; then
    echo "错误: 编译失败"
    exit 1
fi

echo ""
# ========== 3. Git状态检查 ==========
echo "[3/5] 检查Git状态..."
git status --short
echo ""

# ========== 4. Git提交 ==========
echo "[4/5] Git提交..."
read -p "请输入提交信息 (留空取消): " COMMIT_MSG
if [ -z "$COMMIT_MSG" ]; then
    echo "取消提交，工作流结束"
    exit 0
fi

git add .
git commit -m "$COMMIT_MSG"
if [ $? -ne 0 ]; then
    echo "提示: Git提交失败或无更改"
fi

echo ""
# ========== 5. 推送到远程 ==========
echo "[5/5] 推送到远程..."
read -p "是否推送到远程? (y/n): " PUSH
if [ "$PUSH" = "y" ] || [ "$PUSH" = "Y" ]; then
    git push
    if [ $? -ne 0 ]; then
        echo "错误: 推送失败"
        exit 1
    fi
    echo "✓ 推送成功"
fi

echo ""
echo "================================================"
echo "工作流完成!"
echo "================================================"
