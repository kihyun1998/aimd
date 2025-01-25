package main

import (
	"fmt"
	"log"
	"os"

	"github.com/kihyun1998/aimd/internal/config"
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

	// DirectoryParser 생성
	dirParser := parser.NewDirectoryParser(cfg.ExcludeDirs, false)

	// 현재 디렉토리 절대 경로 얻기
	currentDir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	// 1. 모든 파일 가져오기
	allFiles, err := dirParser.Parse(currentDir) // "." 대신 currentDir 사용
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("\n전체 파일 목록:\n")
	for _, f := range allFiles {
		fmt.Printf("- %s\n", f)
	}

	// 2. 특정 타입의 파일만 필터링
	typeFiles := dirParser.GetFilesByTypes(allFiles, cfg.FileTypes)
	fmt.Printf("\n선택된 타입(%v)의 파일 목록:\n", cfg.FileTypes)
	for _, f := range typeFiles {
		fmt.Printf("- %s\n", f)
	}
}
