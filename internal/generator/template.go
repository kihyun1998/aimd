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
	Structure   string
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
