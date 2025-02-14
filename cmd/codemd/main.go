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
	mdGen := generator.NewMarkdownGenerator(fileParser, cfg.OutputPath, cfg.MaxFileSizeMB)

	// 기본 템플릿 설정
	defaultTemplate := "# {{.ProjectName}}\n{{.Structure}}{{range .Files}}## {{.Path}}\n```{{if .Extension}}{{.Extension}}{{end}}\n{{.Content}}\n```\n{{end}}"

	if err := mdGen.SetTemplate(defaultTemplate); err != nil {
		log.Fatal(err)
	}
	// 마크다운 생성
	if err := mdGen.Generate(typeFiles); err != nil {
		log.Fatal(err)
	}
}
