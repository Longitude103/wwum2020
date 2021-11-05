package logging

import (
	"log"
	"os"
)

type TheLogger struct {
	file *os.File
	lg   *log.Logger
	li   *log.Logger
	ld   *log.Logger
	le   *log.Logger
}

func (l *TheLogger) Info(message string) {
	l.li.Println(message)
}

func (l *TheLogger) Debug(message string) {
	l.ld.Println(message)
}

func (l *TheLogger) Infof(template string, args ...interface{}) {
	l.li.Printf(template, args...)
}

// Errorf uses fmt.Printf to log messages
func (l *TheLogger) Errorf(template string, args ...interface{}) {
	l.le.Printf(template, args...)
}

func (l *TheLogger) Error(message string) {
	l.le.Println(message)
}

func (l *TheLogger) close() {
	_ = l.file.Close()
}

func NewLogger(fileName string) *TheLogger {
	f, _ := os.OpenFile(fileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	l := TheLogger{file: f}

	logger := log.New(f, "prefix: ", log.LstdFlags)
	l.lg = logger

	loggerInfo := log.New(f, "Info: ", log.LstdFlags)
	l.li = loggerInfo

	lDebug := log.New(f, "Debug: ", log.LstdFlags)
	l.ld = lDebug

	lError := log.New(f, "Error: ", log.LstdFlags)
	l.le = lError

	return &l
}
