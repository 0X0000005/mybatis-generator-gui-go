#!/bin/bash
# Git版本标签创建脚本

if [ -z "$1" ]; then
    echo "使用方式: ./tag_version.sh v1.0.0"
    echo "示例: ./tag_version.sh v1.0.1"
    exit 1
fi

VERSION=$1

echo "创建版本标签: $VERSION"
echo "================================"

# 检查是否有未提交的更改
if  [[ -n $(git status -s) ]]; then
    echo "警告: 有未提交的更改"
    echo "请先提交所有更改后再创建标签"
    exit 1
fi

# 创建标签
git tag -a $VERSION -m "发布版本 $VERSION"
if [ $? -ne 0 ]; then
    echo "错误: 创建标签失败"
    exit 1
fi

echo "标签创建成功: $VERSION"

# 推送标签
read -p "是否推送标签到远程仓库? (y/n) " -n 1 -r
echo
if [[ $REPLY =~ ^[Yy]$ ]]; then
    git push origin $VERSION
    if [ $? -eq 0 ]; then
        echo "标签已推送到远程仓库"
    else
        echo "警告: 推送标签失败"
    fi
fi

echo "================================"
echo "完成!"
