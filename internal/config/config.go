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
	ShowVersion   bool
}

// Usage 메시지 설정
func SetUsage(programName string) {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "사용법: %s [옵션]\n\n옵션:\n", programName)
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\n예시:\n")
		fmt.Fprintf(os.Stderr, "  %s -version\n", programName)
		fmt.Fprintf(os.Stderr, "  %s -type go,java\n", programName)
		fmt.Fprintf(os.Stderr, "  %s -type go -exclude vendor,node_modules\n", programName)
	}
}

// 플래그 파싱
func ParseFlags() (*Config, error) {
	var (
		types         string
		output        string
		exclude       string
		useCodeIgnore bool
		showVersion   bool
	)

	flag.StringVar(&types, "type", "", "파일 확장자들 (쉼표로 구분)")
	flag.StringVar(&types, "t", "", "파일 확장자들 (쉼표로 구분) (짧은 버전)")

	flag.StringVar(&output, "out", "CODE.md", "출력 파일 경로")
	flag.StringVar(&output, "o", "CODE.md", "출력 파일 경로 (짧은 버전)")

	flag.StringVar(&exclude, "exclude", "", "제외할 디렉토리들 (쉼표로 구분)")
	flag.StringVar(&exclude, "e", "", "제외할 디렉토리들 (쉼표로 구분) (짧은 버전)")

	flag.BoolVar(&useCodeIgnore, "codeignore", false, ".codeignore 파일 사용 여부")
	flag.BoolVar(&useCodeIgnore, "c", false, ".codeignore 파일 사용 여부 (짧은 버전)")

	flag.BoolVar(&showVersion, "version", false, "버전 정보 출력")
	flag.BoolVar(&showVersion, "v", false, "버전 정보 출력 (짧은 버전)")

	flag.Parse()

	// -v 또는 -version 플래그만 있는 경우
	if showVersion && len(os.Args) == 2 {
		return &Config{ShowVersion: true}, nil
	}

	return &Config{
		FileTypes:     strings.Split(types, ","),
		OutputPath:    output,
		ExcludeDirs:   strings.Split(exclude, ","),
		UseCodeIgnore: useCodeIgnore,
		ShowVersion:   showVersion,
	}, nil
}
