package ignore

import (
	"bufio"
	"os"
	"path/filepath"
	"strings"
)

// CodeIgnorePattern은 하나의 codeignore 패턴을 나타냅니다
type CodeIgnorePattern struct {
	pattern     string
	isNegative  bool
	isDirectory bool
}

// IsMatch는 패턴이 주어진 경로와 매칭되는지 확인합니다
func (p *CodeIgnorePattern) IsMatch(path string) bool {
	if path == "" {
		return false
	}

	pattern := p.pattern
	if !strings.HasPrefix(pattern, "/") {
		pattern = "**/" + pattern
	}
	pattern = strings.TrimPrefix(pattern, "/")

	if p.isDirectory {
		pattern += "/**"
	}

	matched, err := filepath.Match(pattern, path)
	if err != nil {
		return false
	}

	if strings.Contains(pattern, "**") {
		pathParts := strings.Split(path, "/")
		patternParts := strings.Split(pattern, "/")
		return matchWithDoublestar(patternParts, pathParts)
	}

	return matched
}

func (p *CodeIgnorePattern) IsNegative() bool {
	return p.isNegative
}

func (p *CodeIgnorePattern) IsDirectory() bool {
	return p.isDirectory
}

// CodeIgnore는 .codeignore 파일의 패턴들을 관리합니다
type CodeIgnore struct {
	patterns []Pattern
	root     string
}

// NewCodeIgnore는 주어진 경로의 .codeignore 파일을 읽어서 새로운 CodeIgnore 객체를 생성합니다
func NewCodeIgnore(path string) (Ignorer, error) {
	ci := &CodeIgnore{
		root: filepath.Dir(path),
	}
	if err := ci.LoadFromFile(path); err != nil {
		return nil, err
	}
	return ci, nil
}

// LoadFromFile은 .codeignore 파일을 읽어서 패턴을 로드합니다
func (ci *CodeIgnore) LoadFromFile(path string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		pattern := strings.TrimSpace(scanner.Text())
		if pattern == "" || strings.HasPrefix(pattern, "#") {
			continue
		}
		if err := ci.AddPattern(pattern); err != nil {
			return err
		}
	}
	return scanner.Err()
}

// AddPattern은 새로운 무시 패턴을 추가합니다
func (ci *CodeIgnore) AddPattern(pattern string) error {
	p := &CodeIgnorePattern{
		pattern: pattern,
	}

	if strings.HasPrefix(pattern, "!") {
		p.isNegative = true
		pattern = strings.TrimPrefix(pattern, "!")
		p.pattern = pattern
	}

	if strings.HasSuffix(pattern, "/") {
		p.isDirectory = true
		pattern = strings.TrimSuffix(pattern, "/")
		p.pattern = pattern
	}

	ci.patterns = append(ci.patterns, p)
	return nil
}

// ShouldIgnore는 주어진 경로가 무시되어야 하는지 확인합니다
func (ci *CodeIgnore) ShouldIgnore(path string) bool {
	relPath, err := filepath.Rel(ci.root, path)
	if err != nil {
		return false
	}

	relPath = filepath.ToSlash(relPath)

	var (
		isMatched bool
		lastMatch Pattern
	)

	for _, pattern := range ci.patterns {
		if pattern.IsMatch(relPath) {
			isMatched = true
			lastMatch = pattern
		}
	}

	if !isMatched || lastMatch == nil {
		return false
	}

	return !lastMatch.IsNegative()
}

// matchWithDoublestar는 ** 패턴을 포함한 매칭을 처리합니다
func matchWithDoublestar(patternParts, pathParts []string) bool {
	patternIdx := 0
	pathIdx := 0

	for patternIdx < len(patternParts) && pathIdx < len(pathParts) {
		if patternParts[patternIdx] == "**" {
			patternIdx++
			if patternIdx >= len(patternParts) {
				return true
			}
			for ; pathIdx < len(pathParts); pathIdx++ {
				if match, _ := filepath.Match(patternParts[patternIdx], pathParts[pathIdx]); match {
					if matchWithDoublestar(patternParts[patternIdx:], pathParts[pathIdx:]) {
						return true
					}
				}
			}
			return false
		}

		matched, _ := filepath.Match(patternParts[patternIdx], pathParts[pathIdx])
		if !matched {
			return false
		}
		patternIdx++
		pathIdx++
	}

	return patternIdx == len(patternParts) && pathIdx == len(pathParts)
}
