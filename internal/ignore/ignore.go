package ignore

// Pattern은 무시 패턴을 나타내는 인터페이스입니다
type Pattern interface {
	IsMatch(path string) bool
	IsNegative() bool
	IsDirectory() bool
}

// Ignorer는 파일 무시 규칙을 처리하는 인터페이스입니다
type Ignorer interface {
	AddPattern(pattern string) error
	ShouldIgnore(path string) bool
	LoadFromFile(path string) error
}
