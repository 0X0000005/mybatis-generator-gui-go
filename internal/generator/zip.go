package generator

import (
	"archive/zip"
	"fmt"
	"io"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const (
	// TempDirName 临时目录名（相对于应用当前目录）
	TempDirName = "temp"
	// FileExpireDuration 文件过期时间（5分钟）
	FileExpireDuration = 5 * time.Minute
)

// CreateZipArchive 创建ZIP归档文件
func CreateZipArchive(files []string, projectFolder, tableName string) (string, error) {
	if len(files) == 0 {
		return "", fmt.Errorf("没有文件需要打包")
	}

	// 获取当前工作目录
	currentDir, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("获取当前目录失败: %v", err)
	}

	// 确保临时目录存在（应用当前目录下的temp）
	tempDir := filepath.Join(currentDir, TempDirName)
	if err := os.MkdirAll(tempDir, 0755); err != nil {
		return "", fmt.Errorf("创建临时目录失败: %v", err)
	}

	// 生成文件名：mgg_表名_时间戳_随机4位.zip
	timestamp := time.Now().Format("20060102_150405")
	rand.Seed(time.Now().UnixNano()) // Initialize random seed
	random := generateRandomString(4)
	zipName := fmt.Sprintf("mgg_%s_%s_%s.zip", tableName, timestamp, random)
	zipPath := filepath.Join(tempDir, zipName)

	// 创建ZIP文件
	zipFile, err := os.Create(zipPath)
	if err != nil {
		return "", fmt.Errorf("创建ZIP文件失败: %v", err)
	}
	defer zipFile.Close()

	zipWriter := zip.NewWriter(zipFile)
	defer zipWriter.Close()

	// 添加每个文件到ZIP
	for _, file := range files {
		if err := addFileToZip(zipWriter, file, projectFolder); err != nil {
			return "", err
		}
	}

	return zipPath, nil
}

// generateRandomString 生成随机字符串
func generateRandomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyz0123456789"
	result := make([]byte, length)
	for i := range result {
		result[i] = charset[rand.Intn(len(charset))]
	}
	return string(result)
}

// addFileToZip 添加文件到ZIP归档
func addFileToZip(zipWriter *zip.Writer, filename string, basePath string) error {
	// 打开文件
	file, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("打开文件失败: %v", err)
	}
	defer file.Close()

	// 获取文件信息
	info, err := file.Stat()
	if err != nil {
		return fmt.Errorf("获取文件信息失败: %v", err)
	}

	// 计算ZIP中的相对路径
	relPath, err := filepath.Rel(basePath, filename)
	if err != nil {
		relPath = filepath.Base(filename)
	}

	// 统一使用斜杠作为路径分隔符(ZIP标准)
	relPath = strings.ReplaceAll(relPath, "\\", "/")

	// 创建ZIP文件头
	header, err := zip.FileInfoHeader(info)
	if err != nil {
		return fmt.Errorf("创建文件头失败: %v", err)
	}
	header.Name = relPath
	header.Method = zip.Deflate

	// 写入文件头
	writer, err := zipWriter.CreateHeader(header)
	if err != nil {
		return fmt.Errorf("写入文件头失败: %v", err)
	}

	// 复制文件内容
	_, err = io.Copy(writer, file)
	return err
}

// CleanExpiredZips 清理过期的ZIP文件
func CleanExpiredZips() {
	// 获取当前工作目录
	currentDir, err := os.Getwd()
	if err != nil {
		return
	}

	tempDir := filepath.Join(currentDir, TempDirName)

	// 检查目录是否存在
	if _, err := os.Stat(tempDir); os.IsNotExist(err) {
		return
	}

	// 读取目录中的所有文件
	entries, err := os.ReadDir(tempDir)
	if err != nil {
		return
	}

	now := time.Now()
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		// 只处理mgg_开头的zip文件
		if !strings.HasPrefix(entry.Name(), "mgg_") || !strings.HasSuffix(entry.Name(), ".zip") {
			continue
		}

		filePath := filepath.Join(tempDir, entry.Name())
		info, err := entry.Info()
		if err != nil {
			continue
		}

		// 如果文件超过过期时间，删除它
		if now.Sub(info.ModTime()) > FileExpireDuration {
			os.Remove(filePath)
		}
	}
}

// StartCleanupScheduler 启动定时清理任务
func StartCleanupScheduler() {
	ticker := time.NewTicker(5 * time.Minute)
	go func() {
		for range ticker.C {
			CleanExpiredZips()
		}
	}()
}
