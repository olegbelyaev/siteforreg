package errorwrapper

import "log"

// Panic - выводит panic(msg + err.Error()), если err!=nil и len(msg)>0
func Panic(msg string, err error) {
	if err == nil {
		return
	}
	if len(msg) == 0 {
		return
	}
	panic(msg + " : " + err.Error())
}

// Log - prints msg + err.Error() to log if err!=nil
func Log(msg string, err error) {
	if err == nil {
		return
	}
	if len(msg) == 0 {
		return
	}
	log.Print(msg + " : " + err.Error())
}
