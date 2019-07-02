package tasmota

import "log"

type Loglevel int

const (
	DEBUG Loglevel = 1
	INFO  Loglevel = 2
	WARN  Loglevel = 3
	ERROR Loglevel = 4
	OFF   Loglevel = 5
)

var Level = INFO

func info(module string, message string, v ...interface{}) {
	if Level <= INFO {
		log.Printf("["+module+"]\t INFO: "+message, v...)
	}
}

func debug(module string, message string, v ...interface{}) {
	if Level <= DEBUG {
		log.Printf("["+module+"]\t DEBUG: "+message, v...)
	}
}

func fatal(module string, message string, v ...interface{}) {
	if Level <= ERROR {
		log.Printf("["+module+"]\t FATAL: "+message, v...)
	}
}

func warn(module string, message string, v ...interface{}) {
	if Level <= WARN {
		log.Printf("["+module+"]\t WARN: "+message, v...)
	}
}
