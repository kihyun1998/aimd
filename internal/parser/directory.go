package parser

import (
	"os"
	"path/filepath"

	"github.com/kihyun1998/codemd/pkg/utils"
)

// DirectoryParser 구현체
type directoryParser struct {
	excludeDirs   []string
	includeHidden bool // 숨김 파일 포함 여부 추가
	gitIgnore     *utils.GitIgnore
}

func NewDirectoryParser(excludeDirs []string, includeHidden bool, useGitIgnore bool) DirectoryParser {
	var gitIgnore *utils.GitIgnore
	if useGitIgnore {
		if gi, err := utils.NewGitIgnore(".gitignore"); err == nil {
			gitIgnore = gi
		}
	}

	return &directoryParser{
		excludeDirs:   excludeDirs,
		includeHidden: includeHidden,
		gitIgnore:     gitIgnore,
	}
}

// 모든 파일 가져오기
func (d *directoryParser) Parse(root string) ([]string, error) {
	var files []string

	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return NewParseError(path, err)
		}

		// .gitIgnore 규칙 체크
		if d.gitIgnore != nil && d.gitIgnore.ShouldIgnore(path) {
			if info.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}

		// 숨김 파일 체크
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

	if err != nil {
		return nil, err
	}

	return files, err
}

// 특정 타입의 파일만 필터링 (마크다운 생성용)
func (d *directoryParser) GetFilesByTypes(allFiles []string, types []string) []string {
	if len(types) == 0 {
		return allFiles
	}

	var filtered []string
	for _, file := range allFiles {
		ext := filepath.Ext(file)
		if ext != "" {
			ext = ext[1:] // 점(.) 제거
			for _, t := range types {
				if ext == t {
					filtered = append(filtered, file)
					break
				}
			}
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
