package generator

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// CreateZipArchive 创建ZIP归档文件
func CreateZipArchive(files []string, projectFolder string) (string, error) {
	if len(files) == 0 {
		return "", fmt.Errorf("没有文件需要打包")
	}

	// 创建临时ZIP文件
	zipPath := filepath.Join(os.TempDir(), fmt.Sprintf("mybatis-gen-%d.zip", os.Getpid()))

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
