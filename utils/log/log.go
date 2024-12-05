package log

import (
	"log"
	"os"
)

var (
	f      *os.File
	err    error
	logger *log.Logger
)

func init() {
	f, err = os.OpenFile("todo-log-stream.log",
		os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}

	logger = log.New(f, "[Logger]", log.LstdFlags)
}

func GetLogger() *log.Logger {
	return logger
}

func Error(message string, v ...interface{}) {
	logger.Printf(" [Error] "+message, v)
}

func Info(message string, v ...interface{}) {
	logger.Printf(" [INFO] "+message, v)
}

func Debug(message string, v ...interface{}) {
	logger.Printf(" [DEBUG] "+message, v)
}
