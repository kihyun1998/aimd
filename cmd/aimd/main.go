package main

import (
	"fmt"
	"log"
	"os"

	"github.com/kihyun1998/aimd/internal/config"
)

func main() {
	// 커스텀 usage 메시지 설정
	config.SetUsage(os.Args[0])

	// 설정 파싱
	cfg, err := config.ParseFlags()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("설정: %+v\n", cfg)
}
