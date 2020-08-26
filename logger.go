package main

import "log"
import "os"
import "time"

func InitLog(appName string) {

	goPath := os.Getenv("GOPATH")

	var logDir = goPath + "/src/" + appName + "/logs"

	if isExist(logDir) != true {
		mkdirErr := os.Mkdir(logDir, 0755)

		if mkdirErr != nil {
			panic("Mkdir [" + logDir + "] fail!")
		}
	}

	var logFile = logDir + "/error." + time.Now().Format(dateFormat) + ".log"

	file, err := os.OpenFile(logFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)

	if err != nil {
		panic("open file [" + logFile + "] fall!")
	}

	// os.Chmod(logFile, 0644)

	var prefix string = "[error]"

	log.SetOutput(file)
	log.SetPrefix(prefix)
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}

func isExist(path string) bool {
	_, err := os.Stat(path)

	if err != nil {
		if os.IsExist(err) {
			return true
		}

		if os.IsNotExist(err) {
			return false
		}
		return false
	}

	return true
}
