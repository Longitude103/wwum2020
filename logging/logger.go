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

// Info is for information that you want to log, uses fmt.Println for the message
func (l *TheLogger) Info(message string) {
	l.li.Println(message)
}

// Debug is for debug items you want to log, uses fmt.Println for the message
func (l *TheLogger) Debug(message string) {
	l.ld.Println(message)
}

// Debugf is for debugging with a template string and args, uses fmt.Printf for the output
func (l *TheLogger) Debugf(template string, args ...interface{}) {
	l.ld.Printf(template, args...)
}

// Infof is for logging information with a template string and args, uses fmt.Printf for the output
func (l *TheLogger) Infof(template string, args ...interface{}) {
	l.li.Printf(template, args...)
}

// Errorf is for error output with a template string and args, uses fmt.Printf to log messages
func (l *TheLogger) Errorf(template string, args ...interface{}) {
	l.le.Printf(template, args...)
}

// Error is for error messages, uses fmt.Println for messages
func (l *TheLogger) Error(message string) {
	l.le.Println(message)
}

// Close closes the log file that is uses
func (l *TheLogger) Close() {
	_ = l.file.Close()
}

// NewLogger returns a new instance of TheLogger and takes in a fileName for the log file.
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
