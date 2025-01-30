# codemd
## cmd/codemd/main.go
```go
package main

import (
	"fmt"
	"log"
	"os"

	"github.com/kihyun1998/codemd/internal/config"
	"github.com/kihyun1998/codemd/internal/generator"
	"github.com/kihyun1998/codemd/internal/parser"
	"github.com/kihyun1998/codemd/internal/version"
)

func main() {
	// 커스텀 usage 메시지 설정
	config.SetUsage(os.Args[0])

	// 설정 파싱
	cfg, err := config.ParseFlags()
	if err != nil {
		log.Fatal(err)
	}

	// 버전 출력
	if cfg.ShowVersion {
		fmt.Println(version.GetVersionInfo())
		os.Exit(0)
	}

	// 파서 생성
	dirParser := parser.NewDirectoryParser(cfg.ExcludeDirs, false, cfg.UseCodeIgnore)
	fileParser := parser.NewFileParser()

	// 현재 디렉토리 경로
	currentDir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	// 파일 목록 가져오기
	allFiles, err := dirParser.Parse(currentDir)
	if err != nil {
		log.Fatal(err)
	}

	// 타입별 필터링
	typeFiles := dirParser.GetFilesByTypes(allFiles, cfg.FileTypes)

	// 마크다운 생성기 생성
	mdGen := generator.NewMarkdownGenerator(fileParser, cfg.OutputPath)

	// 기본 템플릿 설정
	defaultTemplate := "# {{.ProjectName}}\n{{range .Files}}## {{.Path}}\n```{{if .Extension}}{{.Extension}}{{end}}\n{{.Content}}\n```\n{{end}}"
	if err := mdGen.SetTemplate(defaultTemplate); err != nil {
		log.Fatal(err)
	}
	// 마크다운 생성
	if err := mdGen.Generate(typeFiles); err != nil {
		log.Fatal(err)
	}
}

```
## internal/config/config.go
```go
package config

import (
	"flag"
	"fmt"
	"os"
	"strings"
)

type Config struct {
	FileTypes     []string
	OutputPath    string
	ExcludeDirs   []string
	UseCodeIgnore bool
	ShowVersion   bool
}

// Usage 메시지 설정
func SetUsage(programName string) {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "사용법: %s [옵션]\n\n옵션:\n", programName)
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\n예시:\n")
		fmt.Fprintf(os.Stderr, "  %s -version\n", programName)
		fmt.Fprintf(os.Stderr, "  %s -type go,java\n", programName)
		fmt.Fprintf(os.Stderr, "  %s -type go -exclude vendor,node_modules\n", programName)
	}
}

// 플래그 파싱
func ParseFlags() (*Config, error) {
	var (
		types         string
		output        string
		exclude       string
		useCodeIgnore bool
		showVersion   bool
	)

	// 긴 버전과 짧은 버전의 플래그 모두 지원
	flag.StringVar(&types, "type", "", "파일 확장자들 (쉼표로 구분)")
	flag.StringVar(&types, "t", "", "파일 확장자들 (쉼표로 구분) (짧은 버전)")

	flag.StringVar(&output, "out", "CODE.md", "출력 파일 경로")
	flag.StringVar(&output, "o", "CODE.md", "출력 파일 경로 (짧은 버전)")

	flag.StringVar(&exclude, "exclude", "", "제외할 디렉토리들 (쉼표로 구분)")
	flag.StringVar(&exclude, "e", "", "제외할 디렉토리들 (쉼표로 구분) (짧은 버전)")

	flag.BoolVar(&useCodeIgnore, "codeignore", false, ".codeignore 파일 사용 여부")
	flag.BoolVar(&useCodeIgnore, "c", false, ".codeignore 파일 사용 여부 (짧은 버전)")

	flag.BoolVar(&showVersion, "version", false, "버전 정보 출력")
	flag.BoolVar(&showVersion, "v", false, "버전 정보 출력 (짧은 버전)")

	flag.Parse()

	// -v 또는 -version 플래그만 있는 경우
	if showVersion && len(os.Args) == 2 {
		return &Config{ShowVersion: true}, nil
	}

	// -t 또는 -type 플래그가 없으면 에러
	if types == "" {
		return nil, fmt.Errorf("필수 플래그 누락: -type 또는 -t")
	}

	return &Config{
		FileTypes:     strings.Split(types, ","),
		OutputPath:    output,
		ExcludeDirs:   strings.Split(exclude, ","),
		UseCodeIgnore: useCodeIgnore,
		ShowVersion:   showVersion,
	}, nil
}

```
## internal/generator/markdown.go
```go
package generator

import (
	"os"
	"path/filepath"

	"github.com/kihyun1998/codemd/internal/parser"
)

type MarkdownGenerator interface {
	Generate(files []string) error
	SetTemplate(template string) error
}

// 마크다운 생성기 구조체
type markdownGenerator struct {
	fileParser  parser.FileParser
	processor   *templateProcessor
	outputPath  string
	rootDir     string
	projectName string
}

// 생성자
func NewMarkdownGenerator(fp parser.FileParser, outputPath string) MarkdownGenerator {
	rootDir, err := os.Getwd()
	if err != nil {
		rootDir = ""
	}

	// 프로젝트 이름 추출
	projectName := filepath.Base(rootDir)

	return &markdownGenerator{
		fileParser:  fp,
		outputPath:  outputPath,
		rootDir:     rootDir,
		projectName: projectName,
	}
}

// 상대 경로 변환 함수
func (mg *markdownGenerator) toRelativePath(absolutePath string) string {
	if mg.rootDir == "" {
		return absolutePath
	}

	relativePath, err := filepath.Rel(mg.rootDir, absolutePath)
	if err != nil {
		return absolutePath
	}

	// 윈도우 스타일 경로를 UNIX 스타일 경로로 변환
	return filepath.ToSlash(relativePath)
}

// 템플릿 설정
func (mg *markdownGenerator) SetTemplate(template string) error {
	processor, err := NewTemplateProcessor(template)
	if err != nil {
		return err
	}
	mg.processor = processor
	return nil
}

// 마크다운 생성
func (mg *markdownGenerator) Generate(files []string) error {
	var fileDataList []FileData

	for _, file := range files {
		content, err := mg.fileParser.ReadContent(file)
		if err != nil {
			return err
		}

		ext := filepath.Ext(file)
		if ext != "" {
			ext = ext[1:]
		}

		// 상대 경로로 변환
		relativePath := mg.toRelativePath(file)

		fileDataList = append(fileDataList, FileData{
			Path:      relativePath,
			Content:   content,
			Extension: ext,
		})
	}

	data := TemplateData{
		ProjectName: mg.projectName,
		Files:       fileDataList,
	}

	result, err := mg.processor.Execute(data)
	if err != nil {
		return err
	}

	return os.WriteFile(mg.outputPath, []byte(result), 0644)
}

```
## internal/generator/template.go
```go
package generator

import (
	"bytes"
	"path/filepath"
	"text/template"
)

// 템플릿 처리기 구조체
type templateProcessor struct {
	tmpl *template.Template
}

// 템플릿 데이터 구조체
type FileData struct {
	Path      string
	Content   string
	Extension string
}

type TemplateData struct {
	ProjectName string
	Files       []FileData
}

// 생성자 함수
func NewTemplateProcessor(templateStr string) (*templateProcessor, error) {
	tmpl, err := template.New("markdown").Parse(templateStr)
	if err != nil {
		return nil, err
	}
	return &templateProcessor{tmpl: tmpl}, nil
}

// 템플릿 실행
func (tp *templateProcessor) Execute(data TemplateData) (string, error) {
	var buf bytes.Buffer
	if err := tp.tmpl.Execute(&buf, data); err != nil {
		return "", err
	}
	return buf.String(), nil
}

func (tp *templateProcessor) getExtension(path string) string {
	ext := filepath.Ext(path)
	if ext != "" {
		return ext[1:]
	}
	return ""
}

```
## internal/ignore/codeignore.go
```go
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

```
## internal/ignore/codeignore_test.go
```go
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

```
## internal/ignore/ignore.go
```go
package ignore

// Pattern은 무시 패턴을 나타내는 인터페이스입니다
type Pattern interface {
	IsMatch(path string) bool
	IsNegative() bool
	IsDirectory() bool
}

// Ignorer는 파일 무시 규칙을 처리하는 인터페이스입니다
type Ignorer interface {
	AddPattern(pattern string) error
	ShouldIgnore(path string) bool
	LoadFromFile(path string) error
}

```
## internal/parser/directory.go
```go
package parser

import (
	"os"
	"path/filepath"

	"github.com/kihyun1998/codemd/internal/ignore"
	"github.com/kihyun1998/codemd/pkg/utils"
)

// DirectoryParser 구현체
type directoryParser struct {
	excludeDirs   []string
	includeHidden bool // 숨김 파일 포함 여부 추가
	ignorer       ignore.Ignorer
	rootDir       string
}

func NewDirectoryParser(excludeDirs []string, includeHidden bool, useCodeIgnore bool) DirectoryParser {
	var ignorer ignore.Ignorer

	rootDir, err := os.Getwd()
	if err != nil {
		rootDir = "."
	}

	if useCodeIgnore {
		codeIgnorePath := filepath.Join(rootDir, ".codeignore")
		if ci, err := ignore.NewCodeIgnore(codeIgnorePath); err == nil {
			ignorer = ci
		}
	}

	return &directoryParser{
		excludeDirs:   excludeDirs,
		includeHidden: includeHidden,
		ignorer:       ignorer,
		rootDir:       rootDir,
	}
}

// 모든 파일 가져오기
func (d *directoryParser) Parse(root string) ([]string, error) {
	var files []string

	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return NewParseError(path, err)
		}

		// .codeignore 규칙 체크
		if d.ignorer != nil {
			if d.ignorer.ShouldIgnore(path) {
				if info.IsDir() {
					return filepath.SkipDir
				}
				return nil
			}
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

	return files, nil
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

```
## internal/parser/error.go
```go
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

```
## internal/parser/file.go
```go
package parser

import (
	"io"
	"os"
)

// 파일 파서 구현체
type fileParser struct{}

// 생성자 함수
func NewFileParser() FileParser {
	return &fileParser{}
}

// ReadContent 구현
func (fp *fileParser) ReadContent(path string) (string, error) {
	// 파일 열기
	file, err := os.Open(path)
	if err != nil {
		return "", NewParseError(path, err)
	}
	defer file.Close()

	// 파일 내용 읽기
	content, err := io.ReadAll(file)
	if err != nil {
		return "", NewParseError(path, err)
	}

	return string(content), nil
}

```
## internal/parser/parser.go
```go
package parser

type DirectoryParser interface {
	Parse(root string) ([]string, error)
	// FilterByExtenstion(files []string, ext string) []string
	GetFilesByTypes(allFiles []string, types []string) []string
}

type FileParser interface {
	ReadContent(path string) (string, error)
}

```
## internal/version/version.go
```go
package version

// Version 정보
var (
	// 메이저 버전
	Major = "1"
	// 마이너 버전
	Minor = "1"
	// 패치 버전
	Patch = "0"
)

// GetVersion은 현재 버전 문자열을 반환합니다
func GetVersion() string {
	return Major + "." + Minor + "." + Patch
}

// GetVersionInfo는 버전 정보를 상세히 반환합니다
func GetVersionInfo() string {
	return "codemd version " + GetVersion()
}

```
## pkg/utils/file_utils.go
```go
package utils

import (
	"os"
	"strings"
)

// 숨김 파일/디렉토리 체크
func IsHidden(path string) bool {
	name := strings.TrimSpace(path)
	return strings.HasPrefix(name, ".") || strings.HasPrefix(name, "_")
}

// 파일 정보 구조체
type FileInfo struct {
	Path     string
	IsHidden bool
	Ext      string
}

// 파일 존재 여부 확인
func FileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

```
## test/generate_test.go
```go
package test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/kihyun1998/codemd/internal/generator"
	"github.com/kihyun1998/codemd/internal/parser"
)

func TestTemplateProcessor(t *testing.T) {
	tests := []struct {
		name     string
		template string
		data     generator.TemplateData
		want     string
		wantErr  bool
	}{
		{
			name:     "기본 템플릿 테스트",
			template: "# Files\n{{range .Files}}## {{.Path}}\n```\n{{.Content}}\n```\n{{end}}",
			data: generator.TemplateData{
				Files: []generator.FileData{
					{
						Path:    "test.go",
						Content: "package main",
					},
				},
			},
			want:    "# Files\n## test.go\n```\npackage main\n```\n",
			wantErr: false,
		},
		{
			name:     "잘못된 템플릿 문법",
			template: "{{.InvalidSyntax}",
			data:     generator.TemplateData{},
			want:     "",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tp, err := generator.NewTemplateProcessor(tt.template)
			if err != nil && !tt.wantErr {
				t.Errorf("NewTemplateProcessor() error = %v", err)
				return
			}

			if tp != nil {
				got, err := tp.Execute(tt.data)
				if (err != nil) != tt.wantErr {
					t.Errorf("Execute() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
				if got != tt.want {
					t.Errorf("Execute() = %v, want %v", got, tt.want)
				}
			}
		})
	}
}

func TestMarkdownGenerator(t *testing.T) {
	// 임시 디렉토리 및 파일 설정
	tempDir := t.TempDir()
	testFile := filepath.Join(tempDir, "test.go")
	outputPath := filepath.Join(tempDir, "CODE.md")
	testContent := "package main\n\nfunc main() {}"

	if err := os.WriteFile(testFile, []byte(testContent), 0644); err != nil {
		t.Fatal(err)
	}

	// 테스트용 파서 생성
	fp := parser.NewFileParser()

	// 테스트용 템플릿
	testTemplate := "# Files\n{{range .Files}}## {{.Path}}\n```go\n{{.Content}}\n```\n{{end}}"

	tests := []struct {
		name     string
		files    []string
		template string
		wantErr  bool
	}{
		{
			name:     "정상 마크다운 생성",
			files:    []string{testFile},
			template: testTemplate,
			wantErr:  false,
		},
		{
			name:     "존재하지 않는 파일",
			files:    []string{"없는파일.go"},
			template: testTemplate,
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mg := generator.NewMarkdownGenerator(fp, outputPath)

			err := mg.SetTemplate(tt.template)
			if err != nil {
				t.Errorf("SetTemplate() error = %v", err)
				return
			}

			err = mg.Generate(tt.files)
			if (err != nil) != tt.wantErr {
				t.Errorf("Generate() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !tt.wantErr {
				// 생성된 파일 확인
				_, err := os.Stat(outputPath)
				if err != nil {
					t.Errorf("Generated file not found: %v", err)
				}
			}
		})
	}
}

```
## test/parser_test.go
```go
package test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/kihyun1998/codemd/internal/parser"
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
			p := parser.NewDirectoryParser(tt.excludeDirs, tt.includeHidden, true)

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
	p := parser.NewDirectoryParser(nil, false, true)
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

func TestFileParser(t *testing.T) {
	// 임시 디렉토리 생성
	tempDir := t.TempDir()

	// 테스트 파일 생성
	testContent := "테스트 컨텐츠"
	testPath := filepath.Join(tempDir, "test.txt")
	if err := os.WriteFile(testPath, []byte(testContent), 0644); err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name    string
		path    string
		want    string
		wantErr bool
	}{
		{
			name:    "정상 파일 읽기",
			path:    testPath,
			want:    testContent,
			wantErr: false,
		},
		{
			name:    "존재하지 않는 파일",
			path:    "없는파일.txt",
			want:    "",
			wantErr: true,
		},
	}

	fp := parser.NewFileParser()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := fp.ReadContent(tt.path)

			// 에러 검증
			if (err != nil) != tt.wantErr {
				t.Errorf("ReadContent() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			// 결과 검증
			if got != tt.want {
				t.Errorf("ReadContent() = %v, want %v", got, tt.want)
			}
		})
	}
}

```
