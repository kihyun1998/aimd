package test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/kihyun1998/codemd/internal/file"
)

func TestFileSplitter(t *testing.T) {
	// 테스트 디렉토리 생성
	tempDir := t.TempDir()
	testFile := filepath.Join(tempDir, "CODE.md")

	// 1MB의 테스트 데이터 생성
	testData := make([]byte, 1024*1024*100)
	for i := range testData {
		testData[i] = 'A'
	}
	content := string(testData)

	tests := []struct {
		name        string
		maxSizeMB   int64
		content     string
		wantFiles   int
		shouldSplit bool
	}{
		{
			name:        "파일 분할 불필요",
			maxSizeMB:   10,
			content:     content,
			wantFiles:   10,
			shouldSplit: false,
		},
		{
			name:        "파일 분할 필요",
			maxSizeMB:   1, // 1MB
			content:     content,
			wantFiles:   100,
			shouldSplit: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			splitter := file.NewFileSplitter(tt.maxSizeMB)
			err := splitter.SplitIfNeeded(tt.content, testFile)

			if err != nil {
				t.Errorf("SplitIfNeeded() error = %v", err)
				return
			}

			// 파일 개수 확인
			files, err := os.ReadDir(tempDir)
			if err != nil {
				t.Fatal(err)
			}

			fileCount := 0
			for _, f := range files {
				if !f.IsDir() {
					fileCount++
				}
			}

			if fileCount != tt.wantFiles {
				t.Errorf("파일 개수가 일치하지 않음. got %v, want %v", fileCount, tt.wantFiles)
			}
		})
	}
}
