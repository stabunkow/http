package logging

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"stabunkow/http/pkg/setting"
	"time"
)

type Level int

const (
	DEBUG Level = iota
	INFO
	WARN
	ERROR
	FATAL
)

var (
	F                  *os.File
	DefaultPrefix      = ""
	DefaultCallerDepth = 2

	levelFlags = []string{"DEBUG", "INFO", "WARN", "ERROR", "FATAL"}

	logger *log.Logger
)

func Setup() {
	var err error
	path := getLogFilePath() + "/" + getLogFileName()

	F, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		log.Fatalf("logging.Setup err: %v", err)
	}

	logger = log.New(F, DefaultPrefix, log.LstdFlags)
}

func getLogFilePath() string {
	return setting.AppSetting.LogPath
}

func getLogFileName() string {
	return fmt.Sprintf("%s-%s.log",
		setting.AppSetting.LogSaveName,
		time.Now().Format("2006-01-02"),
	)
}

// setPrefix set the prefix of the log output
func setPrefix(level Level) {
	_, file, line, ok := runtime.Caller(DefaultCallerDepth)

	logPrefix := ""
	if ok {
		logPrefix = fmt.Sprintf("[%s][%s:%d]", levelFlags[level], filepath.Base(file), line)
	} else {
		logPrefix = fmt.Sprintf("[%s]", levelFlags[level])
	}

	logger.SetPrefix(logPrefix)
}

func Debug(v ...interface{}) {
	setPrefix(DEBUG)
	logger.Println(v...)
}

func Info(v ...interface{}) {
	setPrefix(INFO)
	logger.Println(v...)
}

func Warn(v ...interface{}) {
	setPrefix(WARN)
	logger.Println(v...)
}

func Error(v ...interface{}) {
	setPrefix(ERROR)
	logger.Println(v...)
}

func Fatal(v ...interface{}) {
	setPrefix(FATAL)
	logger.Fatalln(v...)
}
