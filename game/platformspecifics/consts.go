package platformspecifics

//goland:noinspection GoNameStartsWithPackageName
const (
	NDKTag  = "NDKBridge"
	GameTag = "RaylibGoGame"
)
const (
	AndroidLogUnknown = iota
	AndroidLogDefault
	AndroidLogVerbose
	AndroidLogDebug
	AndroidLogInfo
	AndroidLogWarn
	AndroidLogError
	AndroidLogFatal
	AndroidLogSilent
)

type Insets struct {
	Left   int32
	Top    int32
	Right  int32
	Bottom int32
}

type PlatformFunctions interface {
	GetWindowSize() (int32, int32)
	GetInsets() Insets

	LogIt(priority int, tag string, text string) int
}
