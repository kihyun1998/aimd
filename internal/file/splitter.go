package file

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// FileSplitter는 파일 분할 인터페이스
type FileSplitter interface {
	SplitIfNeeded(content string, basePath string) error
}

// 파일 분할을 위한 구조체
type fileSplitter struct {
	maxFileSize int64 // 최대 파일 크기 (바이트)
}

// NewFileSplitter는 FileSplitter 인스턴스를 생성
func NewFileSplitter(maxSizeMB int64) FileSplitter {
	return &fileSplitter{
		maxFileSize: maxSizeMB * 1024 * 1024, // MB를 바이트로 변환
	}
}

// SplitIfNeeded는 콘텐츠를 여러 파일로 분할
func (fs *fileSplitter) SplitIfNeeded(content string, basePath string) error {
	contentSize := int64(len(content))

	// 최대 크기를 초과하지 않으면 단일 파일로 저장
	if contentSize <= fs.maxFileSize {
		return os.WriteFile(basePath, []byte(content), 0644)
	}

	// 파일 분할이 필요한 경우
	parts := fs.splitContent(content)

	// 각 부분을 개별 파일로 저장
	for i, part := range parts {
		fileName := fs.generateFileName(basePath, i+1)
		if err := os.WriteFile(fileName, []byte(part), 0644); err != nil {
			return fmt.Errorf("파일 분할 저장 실패: %w", err)
		}
	}

	return nil
}

// splitContent는 콘텐츠를 여러 부분으로 분할
func (fs *fileSplitter) splitContent(content string) []string {
	var parts []string
	contentSize := int64(len(content))
	numParts := (contentSize + fs.maxFileSize - 1) / fs.maxFileSize

	for i := int64(0); i < numParts; i++ {
		start := i * fs.maxFileSize
		end := (i + 1) * fs.maxFileSize
		if end > contentSize {
			end = contentSize
		}
		parts = append(parts, content[start:end])
	}

	return parts
}

// generateFileName은 분할된 파일의 이름을 생성
func (fs *fileSplitter) generateFileName(basePath string, index int) string {
	ext := filepath.Ext(basePath)
	baseWithoutExt := strings.TrimSuffix(basePath, ext)
	return fmt.Sprintf("%s%d%s", baseWithoutExt, index, ext)
}
