package errorwrapper

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
