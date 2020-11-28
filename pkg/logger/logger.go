package logger

import (
	"log"
	"os"
)

type level int

const (
	LevelDebug    = level(0)
	LevelInfo     = level(1)
	LevelWarning  = level(2)
	LevelError    = level(3)
	LevelCritical = level(4)
)

var CurrentLevel = LevelWarning

func Debug(msg string) {
	if CurrentLevel <= LevelDebug {
		logFormatted(msg, "d")
	}
}

func Info(msg string) {
	if CurrentLevel <= LevelInfo {
		logFormatted(msg, "i")
	}
}

func Warning(msg string) {
	if CurrentLevel <= LevelWarning {
		logFormatted(msg, "w")
	}
}

func Error(err error) {
	if CurrentLevel <= LevelError {
		logFormatted(err, "e")
	}
}

func Critical(err error) {
	if CurrentLevel <= LevelCritical {
		logFormatted(err, "c")
		os.Exit(1)
	}
}

func logFormatted(msg interface{}, lv string) {
	log.Printf("[%v] %s", lv, msg)
}
