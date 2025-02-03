package osutils

import (
	"fmt"
	"runtime"
)

func GetTraceInfo(skip ...int) string {
	if len(skip) == 0 {
		skip = []int{1}
	}
	_, file, no, ok := runtime.Caller(skip[0])
	if ok {
		return fmt.Sprintf("%s:%d", file, no)
	}
	return "unknown"
}
