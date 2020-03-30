package log

import (
	"os"
	"runtime"

	"github.com/ResultadosDigitais/x9/config"
	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
)

type Logger struct {
	Logger *logrus.Logger
}

var logger Logger

func Init() {
	level := logrus.InfoLevel
	if *config.Opts.Debug {
		level = logrus.DebugLevel
	}
	logger = Logger{
		Logger: &logrus.Logger{
			Out:       os.Stdout,
			Hooks:     make(logrus.LevelHooks),
			Level:     level,
			Formatter: &logrus.JSONFormatter{},
		},
	}
}

func Fatal(msg string, extraFields map[string]interface{}) {
	logger.Logger.WithFields(getFields(extraFields)).Fatal(msg)
}

func Error(msg string, extraFields map[string]interface{}) {
	logger.Logger.WithFields(getFields(extraFields)).Error(msg)
}

func Warn(msg string, extraFields map[string]interface{}) {
	logger.Logger.WithFields(getFields(extraFields)).Warn(msg)
}

func Info(msg string, extraFields map[string]interface{}) {
	logger.Logger.WithFields(getFields(extraFields)).Info(msg)

}

func Debug(msg string, extraFields map[string]interface{}) {
	logger.Logger.WithFields(getFields(extraFields)).Debug(msg)
}

func getFields(extraFields map[string]interface{}) logrus.Fields {
	programCounter, file, line, _ := runtime.Caller(2)
	function := runtime.FuncForPC(programCounter)
	fields := log.Fields{
		"file":     file,
		"line":     line,
		"function": function.Name(),
	}
	for k, v := range extraFields {
		fields[k] = v
	}
	return fields
}
