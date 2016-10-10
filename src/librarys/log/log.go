package log

import (
	"conf"
	"log"
	"os"
)

var logger *log.Logger

func init() {
	var logPath string
	logPath, ok := conf.Conf["logPath"].(string)
	if !ok {
		return
	}
	logFile, err := os.OpenFile(logPath, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0666)

	if err != nil {
		log.Fatalln("open log file failed!")
	}

	logger = log.New(logFile, "",
		log.Ldate|
			log.Ltime|
			log.Lmicroseconds|
			log.Llongfile)
}

func Fatal(v ...interface{}) {
	logger.SetPrefix("[Fatal]")
	logger.Println(v...)
}

func FatalChk(err error, v ...interface{}) {
	if err != nil {
		Fatal(err, v)
	}
}

func Debug(v ...interface{}) {
	logger.SetPrefix("[Fatal]")
	logger.Print(v)
}
