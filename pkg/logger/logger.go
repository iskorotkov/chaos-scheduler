package logger

import (
	"log"
	"os"
	"strings"
)

type level int

const (
	levelDebug    = level(0)
	levelInfo     = level(1)
	levelWarning  = level(2)
	levelError    = level(3)
	levelCritical = level(4)
)

var currentLevel = levelWarning

func SetLevel(level string) {
	level = strings.ToLower(level)

	switch level {
	case "0", "debug":
		currentLevel = levelDebug
	case "1", "info":
		currentLevel = levelInfo
	case "2", "warning":
		currentLevel = levelWarning
	case "3", "error":
		currentLevel = levelError
	case "4", "critical":
		currentLevel = levelCritical
	default:
		Warning("couldn't parse log level")
	}
}

func Debug(msg string) {
	if currentLevel <= levelDebug {
		logFormatted(msg, "d")
	}
}

func Info(msg string) {
	if currentLevel <= levelInfo {
		logFormatted(msg, "i")
	}
}

func Warning(msg string) {
	if currentLevel <= levelWarning {
		logFormatted(msg, "w")
	}
}

func Error(err error) {
	if currentLevel <= levelError {
		logFormatted(err, "e")
	}
}

func Critical(err error) {
	if currentLevel <= levelCritical {
		logFormatted(err, "c")
		os.Exit(1)
	}
}

func logFormatted(msg interface{}, lv string) {
	log.Printf("[%v] %s", lv, msg)
}
