package utils

import (
	"bufio"
	"os"
	"path/filepath"
	"strings"
)

type GitIgnore struct {
	patterns []string
}

func NewGitIgnore(path string) (*GitIgnore, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var patterns []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		pattern := strings.TrimSpace(scanner.Text())

		// 빈줄 또는 주석은 무시
		if pattern != "" && !strings.HasPrefix(pattern, "#") {
			patterns = append(patterns, pattern)
		}
	}

	return &GitIgnore{patterns: patterns}, nil
}

func (gi *GitIgnore) ShouldIgnore(path string) bool {
	for _, pattern := range gi.patterns {
		matched, err := filepath.Match(pattern, path)
		if err == nil && matched {
			return true
		}
	}
	return false
}
