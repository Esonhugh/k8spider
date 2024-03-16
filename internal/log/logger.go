package log

import (
	"io"
	"os"

	log "github.com/sirupsen/logrus"
	prefixed "github.com/x-cray/logrus-prefixed-formatter"
)

const (
	LogLevelTrace = "trace"
	LogLevelDebug = "debug"
	LogLevelInfo  = "info"
	LogLevelWarn  = "warn"
	LogLevelError = "error"
	LogLevelFatal = "fatal"
	LogLevelPanic = "panic"
)

// Init func is a function to init logrus with specific log level
func Init(level string) {
	var logWriter io.Writer
	if logLevel(level) >= log.DebugLevel {
		f, e := os.OpenFile("debug.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0600)
		if e == nil {
			logWriter = io.MultiWriter(os.Stdout, f)
		} else {
			logWriter = os.Stdout
		}
	} else {
		logWriter = os.Stdout
	}
	log.SetOutput(logWriter)
	log.SetFormatter(logFormat())
	log.SetLevel(logLevel(level))
}

// logLevel search level strings return correct Level
func logLevel(level string) log.Level {
	switch level {
	case LogLevelTrace:
		return log.TraceLevel
	case LogLevelDebug:
		return log.DebugLevel
	case LogLevelInfo:
		return log.InfoLevel
	case LogLevelWarn:
		return log.WarnLevel
	case LogLevelError:
		return log.ErrorLevel
	case LogLevelFatal:
		return log.FatalLevel
	case LogLevelPanic:
		return log.PanicLevel
	default:
		return log.InfoLevel
	}
}

// logFormat sets log format by using prefixed "x-cray/logrus-prefixed-formatter"
func logFormat() log.Formatter {
	formatter := new(prefixed.TextFormatter)
	formatter.FullTimestamp = true
	formatter.TimestampFormat = "2006-01-02 15:04:05"
	formatter.SetColorScheme(&prefixed.ColorScheme{
		PrefixStyle:    "blue+b",
		TimestampStyle: "white+h",
	})
	return formatter
}
