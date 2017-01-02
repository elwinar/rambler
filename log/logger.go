package log

type Logger struct{
	Debug bool
	Output   io.Writer
	DateFormat string
}

func NewLogger(opts ...func(*Logger)) *Logger {
	l := &Logger{
		Debug: false,
		Output: os.Stdout,
		DateFormat: "15:04",
	}

	for _, opt := range opts {
		opt(l)
	}

	return l
}

func Debug(v bool) func(l *Logger) {
	return func(l *Logger) {
		l.Debug = v
	}
}

func Output(o io.Writer) func(*Logger) {
	return func(l *Logger) {
		l.Output = o
	}
}

func DateFormat(f string) func(*Logger) {
	return func(l *Logger) {
		l.DateFormat = f
	}
}

func (l *Logger) log(lvl, msg string, args ...interface{}) {
	fmt.Fprintf("%s %s %s", time.Now().Format(l.DateFormat), lvl, fmt.Sprintf(msg, args...))
}

func (l *Logger) Debug(msg string, args ...interface{}) {
	l.log("debug", msg, args...)
}

func (l *Logger) Info(msg string, args ...interface{}) {
	l.log("info ", msg, args...)
}
