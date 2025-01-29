package ignore

import (
	"os"
	"path/filepath"
	"testing"
)

// TestNewCodeIgnore는 CodeIgnore 인스턴스 생성을 테스트
func TestNewCodeIgnore(t *testing.T) {
	// 임시 디렉토리 생성
	tmpDir := t.TempDir()

	// 테스트용 .codeignore 파일 생성
	codeignorePath := filepath.Join(tmpDir, ".codeignore")
	content := `# 주석
*.log
/node_modules/
!important.log
build/
*.tmp
`
	err := os.WriteFile(codeignorePath, []byte(content), 0644)
	if err != nil {
		t.Fatalf(".codeignore 파일 생성 실패: %v", err)
	}

	// CodeIgnore 인스턴스 생성 테스트
	ignorer, err := NewCodeIgnore(codeignorePath)
	if err != nil {
		t.Fatalf("NewCodeIgnore 실패: %v", err)
	}

	// CodeIgnore로 타입 변환
	ci, ok := ignorer.(*CodeIgnore)
	if !ok {
		t.Fatal("Ignorer를 *CodeIgnore로 변환 실패")
	}

	// 패턴 수 확인
	expectedPatterns := 5 // 주석 제외
	if len(ci.patterns) != expectedPatterns {
		t.Errorf("패턴 수가 일치하지 않음. got %d, want %d", len(ci.patterns), expectedPatterns)
	}
}

// TestCodeIgnorePatterns는 다양한 codeignore 패턴 매칭을 테스트
func TestCodeIgnorePatterns(t *testing.T) {
	tests := []struct {
		name     string
		patterns []string
		paths    map[string]bool
	}{
		{
			name:     "기본 글로브 패턴",
			patterns: []string{"*.log"},
			paths: map[string]bool{
				"test.log":     true,
				"dir/test.log": true,
				"test.txt":     false,
				"log.txt":      false,
			},
		},
		{
			name:     "디렉토리 패턴",
			patterns: []string{"node_modules/"},
			paths: map[string]bool{
				"node_modules/file.js":     true,
				"node_modules/dir/file.js": true,
				"dir/node_modules/file.js": true,
				"nodemodules/file.js":      false,
			},
		},
		{
			name: "네거티브 패턴",
			patterns: []string{
				"*.log",          // 먼저 모든 로그 파일을 무시
				"!important.log", // important.log는 예외 처리
			},
			paths: map[string]bool{
				"important.log":     false, // 무시하지 않음
				"test.log":          true,  // 무시함
				"dir/important.log": false, // 하위 디렉토리도 무시하지 않음
				"logs/test.log":     true,  // 다른 로그 파일은 무시
			},
		},
		{
			name:     "중첩 디렉토리 패턴",
			patterns: []string{"**/temp/**"},
			paths: map[string]bool{
				"temp/file.txt":            true,
				"dir/temp/file.txt":        true,
				"dir/temp/subdir/file.txt": true,
				"template/file.txt":        false,
			},
		},
		{
			name: "복잡한 네거티브 패턴",
			patterns: []string{
				"*.log",
				"!important.log",
				"trace.*",
			},
			paths: map[string]bool{
				"debug.log":     true,  // *.log에 의해 무시
				"important.log": false, // !important.log에 의해 예외 처리
				"trace.log":     true,  // trace.*에 의해 다시 무시
				"trace.txt":     true,  // trace.*에 의해 무시
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 임시 디렉토리 생성
			tmpDir := t.TempDir()
			codeignorePath := filepath.Join(tmpDir, ".codeignore")

			// 패턴들을 .codeignore 파일에 작성
			content := ""
			for _, pattern := range tt.patterns {
				content += pattern + "\n"
			}
			err := os.WriteFile(codeignorePath, []byte(content), 0644)
			if err != nil {
				t.Fatalf(".codeignore 파일 생성 실패: %v", err)
			}

			ignorer, err := NewCodeIgnore(codeignorePath)
			if err != nil {
				t.Fatalf("NewCodeIgnore 실패: %v", err)
			}

			for path, shouldIgnore := range tt.paths {
				fullPath := filepath.Join(tmpDir, path)
				// 디렉토리 생성
				dir := filepath.Dir(fullPath)
				if err := os.MkdirAll(dir, 0755); err != nil {
					t.Fatalf("디렉토리 생성 실패: %v", err)
				}
				// 파일 생성
				if err := os.WriteFile(fullPath, []byte("test"), 0644); err != nil {
					t.Fatalf("테스트 파일 생성 실패: %v", err)
				}

				got := ignorer.ShouldIgnore(fullPath)
				if got != shouldIgnore {
					t.Errorf("패턴 %q에 대해 경로 %q의 결과가 잘못됨. got %v, want %v",
						tt.patterns, path, got, shouldIgnore)
				}
			}
		})
	}
}

// TestCodeIgnoreEdgeCases는 특수한 경우의 codeignore 패턴을 테스트
func TestCodeIgnoreEdgeCases(t *testing.T) {
	tests := []struct {
		name     string
		patterns []string
		path     string
		want     bool
	}{
		{
			name:     "빈 패턴",
			patterns: []string{""},
			path:     "test.txt",
			want:     false,
		},
		{
			name:     "주석만 있는 경우",
			patterns: []string{"# 주석입니다"},
			path:     "test.txt",
			want:     false,
		},
		{
			name: "복잡한 패턴 조합",
			patterns: []string{
				"*.log",
				"!important.log",
				"logs/",
			},
			path: "logs/important.log",
			want: true, // 디렉토리 패턴이 네거티브 패턴보다 우선
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpDir := t.TempDir()
			codeignorePath := filepath.Join(tmpDir, ".codeignore")
			content := ""
			for _, pattern := range tt.patterns {
				content += pattern + "\n"
			}
			err := os.WriteFile(codeignorePath, []byte(content), 0644)
			if err != nil {
				t.Fatalf(".codeignore 파일 생성 실패: %v", err)
			}

			ignorer, err := NewCodeIgnore(codeignorePath)
			if err != nil {
				t.Fatalf("NewCodeIgnore 실패: %v", err)
			}

			fullPath := filepath.Join(tmpDir, tt.path)
			got := ignorer.ShouldIgnore(fullPath)
			if got != tt.want {
				t.Errorf("%q 경로에 대한 결과가 잘못됨. got %v, want %v",
					tt.path, got, tt.want)
			}
		})
	}
}

// TestCodeIgnoreWithDirectoryParser는 DirectoryParser와의 통합을 테스트
func TestCodeIgnoreWithDirectoryParser(t *testing.T) {
	// 임시 디렉토리 구조 생성
	tmpDir := t.TempDir()
	files := map[string]bool{ // 경로와 무시 여부
		"test.go":               false,
		"test.log":              true,
		"build/output.go":       true,
		"src/main.go":           false,
		"node_modules/index.js": true,
	}

	// .codeignore 파일 생성
	codeignoreContent := `*.log
build/
node_modules/
`
	err := os.WriteFile(filepath.Join(tmpDir, ".codeignore"), []byte(codeignoreContent), 0644)
	if err != nil {
		t.Fatalf(".codeignore 파일 생성 실패: %v", err)
	}

	// 테스트 파일 생성
	for path := range files {
		fullPath := filepath.Join(tmpDir, path)
		dir := filepath.Dir(fullPath)
		if err := os.MkdirAll(dir, 0755); err != nil {
			t.Fatalf("디렉토리 생성 실패: %v", err)
		}
		if err := os.WriteFile(fullPath, []byte("test"), 0644); err != nil {
			t.Fatalf("파일 생성 실패: %v", err)
		}
	}

	// CodeIgnore 인스턴스 생성
	ignorer, err := NewCodeIgnore(filepath.Join(tmpDir, ".codeignore"))
	if err != nil {
		t.Fatalf("NewCodeIgnore 실패: %v", err)
	}

	// 각 파일에 대해 CodeIgnore 규칙 테스트
	for path, shouldIgnore := range files {
		fullPath := filepath.Join(tmpDir, path)
		got := ignorer.ShouldIgnore(fullPath)
		if got != shouldIgnore {
			t.Errorf("경로 %q에 대한 결과가 잘못됨. got %v, want %v",
				path, got, shouldIgnore)
		}
	}
}

func TestFlutterPluginPattern(t *testing.T) {
	// 임시 디렉토리 생성
	tmpDir := t.TempDir()

	// 테스트용 디렉토리 구조 생성
	testFiles := map[string]bool{
		"lib/plugin.dart":             false, // 추적해야 함
		"lib/src/implementation.dart": false, // 추적해야 함
		"windows/plugin.cpp":          false, // 추적해야 함
		"windows/include/header.h":    false, // 추적해야 함
		"android/build.gradle":        true,  // 무시해야 함
		"ios/plugin.podspec":          true,  // 무시해야 함
		"linux/CMakeLists.txt":        true,  // 무시해야 함
		"macos/plugin.swift":          true,  // 무시해야 함
		"web/plugin.js":               true,  // 무시해야 함
		"pubspec.yaml":                true,  // 무시해야 함
		"README.md":                   true,  // 무시해야 함
	}

	// .codeignore 파일 생성
	codeignoreContent := `# 모든 파일을 무시
*
# lib 디렉토리와 windows 디렉토리만 추적
!lib/
!windows/
# 기타 플랫폼 디렉토리는 무시
android/
ios/
linux/
macos/
web/
`
	codeignorePath := filepath.Join(tmpDir, ".codeignore")
	err := os.WriteFile(codeignorePath, []byte(codeignoreContent), 0644)
	if err != nil {
		t.Fatalf(".codeignore 파일 생성 실패: %v", err)
	}

	// 테스트 파일 구조 생성
	for path := range testFiles {
		fullPath := filepath.Join(tmpDir, path)
		dir := filepath.Dir(fullPath)
		if err := os.MkdirAll(dir, 0755); err != nil {
			t.Fatalf("디렉토리 생성 실패: %v", err)
		}
		if err := os.WriteFile(fullPath, []byte("test"), 0644); err != nil {
			t.Fatalf("파일 생성 실패: %v", err)
		}
	}

	// CodeIgnore 인스턴스 생성
	ignorer, err := NewCodeIgnore(codeignorePath)
	if err != nil {
		t.Fatalf("NewCodeIgnore 실패: %v", err)
	}

	// 각 파일에 대해 무시 규칙 테스트
	for path, shouldIgnore := range testFiles {
		fullPath := filepath.Join(tmpDir, path)
		got := ignorer.ShouldIgnore(fullPath)
		if got != shouldIgnore {
			t.Errorf("경로 %q에 대한 결과가 잘못됨. got %v, want %v",
				path, got, shouldIgnore)
		}
	}
}

// 중첩된 디렉토리 패턴 테스트
func TestNestedDirectoryPatterns(t *testing.T) {
	tmpDir := t.TempDir()

	testFiles := map[string]bool{
		"lib/plugin.dart":                  false, // 추적
		"lib/generated/temp.dart":          true,  // 무시
		"windows/plugin.cpp":               false, // 추적
		"windows/build/temp.cpp":           true,  // 무시
		"windows/include/header.h":         false, // 추적
		"windows/include/generated/temp.h": true,  // 무시
	}

	// .codeignore 파일 생성
	codeignoreContent := `# 기본 추적 제외
*
# 필요한 디렉토리 추적
!lib/
!windows/
# 생성된 파일들 무시
**/generated/**
**/build/**
`
	codeignorePath := filepath.Join(tmpDir, ".codeignore")
	err := os.WriteFile(codeignorePath, []byte(codeignoreContent), 0644)
	if err != nil {
		t.Fatalf(".codeignore 파일 생성 실패: %v", err)
	}

	// 테스트 파일 구조 생성
	for path := range testFiles {
		fullPath := filepath.Join(tmpDir, path)
		dir := filepath.Dir(fullPath)
		if err := os.MkdirAll(dir, 0755); err != nil {
			t.Fatalf("디렉토리 생성 실패: %v", err)
		}
		if err := os.WriteFile(fullPath, []byte("test"), 0644); err != nil {
			t.Fatalf("파일 생성 실패: %v", err)
		}
	}

	// CodeIgnore 인스턴스 생성
	ignorer, err := NewCodeIgnore(codeignorePath)
	if err != nil {
		t.Fatalf("NewCodeIgnore 실패: %v", err)
	}

	// 각 파일에 대해 무시 규칙 테스트
	for path, shouldIgnore := range testFiles {
		fullPath := filepath.Join(tmpDir, path)
		got := ignorer.ShouldIgnore(fullPath)
		if got != shouldIgnore {
			t.Errorf("경로 %q에 대한 결과가 잘못됨. got %v, want %v",
				path, got, shouldIgnore)
		}
	}
}
