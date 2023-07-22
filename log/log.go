package log

import (
	"os"

	log "github.com/sirupsen/logrus"
)

func Init() {
	file := "/tmp/thriftls.log"
	logFile, err := os.OpenFile(file, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0766)
	if err != nil {
		panic(err)
	}
	log.SetLevel(log.DebugLevel)
	log.SetOutput(logFile) // 将文件设置为log输出的文件
}

type Logger struct {

}
