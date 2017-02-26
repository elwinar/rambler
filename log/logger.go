package log

import (
	"fmt"
	"io"
	"os"
	"time"
)

type Logger struct {
	PrintDebug bool
	Output     io.Writer
	DateFormat string
}

func NewLogger(opts ...func(*Logger)) *Logger {
	l := &Logger{
		PrintDebug: false,
		Output:     os.Stdout,
		DateFormat: "15:04",
	}

	for _, opt := range opts {
		opt(l)
	}

	return l
}

func (l *Logger) log(lvl, msg string, args ...interface{}) {
	fmt.Fprintf(l.Output, "%s %s %s\n", time.Now().Format(l.DateFormat), lvl, fmt.Sprintf(msg, args...))
}

func (l *Logger) Debug(msg string, args ...interface{}) {
	if l.PrintDebug {
		l.log("debug", msg, args...)
	}
}

func (l *Logger) Info(msg string, args ...interface{}) {
	l.log("info ", msg, args...)
}
