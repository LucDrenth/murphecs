package utils

import (
	"fmt"
	"runtime"
	"strings"
)

func Caller(stackLevel int, numberOfPackagesToInclude int) string {
	_, fullPath, line, ok := runtime.Caller(stackLevel)
	if !ok {
		return ""
	}

	var path string
	pathSplit := strings.Split(fullPath, "/")
	if len(pathSplit) <= numberOfPackagesToInclude {
		path = strings.Join(pathSplit, "/")
	} else {
		path = strings.Join(pathSplit[len(pathSplit)-numberOfPackagesToInclude:], "/")
	}

	return fmt.Sprintf("%s:%d", path, line)
}
