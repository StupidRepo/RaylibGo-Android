//go:build (linux || windows || darwin) && !android

package platformspecifics

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

type Platform struct{}

const (
	GLSLVersion = 330
)

func (p *Platform) GetOS() PlatformEnum {
	return PlatformDesktop
}

func (p *Platform) GetWindowSize() (int32, int32) {
	return int32(rl.GetRenderWidth()), int32(rl.GetRenderHeight())
}
func (p *Platform) GetInsets() Insets {
	return Insets{} // desktop has none
}

func (p *Platform) LogIt(priority int, tag string, text string) int {
	// for desktop, we'll just print to the console
	logMessage := "[" + tag + "] " + text
	switch priority {
	case AndroidLogVerbose, AndroidLogDebug:
		println("DEBUG: " + logMessage)
	case AndroidLogInfo:
		println("INFO: " + logMessage)
	case AndroidLogWarn:
		println("WARN: " + logMessage)
	case AndroidLogError:
		println("ERROR: " + logMessage)
	case AndroidLogFatal:
		println("FATAL: " + logMessage)
	default:
		println("UNKNOWN: " + logMessage)
	}
	return 0 // success
}
