package cache

import (
	"runtime"
	"strings"
)

// this is used to get the caller function for which we should log the latency
func getRequiredCallerFunctionName() string {
	result := "unknown"
	counters := make([]uintptr, 8)
	n := runtime.Callers(0, counters)
	if n > 0 {
		frames := runtime.CallersFrames(counters[:n])
		for {
			f, ok := frames.Next()
			if !ok {
				break
			}
			if strings.Contains(f.Function, "github.com/angel-one/go-cache-client") {
				s := strings.Split(f.Function, ".")
				result = s[len(s)-1]
			}
		}
	}
	return result
}
