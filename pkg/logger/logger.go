package logger

import (
	"log"
	"os"
)

type Log struct {
	Info *log.Logger
	Err  *log.Logger
}

func New() *Log {
	logInfo := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	logError := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime)
	return &Log{
		Info: logInfo,
		Err:  logError,
	}
}
