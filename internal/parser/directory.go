package parser

import (
	"os"
	"path/filepath"

	"github.com/kihyun1998/aimd/pkg/utils"
)

// DirectoryParser 구현체
type directoryParser struct {
	excludeDirs   []string
	includeHidden bool // 숨김 파일 포함 여부 추가
}

func NewDirectoryParser(excludeDirs []string, includeHidden bool) DirectoryParser {
	return &directoryParser{
		excludeDirs:   excludeDirs,
		includeHidden: includeHidden,
	}
}

func (d *directoryParser) Parse(root string) ([]string, error) {
	var files []string

	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// 숨김 파일/디렉토리 처리
		if !d.includeHidden && utils.IsHidden(info.Name()) {
			if info.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}

		if info.IsDir() {
			if d.isExcluded(info.Name()) {
				return filepath.SkipDir
			}
			return nil
		}

		files = append(files, path)
		return nil
	})

	return files, err
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
