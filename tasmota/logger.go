package tasmota

import "log"

func info(module string, message string, v ...interface{}) {
	log.Printf("["+module+"]\t INFO: "+message, v...)
}

func debug(module string, message string, v ...interface{}) {
	log.Printf("["+module+"]\t DEBUG: "+message, v...)
}

func fatal(module string, message string, v ...interface{}) {
	log.Fatalf("["+module+"]\t DEBUG: "+message, v...)
}
