# CodeMD v1.3.0

## 🎉 새로운 기능

### 대용량 파일 자동 분할 기능
- 출력 파일이 지정된 크기를 초과할 경우 자동으로 여러 파일로 분할
- 기본 최대 파일 크기는 10MB로 설정
- 분할된 파일은 자동으로 넘버링 (예: CODE1.md, CODE2.md, CODE3.md)

### 새로운 CLI 옵션
- `-maxsize` 또는 `-m`: 출력 파일의 최대 크기를 MB 단위로 설정
  ```bash
  # 예시
  codemd -type go -maxsize 20  # 20MB 크기로 분할
  codemd -t go -m 15          # 15MB 크기로 분할
  ```

## 🔧 기술적 개선사항
- SOLID 원칙을 준수한 파일 분할 시스템 구현
- 클린 아키텍처 기반의 모듈화 설계
- 메모리 사용 최적화
- 대용량 파일 처리 성능 개선

## 📚 사용 예시
```bash
# 기본 사용 (10MB 제한)
codemd -type go

# 20MB 크기로 파일 분할 설정
codemd -type go -maxsize 20

# 여러 옵션 조합
codemd -type go,java -maxsize 15 -exclude vendor,node_modules
```

## 🔍 결과 예시
프로젝트 크기가 25MB인 경우:
```
CODE1.md (10MB)
CODE2.md (10MB)
CODE3.md (5MB)
```

## ⚠️ 주의사항
- 파일 분할 크기는 0보다 커야 합니다
- 분할된 파일은 순차적으로 넘버링됩니다
- 기존 동일 이름의 파일이 있다면 덮어쓰기됩니다

## 📋 Documentation
자세한 사용법은 [README.md](./README.md)를 참조해주세요.

## 🔄 업그레이드 방법
```bash
go install github.com/kihyun1998/codemd@latest
```

## 🐛 버그 리포트
버그를 발견하시면 GitHub Issues에 리포트해주세요.
