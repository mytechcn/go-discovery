// Package utils 文件操作工具
package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func WriteFile(body string) {
	ms := time.Now().UnixMilli()
	path := filepath.Join("tmp", fmt.Sprintf("%d.device.xml", ms))
	os.MkdirAll("tmp", 0o755)
	_ = os.WriteFile(path, []byte(body), 0o644)
}

func WriteJSONFile(filename, body string) {
	os.MkdirAll("tmpjson", 0o755)
	path := filepath.Join("tmpjson", filename)
	_ = os.WriteFile(path, []byte(body), 0o644)
}

// GetFileNameNoExt 获取文件名(剔除后缀)
func GetFileNameNoExt(path string) string {
	filename := filepath.Base(path)
	ext := filepath.Ext(filename)
	return strings.TrimSuffix(filename, ext)
}
