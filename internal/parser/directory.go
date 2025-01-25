package parser

import (
	"os"
	"path/filepath"
)

// DirectoryParser 구현체
type directoryParser struct {
	excludeDirs []string // 제외할 디렉토리 목록
}

// 생성자 함수
func NewDirectoryParser(excludeDirs []string) DirectoryParser {
	return &directoryParser{
		excludeDirs: excludeDirs,
	}
}

// Parse 구현 - 재귀적 디렉토리 탐색
func (d *directoryParser) Parse(root string) ([]string, error) {
	var files []string

	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// 디렉토리면 제외 디렉토리인지 확인
		if info.IsDir() {
			if d.isExcluded(info.Name()) {
				return filepath.SkipDir
			}
			return nil
		}

		// 파일이면 목록에 추가
		files = append(files, path)
		return nil
	})

	if err != nil {
		return nil, err
	}

	return files, nil
}

// 파일 확장자 필터링
func (d *directoryParser) FilterByExtenstion(files []string, ext string) []string {
	var filtered []string
	for _, file := range files {
		if filepath.Ext(file) == "."+ext {
			filtered = append(filtered, file)
		}
	}
	return filtered
}

// 제외 디렉토리 확인
func (d *directoryParser) isExcluded(dir string) bool {
	for _, excludeDir := range d.excludeDirs {
		if dir == excludeDir {
			return true
		}
	}
	return false
}
