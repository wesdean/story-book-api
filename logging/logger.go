package logging

import (
	"github.com/google/logger"
	"github.com/gorilla/context"
	"gopkg.in/guregu/null.v3"
	"net/http"
	"os"
	"strings"
)

type Logger struct {
	*logger.Logger
	Name string
}

type LoggerConfig struct {
	Name     null.String
	LogPath  null.String
	Verbose  null.Bool
	ClearLog null.Bool
}

const LOGLEVEL_INFO = 1
const LOGLEVEL_WARNING = 2
const LOGLEVEL_ERROR = 3

func NewLogger(logConfig *LoggerConfig) (*Logger, error) {
	var err error

	var writeType int
	if logConfig.ClearLog.ValueOrZero() {
		writeType = os.O_TRUNC
	} else {
		writeType = os.O_APPEND
	}
	logFile, err := os.OpenFile(logConfig.LogPath.ValueOrZero(), os.O_CREATE|os.O_WRONLY|writeType, 0660)
	if err != nil {
		return nil, err
	}

	_myLogger := logger.Init(logConfig.Name.ValueOrZero(), logConfig.Verbose.ValueOrZero(), true, logFile)
	myLogger := &Logger{Logger: _myLogger, Name: logConfig.Name.ValueOrZero()}

	myLogger.Info("logger initialized")

	return myLogger, nil
}

func CloseLogger(myLogger *Logger) {
	myLogger.Info("logger closed")
	myLogger.Close()
}

func GetLoggerFromRequest(r *http.Request) *Logger {
	logValue, ok := context.GetOk(r, "Logger")
	if ok {
		myLogger, ok := logValue.(*Logger)
		if ok {
			return myLogger
		}
	}
	return nil
}

func Log(myLogger *Logger, level int, v ...interface{}) {
	if myLogger == nil {
		return
	}

	switch level {
	case LOGLEVEL_INFO:
		myLogger.Info(v...)
	case LOGLEVEL_WARNING:
		myLogger.Warning(v...)
	case LOGLEVEL_ERROR:
		myLogger.Error(v...)
	}
}

func Logf(myLogger *Logger, level int, format string, v ...interface{}) {
	if myLogger == nil {
		return
	}

	switch level {
	case LOGLEVEL_INFO:
		myLogger.Infof(format, v)
	case LOGLEVEL_WARNING:
		myLogger.Warningf(format, v)
	case LOGLEVEL_ERROR:
		myLogger.Errorf(format, v)
	}
}

func (myLogger *Logger) MakeLogPrefix(v ...interface{}) (string, []interface{}, int) {
	vCount := len(v)
	return "%s: ", append([]interface{}{myLogger.Name}, v...), vCount
}

func (myLogger *Logger) Info(v ...interface{}) {
	s, v, vCount := myLogger.MakeLogPrefix(v...)
	myLogger.Logger.Infof(s+strings.Trim(strings.Repeat("%v ", vCount), " "), v...)
}

func (myLogger *Logger) Warning(v ...interface{}) {
	s, v, vCount := myLogger.MakeLogPrefix(v...)
	myLogger.Logger.Warningf(s+strings.Trim(strings.Repeat("%v ", vCount), " "), v...)
}

func (myLogger *Logger) Error(v ...interface{}) {
	s, v, vCount := myLogger.MakeLogPrefix(v...)
	myLogger.Logger.Errorf(s+strings.Trim(strings.Repeat("%v ", vCount), " "), v...)
}

func (myLogger *Logger) Infof(format string, v ...interface{}) {
	s, v, _ := myLogger.MakeLogPrefix(v...)
	myLogger.Logger.Infof(s+format, v...)
}

func (myLogger *Logger) Warningf(format string, v ...interface{}) {
	s, v, _ := myLogger.MakeLogPrefix(v...)
	myLogger.Logger.Warningf(s+format, v...)
}

func (myLogger *Logger) Errorf(format string, v ...interface{}) {
	s, v, _ := myLogger.MakeLogPrefix(v...)
	myLogger.Logger.Errorf(s+format, v...)
}
