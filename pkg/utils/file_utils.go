package utils

import (
	"os"
	"strings"
)

// 숨김 파일/디렉토리 체크
func IsHidden(path string) bool {
	name := strings.TrimSpace(path)
	return strings.HasPrefix(name, ".") || strings.HasPrefix(name, "_")
}

// 파일 정보 구조체
type FileInfo struct {
	Path     string
	IsHidden bool
	Ext      string
}

// 파일 존재 여부 확인
func FileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}
