package logging

import (
	"fmt"
	"kus/krzysztof/titler/environment"
	"log"
	"strings"
	"time"
)

// values for log levels
const (
	ERROR = iota
	WARNING
	INFO
	DEBUG
)

// map names to numbers to compare in logging function
var LogLevels map[string]int = map[string]int{
	"ERROR":   ERROR,
	"WARNING": WARNING,
	"INFO":    INFO,
	"DEBUG":   DEBUG,
}

var maxLevelLen int

// reverse log level map
var LogLevelsRev = make(map[int]string, len(LogLevels))

type logWriter struct{}

func (writer logWriter) Write(bytes []byte) (int, error) {
	return fmt.Print(formatTimestamp() + " " + string(bytes))
}

func InitLogging() {
	log.SetFlags(0)
	log.SetOutput(new(logWriter))
	for k, v := range LogLevels {
		LogLevelsRev[v] = k
	}
	maxLevelLen = countMaxLevelLen(&LogLevels)
}

func countMaxLevelLen(m *map[string]int) int {
	l := 0
	cl := 0
	for k := range *m {
		cl = len(k) // assuming ASCII
		if l < cl {
			l = cl
		}
	}
	return l
}

func Log(l int, s string) {
	// by default even when LOGLEVEL does not exists, is empty, does not match, LogLevels[...] will return 0 (ERROR)
	if LogLevels[environment.EnvVars["LOGLEVEL"]] >= l {
		log.Println("[" + LogLevelsRev[l] + strings.Repeat(" ", maxLevelLen-len(LogLevelsRev[l])) + "] " + s)
	}
}

func Logf(l int, s string, v ...any) {
	if LogLevels[environment.EnvVars["LOGLEVEL"]] >= l {
		log.Printf("["+LogLevelsRev[l]+strings.Repeat(" ", maxLevelLen-len(LogLevelsRev[l]))+"] "+s, v...)
	}
}

func formatTimestamp() string {
	t := time.Now().UTC()
	return fmt.Sprintf("%04d-%02d-%02dT%02d:%02d:%02d.%03d", t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second(), t.UnixMilli()%1000)
}
