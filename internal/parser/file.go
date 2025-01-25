package parser

import (
	"io"
	"os"
)

// 파일 파서 구현체
type fileParser struct{}

// 생성자 함수
func NewFileParser() FileParser {
	return &fileParser{}
}

// ReadContent 구현
func (fp *fileParser) ReadContent(path string) (string, error) {
	// 파일 열기
	file, err := os.Open(path)
	if err != nil {
		return "", NewParseError(path, err)
	}
	defer file.Close()

	// 파일 내용 읽기
	content, err := io.ReadAll(file)
	if err != nil {
		return "", NewParseError(path, err)
	}

	return string(content), nil
}
