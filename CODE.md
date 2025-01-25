# Project Documentation
## E:\aimd\cmd\aimd\main.go
```go
package main

import (
	"log"
	"os"

	"github.com/kihyun1998/aimd/internal/config"
	"github.com/kihyun1998/aimd/internal/generator"
	"github.com/kihyun1998/aimd/internal/parser"
)

func main() {
	// 커스텀 usage 메시지 설정
	config.SetUsage(os.Args[0])

	// 설정 파싱
	cfg, err := config.ParseFlags()
	if err != nil {
		log.Fatal(err)
	}

	// 파서 생성
	dirParser := parser.NewDirectoryParser(cfg.ExcludeDirs, false)
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
	defaultTemplate := "# Project Documentation\n{{range .Files}}## {{.Path}}\n```{{if .Extension}}{{.Extension}}{{end}}\n{{.Content}}\n```\n{{end}}"
	if err := mdGen.SetTemplate(defaultTemplate); err != nil {
		log.Fatal(err)
	}
	// 마크다운 생성
	if err := mdGen.Generate(typeFiles); err != nil {
		log.Fatal(err)
	}
}

```
## E:\aimd\internal\config\config.go
```go
package config

import (
	"flag"
	"fmt"
	"os"
	"strings"
)

type Config struct {
	FileTypes   []string
	OutputPath  string
	ExcludeDirs []string
}

// Usage 메시지 설정
func SetUsage(programName string) {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "사용법: %s [옵션]\n\n옵션:\n", programName)
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\n예시:\n")
		fmt.Fprintf(os.Stderr, "  %s -type go,java\n", programName)
		fmt.Fprintf(os.Stderr, "  %s -type go -exclude vendor,node_modules\n", programName)
	}
}

// 플래그 파싱
func ParseFlags() (*Config, error) {
	types := flag.String("type", "", "파일 확장자들 (쉼표로 구분)")
	output := flag.String("out", "CODE.md", "출력 파일 경로")
	exclude := flag.String("exclude", "", "제외할 디렉토리들 (쉼표로 구분)")

	flag.Parse()

	// -type 플래그가 없으면 에러
	if *types == "" {
		return nil, fmt.Errorf("필수 플래그 누락: -type")
	}

	return &Config{
		FileTypes:   strings.Split(*types, ","),
		OutputPath:  *output,
		ExcludeDirs: strings.Split(*exclude, ","),
	}, nil
}

```
## E:\aimd\internal\generator\markdown.go
```go
package generator

import (
	"os"
	"path/filepath"

	"github.com/kihyun1998/aimd/internal/parser"
)

type MarkdownGenerator interface {
	Generate(files []string) error
	SetTemplate(template string) error
}

// 마크다운 생성기 구조체
type markdownGenerator struct {
	fileParser parser.FileParser
	processor  *templateProcessor
	outputPath string
}

// 생성자
func NewMarkdownGenerator(fp parser.FileParser, outputPath string) MarkdownGenerator {
	return &markdownGenerator{
		fileParser: fp,
		outputPath: outputPath,
	}
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

		fileDataList = append(fileDataList, FileData{
			Path:      file,
			Content:   content,
			Extension: ext,
		})
	}

	result, err := mg.processor.Execute(TemplateData{Files: fileDataList})
	if err != nil {
		return err
	}

	return os.WriteFile(mg.outputPath, []byte(result), 0644)
}

```
## E:\aimd\internal\generator\template.go
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
	Files []FileData
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
## E:\aimd\internal\parser\directory.go
```go
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

// 모든 파일 가져오기
func (d *directoryParser) Parse(root string) ([]string, error) {
	var files []string

	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return NewParseError(path, err)
		}

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

```
## E:\aimd\internal\parser\error.go
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
## E:\aimd\internal\parser\file.go
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
## E:\aimd\internal\parser\parser.go
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
## E:\aimd\pkg\utils\file_utils.go
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
## E:\aimd\test\generate_test.go
```go
package test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/kihyun1998/aimd/internal/generator"
	"github.com/kihyun1998/aimd/internal/parser"
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
## E:\aimd\test\parser_test.go
```go
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
