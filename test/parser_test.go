package test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/kihyun1998/aimd/internal/parser"
)

func TestDirectoryParser(t *testing.T) {
	// 테스트 디렉토리 생성
	tempDir := t.TempDir()

	// 테스트 파일 생성
	files := []string{
		"test1.go",
		"test2.go",
		"test3.txt",
		".hidden",
	}

	for _, f := range files {
		path := filepath.Join(tempDir, f)
		if err := os.WriteFile(path, []byte("test"), 0644); err != nil {
			t.Fatal(err)
		}
	}

	// 테스트 케이스
	tests := []struct {
		name          string
		excludeDirs   []string
		includeHidden bool
		wantCount     int
	}{
		{
			name:          "기본 테스트",
			excludeDirs:   nil,
			includeHidden: false,
			wantCount:     3,
		},
		{
			name:          "숨김 파일 포함",
			excludeDirs:   nil,
			includeHidden: true,
			wantCount:     4,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := parser.NewDirectoryParser(tt.excludeDirs, tt.includeHidden)

			got, err := p.Parse(tempDir)
			if err != nil {
				t.Errorf("Parse() error = %v", err)
				return
			}

			if len(got) != tt.wantCount {
				t.Errorf("Parse() = %v, want %v", len(got), tt.wantCount)
			}
		})
	}
}

func TestGetFilesByTypes(t *testing.T) {
	p := parser.NewDirectoryParser(nil, false)
	files := []string{
		"test1.go",
		"test2.go",
		"test3.txt",
	}

	filtered := p.GetFilesByTypes(files, []string{"go"})
	if len(filtered) != 2 {
		t.Errorf("GetFilesByTypes() = %v, want 2", len(filtered))
	}
}
