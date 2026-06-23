package logger

import(
	"log"
	"os"
)

type Logger struct{
	infolog *log.Logger
	errorLog *log.Logger
}

func NewLogger() *Logger{
	return &Logger{
		infolog: log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime),
		errorLog: log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile),
	}
}

func (l *Logger) Info(format string, v ...interface{}){
	l.infolog.Printf(format, v...)
}

func (l *Logger) Error(format string, v ...interface{}){
	l.errorLog.Printf(format, v...)
}