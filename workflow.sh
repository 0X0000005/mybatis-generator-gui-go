#!/bin/bash
# MyBatis Generator GUI - 完整工作流脚本
# 包含: 测试 -> 编译 -> Git提交 -> 推送 -> 标签

VERSION="1.1.0"

echo "================================================"
echo "MyBatis Generator GUI - 完整工作流"
echo "版本: $VERSION"
echo "================================================"
echo ""

# ========== 1. 运行单元测试 ==========
echo "[1/6] 运行单元测试..."
go test ./...
if [ $? -ne 0 ]; then
    echo "错误: 单元测试失败"
    exit 1
fi
echo "✓ 单元测试通过"

echo ""
# ========== 2. 编译 ==========
echo "[2/6] 执行编译..."
bash build.sh
if [ $? -ne 0 ]; then
    echo "错误: 编译失败"
    exit 1
fi

echo ""
# ========== 3. Git状态检查 ==========
echo "[3/6] 检查Git状态..."
git status --short
echo ""

# ========== 4. Git提交 ==========
echo "[4/6] Git提交..."
read -p "请输入提交信息 (留空跳过): " COMMIT_MSG
if [ -z "$COMMIT_MSG" ]; then
    echo "跳过提交"
else
    git add .
    git commit -m "$COMMIT_MSG"
    if [ $? -ne 0 ]; then
        echo "提示: Git提交失败或无更改"
    fi
fi

echo ""
# ========== 5. 推送到远程 ==========
echo "[5/6] 推送到远程..."
read -p "是否推送到远程? (y/n): " PUSH
if [ "$PUSH" = "y" ] || [ "$PUSH" = "Y" ]; then
    git push
    if [ $? -ne 0 ]; then
        echo "错误: 推送失败"
        exit 1
    fi
    echo "✓ 推送成功"
else
    echo "跳过推送"
fi

echo ""
# ========== 6. 创建版本标签 ==========
echo "[6/6] 创建版本标签..."
read -p "是否创建版本标签? (y/n): " CREATE_TAG
if [ "$CREATE_TAG" != "y" ] && [ "$CREATE_TAG" != "Y" ]; then
    echo "跳过标签创建"
else
    read -p "请输入标签名称 (例如: v$VERSION): " TAG_NAME
    if [ -z "$TAG_NAME" ]; then
        echo "取消标签创建"
    else
        read -p "请输入标签说明 (留空使用默认): " TAG_MSG
        if [ -z "$TAG_MSG" ]; then
            TAG_MSG="Release $TAG_NAME"
        fi
        
        # 检查是否有未提交的更改
        if ! git diff-index --quiet HEAD --; then
            echo "警告: 存在未提交的更改，请先提交"
            exit 1
        fi
        
        # 创建标签
        git tag -a "$TAG_NAME" -m "$TAG_MSG"
        if [ $? -ne 0 ]; then
            echo "错误: 创建标签失败"
            exit 1
        fi
        echo "✓ 标签 $TAG_NAME 创建成功"
        
        # 推送标签
        read -p "是否推送标签到远程? (y/n): " PUSH_TAG
        if [ "$PUSH_TAG" = "y" ] || [ "$PUSH_TAG" = "Y" ]; then
            git push origin "$TAG_NAME"
            if [ $? -ne 0 ]; then
                echo "错误: 推送标签失败"
                exit 1
            fi
            echo "✓ 标签推送成功"
        fi
    fi
fi

echo ""
echo "================================================"
echo "工作流完成!"
echo "================================================"
