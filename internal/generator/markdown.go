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
