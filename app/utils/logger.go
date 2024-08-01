package utils

import (
	"log"
	"os"
)

func ErrorLog() *log.Logger {
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	return errorLog

}

func InfoLog() *log.Logger {
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	return infoLog
}
