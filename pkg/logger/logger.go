package logger

import (
	"log"
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
		log.Println(msg)
	}
}

func Info(msg string) {
	if CurrentLevel <= LevelInfo {
		log.Println(msg)
	}
}

func Warning(msg string) {
	if CurrentLevel <= LevelWarning {
		log.Println(msg)
	}
}

func Error(err error) {
	if CurrentLevel <= LevelError {
		log.Println(err)
	}
}

func Critical(err error) {
	if CurrentLevel <= LevelCritical {
		log.Fatalln(err)
	}
}
