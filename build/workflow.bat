@echo off
chcp 65001>nul
REM MyBatis Generator GUI - 完整工作流脚本
REM 包含: 测试 -> 编译 -> Git提交 -> 推送

SET VERSION=1.1.0

echo ================================================
echo MyBatis Generator GUI - 完整工作流
echo 版本: %VERSION%
echo ================================================
echo.

REM 返回项目根目录
cd /d %~dp0\..

REM ========== 1. 运行单元测试 ==========
echo [1/5] 运行单元测试...
go test ./...
if errorlevel 1 (
    echo 错误: 单元测试失败
    pause
    exit /b 1
)
echo ✓ 单元测试通过

echo.
REM ========== 2. 编译 ==========
echo [2/5] 执行编译...
call build\build.bat
if errorlevel 1 (
    echo 错误: 编译失败
    pause
    exit /b 1
)

echo.
REM ========== 3. Git状态检查 ==========
echo [3/5] 检查Git状态...
git status --short
echo.

REM ========== 4. Git提交 ==========
echo [4/5] Git提交...
set /p COMMIT_MSG="请输入提交信息: "
if "%COMMIT_MSG%"=="" (
    echo 取消提交，工作流结束
    pause
    exit /b 0
)

git add .
git commit -m "%COMMIT_MSG%"
if errorlevel 1 (
    echo 提示: Git提交失败或无更改
)

echo.
REM ========== 5. 推送到远程 ==========
echo [5/5] 推送到远程...
set /p PUSH="是否推送到远程? (y/n): "
if /i "%PUSH%"=="y" (
    git push
    if errorlevel 1 (
        echo 错误: 推送失败
        pause
        exit /b 1
    )
    echo ✓ 推送成功
)

echo.
echo ================================================
echo 工作流完成!
echo ================================================
pause
