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
	fileParser parser.FileParser
	processor  *templateProcessor
	outputPath string
	rootDir    string
}

// 생성자
func NewMarkdownGenerator(fp parser.FileParser, outputPath string) MarkdownGenerator {
	rootDir, err := os.Getwd()
	if err != nil {
		rootDir = ""
	}

	return &markdownGenerator{
		fileParser: fp,
		outputPath: outputPath,
		rootDir:    rootDir,
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

	result, err := mg.processor.Execute(TemplateData{Files: fileDataList})
	if err != nil {
		return err
	}

	return os.WriteFile(mg.outputPath, []byte(result), 0644)
}
