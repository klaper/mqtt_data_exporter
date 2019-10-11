package logger

import "log"

type Loglevel int

const (
	DEBUG Loglevel = 1
	INFO  Loglevel = 2
	WARN  Loglevel = 3
	ERROR Loglevel = 4
	OFF   Loglevel = 5
)

var level = INFO

func SetLogLevel(newLevel Loglevel) {
	level = newLevel
}

func Info(module string, message string, v ...interface{}) {
	if level <= INFO {
		log.Printf("["+module+"]\t INFO: "+message, v...)
	}
}

func Debug(module string, message string, v ...interface{}) {
	if level <= DEBUG {
		log.Printf("["+module+"]\t DEBUG: "+message, v...)
	}
}

func Fatal(module string, message string, v ...interface{}) {
	if level <= ERROR {
		log.Printf("["+module+"]\t FATAL: "+message, v...)
	}
}

func Warn(module string, message string, v ...interface{}) {
	if level <= WARN {
		log.Printf("["+module+"]\t WARN: "+message, v...)
	}
}
