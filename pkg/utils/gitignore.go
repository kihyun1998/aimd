package utils

import (
	"bufio"
	"os"
	"path/filepath"
	"strings"
)

// GitIgnorePattern은 하나의 gitignore 패턴을 나타냄
type GitIgnorePattern struct {
	pattern     string // 원본 패턴
	isNegative  bool   // ! 패턴 여부
	isDirectory bool   // 디렉토리 전용 패턴 여부
}

// GitIgnore는 .gitignore 파일의 패턴들을 관리
type GitIgnore struct {
	patterns []*GitIgnorePattern
	root     string // .gitignore 파일이 위치한 루트 경로
}

// NewGitIgnore는 주어진 경로의 .gitignore 파일을 읽어서 새로운 GitIgnore 객체를 생성
func NewGitIgnore(path string) (*GitIgnore, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	root := filepath.Dir(path)
	gi := &GitIgnore{
		root: root,
	}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		pattern := strings.TrimSpace(scanner.Text())

		// 빈 줄이나 주석은 무시
		if pattern == "" || strings.HasPrefix(pattern, "#") {
			continue
		}

		if err := gi.AddPattern(pattern); err != nil {
			return nil, err
		}
	}

	return gi, scanner.Err()
}

// AddPattern은 새로운 gitignore 패턴을 추가
func (gi *GitIgnore) AddPattern(pattern string) error {
	p := &GitIgnorePattern{
		pattern: pattern,
	}

	// 네거티브 패턴 처리
	if strings.HasPrefix(pattern, "!") {
		p.isNegative = true
		pattern = strings.TrimPrefix(pattern, "!")
		p.pattern = pattern
	}

	// 디렉토리 전용 패턴 처리
	if strings.HasSuffix(pattern, "/") {
		p.isDirectory = true
		pattern = strings.TrimSuffix(pattern, "/")
		p.pattern = pattern
	}

	gi.patterns = append(gi.patterns, p)
	return nil
}

// ShouldIgnore는 주어진 경로가 무시되어야 하는지 확인
func (gi *GitIgnore) ShouldIgnore(path string) bool {
	// 상대 경로로 변환
	relPath, err := filepath.Rel(gi.root, path)
	if err != nil {
		return false
	}

	// 윈도우 경로를 유닉스 스타일로 변환
	relPath = filepath.ToSlash(relPath)

	// 마지막 매칭 결과 (기본값은 false)
	ignored := false

	// 모든 패턴을 순회하면서 검사
	for _, pattern := range gi.patterns {
		if pattern.matches(relPath) {
			// 네거티브 패턴은 이전 매칭을 뒤집음
			if pattern.isNegative {
				ignored = false
			} else {
				ignored = true
			}
		}
	}

	return ignored
}

// matches는 패턴이 주어진 경로와 매칭되는지 확인
func (p *GitIgnorePattern) matches(path string) bool {
	if path == "" {
		return false
	}

	// 패턴을 glob 패턴으로 변환
	pattern := p.pattern

	// 패턴이 슬래시로 시작하지 않으면 "**/" 접두사 추가
	if !strings.HasPrefix(pattern, "/") {
		pattern = "**/" + pattern
	}

	// 시작 슬래시 제거
	pattern = strings.TrimPrefix(pattern, "/")

	// 디렉토리 매칭을 위한 처리
	if p.isDirectory {
		pattern += "/**"
	}

	// glob 패턴 매칭
	matched, err := filepath.Match(pattern, path)
	if err != nil {
		// 유효하지 않은 패턴은 false 반환
		return false
	}

	// ** 패턴 처리
	if strings.Contains(pattern, "**") {
		pathParts := strings.Split(path, "/")
		patternParts := strings.Split(pattern, "/")
		return matchWithDoublestar(patternParts, pathParts)
	}

	return matched
}

// matchWithDoublestar는 ** 패턴을 포함한 매칭을 처리
func matchWithDoublestar(patternParts, pathParts []string) bool {
	patternIdx := 0
	pathIdx := 0

	for patternIdx < len(patternParts) && pathIdx < len(pathParts) {
		if patternParts[patternIdx] == "**" {
			patternIdx++
			// ** 다음에 더 이상 패턴이 없으면 나머지 모두 매칭
			if patternIdx >= len(patternParts) {
				return true
			}
			// 다음 세그먼트가 매칭될 때까지 경로 진행
			for ; pathIdx < len(pathParts); pathIdx++ {
				if match, _ := filepath.Match(patternParts[patternIdx], pathParts[pathIdx]); match {
					if matchWithDoublestar(patternParts[patternIdx:], pathParts[pathIdx:]) {
						return true
					}
				}
			}
			return false
		}

		// 일반 글로브 매칭
		matched, _ := filepath.Match(patternParts[patternIdx], pathParts[pathIdx])
		if !matched {
			return false
		}
		patternIdx++
		pathIdx++
	}

	// 남은 부분이 없어야 매칭 성공
	return patternIdx == len(patternParts) && pathIdx == len(pathParts)
}
