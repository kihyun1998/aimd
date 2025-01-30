package version

// Version 정보
var (
	// 메이저 버전
	Major = "1"
	// 마이너 버전
	Minor = "1"
	// 패치 버전
	Patch = "0"
)

// GetVersion은 현재 버전 문자열을 반환합니다
func GetVersion() string {
	return Major + "." + Minor + "." + Patch
}

// GetVersionInfo는 버전 정보를 상세히 반환합니다
func GetVersionInfo() string {
	return "codemd version " + GetVersion()
}
