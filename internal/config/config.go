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
	useCodeIgnore := flag.Bool("codeignore", true, ".codeignore 파일 사용 여부")

	flag.Parse()

	// -type 플래그가 없으면 에러
	if *types == "" {
		return nil, fmt.Errorf("필수 플래그 누락: -type")
	}

	return &Config{
		FileTypes:     strings.Split(*types, ","),
		OutputPath:    *output,
		ExcludeDirs:   strings.Split(*exclude, ","),
		UseCodeIgnore: *useCodeIgnore,
	}, nil
}
