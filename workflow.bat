@echo off
chcp 65001>nul
REM MyBatis Generator GUI - 完整工作流脚本
REM 包含: 测试 -> 编译 -> Git提交 -> 推送 -> 标签

SET VERSION=1.1.0

echo ================================================
echo MyBatis Generator GUI - 完整工作流
echo 版本: %VERSION%
echo ================================================
echo.

REM ========== 1. 运行单元测试 ==========
echo [1/6] 运行单元测试...
go test ./...
if errorlevel 1 (
    echo 错误: 单元测试失败
    pause
    exit /b 1
)
echo ✓ 单元测试通过

echo.
REM ========== 2. 编译 ==========
echo [2/6] 执行编译...
call build.bat
if errorlevel 1 (
    echo 错误: 编译失败
    pause
    exit /b 1
)

echo.
REM ========== 3. Git状态检查 ==========
echo [3/6] 检查Git状态...
git status --short
echo.

REM ========== 4. Git提交 ==========
echo [4/6] Git提交...
set /p COMMIT_MSG="请输入提交信息 (留空跳过): "
if "%COMMIT_MSG%"=="" (
    echo 跳过提交
    goto :push_skip
)

git add .
git commit -m "%COMMIT_MSG%"
if errorlevel 1 (
    echo 提示: Git提交失败或无更改
)

echo.
REM ========== 5. 推送到远程 ==========
:push_skip
echo [5/6] 推送到远程...
set /p PUSH="是否推送到远程? (y/n): "
if /i not "%PUSH%"=="y" (
    echo 跳过推送
    goto :tag_skip
)

git push
if errorlevel 1 (
    echo 错误: 推送失败
    pause
    exit /b 1
)
echo ✓ 推送成功

echo.
REM ========== 6. 创建版本标签 ==========
:tag_skip
echo [6/6] 创建版本标签...
set /p CREATE_TAG="是否创建版本标签? (y/n): "
if /i not "%CREATE_TAG%"=="y" (
    echo 跳过标签创建
    goto :end
)

set /p TAG_NAME="请输入标签名称 (例如: v%VERSION%): "
if "%TAG_NAME%"=="" (
    echo 取消标签创建
    goto :end
)

set /p TAG_MSG="请输入标签说明 (留空使用默认): "
if "%TAG_MSG%"=="" (
    set TAG_MSG=Release %TAG_NAME%
)

REM 检查是否有未提交的更改
git diff-index --quiet HEAD
if errorlevel 1 (
    echo 警告: 存在未提交的更改，请先提交
    pause
    exit /b 1
)

REM 创建标签
git tag -a %TAG_NAME% -m "%TAG_MSG%"
if errorlevel 1 (
    echo 错误: 创建标签失败
    pause
    exit /b 1
)
echo ✓ 标签 %TAG_NAME% 创建成功

REM 推送标签
set /p PUSH_TAG="是否推送标签到远程? (y/n): "
if /i "%PUSH_TAG%"=="y" (
    git push origin %TAG_NAME%
    if errorlevel 1 (
        echo 错误: 推送标签失败
        pause
        exit /b 1
    )
    echo ✓ 标签推送成功
)

:end
echo.
echo ================================================
echo 工作流完成!
echo ================================================
pause
