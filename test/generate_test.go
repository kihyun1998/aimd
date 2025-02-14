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
			mg := generator.NewMarkdownGenerator(fp, outputPath, 10)

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
