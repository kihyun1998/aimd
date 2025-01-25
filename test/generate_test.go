package test

import (
	"testing"

	"github.com/kihyun1998/aimd/internal/generator"
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
