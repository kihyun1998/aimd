# CodeMD (Code to Markdown)

CodeMD는 프로젝트의 소스 코드를 자동으로 마크다운 문서로 변환하는 CLI 도구입니다.

## 주요 기능

- 다양한 프로그래밍 언어 파일 지원
- 재귀적 디렉토리 탐색
- 확장자 기반 파일 필터링
- 커스텀 템플릿 지원
- 숨김 파일/디렉토리 처리

## 설치 방법

```bash
go get github.com/kihyun1998/codemd
```

## 사용 방법

기본 사용:
```bash
codemd -type go
```

추가 옵션 사용:
```bash
codemd -type go,java -exclude vendor,node_modules -out docs/CODE.md
```

### 옵션 설명

- `-type`: 처리할 파일 확장자 (필수, 쉼표로 구분)
- `-out`: 출력 파일 경로 (기본값: CODE.md)
- `-exclude`: 제외할 디렉토리 (선택, 쉼표로 구분)

## 프로젝트 구조

```
codemd/
├── cmd/
│   └── codemd/            # 실행 파일 디렉토리
│       └── main.go        # 메인 진입점
├── internal/
│   ├── config/            # 설정 관련
│   ├── generator/         # 마크다운 생성
│   └── parser/           # 파싱 관련
├── pkg/
│   └── utils/            # 유틸리티
└── test/                 # 테스트 코드
```

## 개발 환경

- Go 1.20 이상
- 모듈 기반 의존성 관리

## 테스트

전체 테스트 실행:
```bash
go test ./...
```

특정 패키지 테스트:
```bash
go test ./internal/parser
```

## 라이선스

MIT License

## 기여하기

1. Fork the Project
2. Create your Feature Branch (`git checkout -b feature/AmazingFeature`)
3. Commit your Changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to the Branch (`git push origin feature/AmazingFeature`)
5. Open a Pull Request

## 작성자

- kihyun1998