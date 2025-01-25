package parser

import "fmt"

type ParseError struct {
	Path string
	Err  error
}

func (e *ParseError) Error() string {
	return fmt.Sprintf("파싱 에러 (경로: %s): %v", e.Path, e.Err)
}

// 에러 생성 함수
func NewParseError(path string, err error) error {
	return &ParseError{
		Path: path,
		Err:  err,
	}
}
