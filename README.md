# CodeMD (Auto Index Markdown Generator)

자동으로 프로젝트의 디렉토리 구조를 분석하여 마크다운 문서를 생성하는 도구입니다.

## 주요 기능
- 지정된 확장자의 파일만 필터링 (예: .go, .java)
- 특정 디렉토리 제외 기능 (예: vendor, node_modules)
- 커스텀 출력 경로 지정
- 재귀적 디렉토리 탐색
- 숨김 파일/디렉토리 처리

## 프로세스 Flow

```mermaid
graph TD
    A[시작] --> B[설정 파싱]
    B --> C{필수 인자 확인}
    C -->|실패| D[에러 출력]
    D --> E[종료]
    C -->|성공| F[버전 체크]
    F -->|버전 출력 요청| G[버전 출력]
    G --> E
    F -->|일반 실행| H[디렉토리 스캔]
    H --> I[codeignore 처리]
    I --> J[파일 필터링]
    J --> K[내용 분석]
    K --> L[마크다운 생성]
    L --> M[파일 저장]
    M --> E[종료]
```

## 설치 방법

### 소스코드에서 빌드
```bash
# 프로젝트 클론
git clone https://github.com/kihyun1998/codemd.git

# 프로젝트 디렉토리로 이동
cd codemd

# 빌드 
go build -o codemd.exe ./cmd/codemd

# 다양한 환경에서 빌드(cmd)
GOOS=linux GOARCH=amd64 go build -o codemd-linux-amd64 ./cmd/codemd
GOOS=darwin GOARCH=amd64 go build -o codemd-darwin-amd64 ./cmd/codemd
GOOS=windows GOARCH=amd64 go build -o codemd-windows-amd64.exe ./cmd/codemd

# 다양한 환경에서 빌드(powershell)
$env:GOOS="linux"; $env:GOARCH="amd64"; go build -o ./build/codemd-linux-amd64 ./cmd/codemd
$env:GOOS="darwin"; $env:GOARCH="amd64"; go build -o ./build/codemd-darwin-amd64 ./cmd/codemd
$env:GOOS="windows"; $env:GOARCH="amd64"; go build -o ./build/codemd-windows-amd64.exe ./cmd/codemd
```

### Go Install 사용
```bash
go install github.com/kihyun1998/codemd@latest
```

## 사용법

### 기본 사용
```bash
codemd -type go

# 버전 확인
codemd -v
# 또는
codemd --version
```

### 추가 옵션 사용
```bash
# 출력 경로 지정
codemd -type go -out docs/index.md

# 특정 디렉토리 제외
codemd -type go -exclude vendor,node_modules

# 여러 확장자 지정
codemd -type go,java,py -out docs/CODE.md
```

### 옵션 설명
- `-type, -t`: 처리할 파일 확장자 (필수, 쉼표로 구분)
- `-out, -o`: 출력 파일 경로 (기본값: CODE.md)
- `-exclude, -e`: 제외할 디렉토리 (선택, 쉼표로 구분)
- `-version, -v`: 버전 정보 출력
- `-codeignore, -c`: .codeignore 파일 사용 여부 (기본값: false)


## 테스트 실행
```bash
# 전체 테스트
go test ./...

# 특정 패키지 테스트
go test ./internal/parser
```

## 기여하기
1. Fork the Project
2. Create your Feature Branch
3. Commit your Changes
4. Push to the Branch
5. Open a Pull Request

## 라이선스
MIT License

## 작성자
- kihyun1998