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
	dirParser := parser.NewDirectoryParser(cfg.ExcludeDirs)

	// 현재 디렉토리에서 파일 목록 가져오기
	files, err := dirParser.Parse(".")
	if err != nil {
		log.Fatal(err)
	}

	// 각 확장자별로 필터링하여 출력
	for _, ext := range cfg.FileTypes {
		filtered := dirParser.FilterByExtenstion(files, ext)
		fmt.Printf("\n%s 파일 목록:\n", ext)
		for _, f := range filtered {
			fmt.Printf("- %s\n", f)
		}
	}
}
