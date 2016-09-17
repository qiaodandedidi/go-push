package librarys

import (
	"conf"
	"log"
	"os"
)

var logger *log.Logger

func init() {
	logFile, err := os.OpenFile(conf.Conf["logPath"], os.O_WRONLY|os.O_APPEND|O_CREATE, 0666)

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
	logger.Fatal(v)
}

func FatalChk(err error) {
	if err != nil {
		Fatal(err)
	}
}

func Debug(v ...interface{}) {
	logger.SetPrefix("[Fatal]")
	logger.Print(v)
}
