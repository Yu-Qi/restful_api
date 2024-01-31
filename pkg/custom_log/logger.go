package customlog

import (
	"context"
	"io"

	"github.com/sirupsen/logrus"
)

type severity string

// Define key
const (
	// log severity
	SeverityDefault   severity = "DEFAULT"
	SeverityDebug     severity = "DEBUG"
	SeverityInfo      severity = "INFO"
	SeverityNotice    severity = "NOTICE"
	SeverityWarning   severity = "WARNING"
	SeverityError     severity = "ERROR"
	SeverityCritical  severity = "CRITICAL"
	SeverityAlert     severity = "ALERT"
	SeverityEmergency severity = "EMERGENCY"
)

func init() {
	logrus.SetFormatter(&logrus.JSONFormatter{
		FieldMap: logrus.FieldMap{
			logrus.FieldKeyMsg:   "message",
			logrus.FieldKeyLevel: "severity",
		},
		DisableTimestamp: true,
	})
	logrus.SetLevel(logrus.DebugLevel)
}

// SetOutput 设定日志输出
func SetOutput(out io.Writer) {
	logrus.SetOutput(out)
}

func writeLog(s severity, message interface{}, data interface{}) {
	logEntry := logrus.WithFields(logrus.Fields{
		"data": data,
	})
	if s == SeverityDebug {
		logEntry.Debug(message)
	} else if s == SeverityInfo {
		logEntry.Info(message)
	} else if s == SeverityWarning {
		logEntry.Warning(message)
	} else if s == SeverityError {
		logEntry.Error(message)
	} else if s == SeverityCritical {
		logEntry.Fatal(message)
	}
}

func withContext(logEntry *logrus.Entry, c context.Context) *logrus.Entry {
	return logEntry
}

func writeLogCtx(c context.Context, s severity, message interface{}, data interface{}) {
	logEntry := logrus.WithFields(logrus.Fields{
		"data": data,
	})
	if s == SeverityDebug {
		logEntry.Debug(message)
	} else if s == SeverityInfo {
		logEntry.Info(message)
	} else if s == SeverityWarning {
		logEntry.Warning(message)
	} else if s == SeverityError {
		logEntry.Error(message)
	} else if s == SeverityCritical {
		logEntry.Fatal(message)
	}
}

// DebugCtx .
func DebugCtx(c context.Context, message interface{}) {
	logEntry := withContext(logrus.WithFields(logrus.Fields{}), c)
	logEntry.Debug(message)
}

// DebugfCtx .
func DebugfCtx(c context.Context, format string, args ...interface{}) {
	logEntry := withContext(logrus.WithFields(logrus.Fields{}), c)
	logEntry.Debugf(format, args...)
}

// DebugWithData .
func DebugWithData(message interface{}, data interface{}) {
	writeLog(SeverityDebug, message, data)
}

// DebugWithDataCtx .
func DebugWithDataCtx(c context.Context, message interface{}, data interface{}) {
	writeLogCtx(c, SeverityDebug, message, data)
}

// InfoCtx .
func InfoCtx(c context.Context, message interface{}) {
	logEntry := withContext(logrus.WithFields(logrus.Fields{}), c)
	logEntry.Info(message)
}

// InfofCtx .
func InfofCtx(c context.Context, format string, args ...interface{}) {
	logEntry := withContext(logrus.WithFields(logrus.Fields{}), c)
	logEntry.Infof(format, args...)
}

// InfoWithData .
func InfoWithData(message interface{}, data interface{}) {
	writeLog(SeverityInfo, message, data)
}

// InfoWithDataCtx .
func InfoWithDataCtx(c context.Context, message interface{}, data interface{}) {
	writeLogCtx(c, SeverityInfo, message, data)
}

// WarningCtx .
func WarningCtx(c context.Context, message interface{}) {
	logEntry := withContext(logrus.WithFields(logrus.Fields{}), c)
	logEntry.Warning(message)
}

// WarningfCtx .
func WarningfCtx(c context.Context, format string, args ...interface{}) {
	logEntry := withContext(logrus.WithFields(logrus.Fields{}), c)
	logEntry.Warningf(format, args...)
}

// WarningWithData .
func WarningWithData(message interface{}, data interface{}) {
	writeLog(SeverityWarning, message, data)
}

// WarningWithDataCtx .
func WarningWithDataCtx(c context.Context, message interface{}, data interface{}) {
	writeLogCtx(c, SeverityWarning, message, data)
}

// ErrorCtx .
func ErrorCtx(c context.Context, message interface{}) {
	logEntry := withContext(logrus.WithFields(logrus.Fields{}), c)
	logEntry.Error(message)
}

// ErrorfCtx .
func ErrorfCtx(c context.Context, format string, args ...interface{}) {
	logEntry := withContext(logrus.WithFields(logrus.Fields{}), c)
	logEntry.Errorf(format, args...)
}

// ErrorWithData .
func ErrorWithData(message interface{}, data interface{}) {
	writeLog(SeverityError, message, data)
}

// ErrorWithDataCtx .
func ErrorWithDataCtx(c context.Context, message interface{}, data interface{}) {
	writeLogCtx(c, SeverityError, message, data)
}

// FatalCtx .
func FatalCtx(c context.Context, message interface{}) {
	logEntry := withContext(logrus.WithFields(logrus.Fields{}), c)
	logEntry.Fatal(message)
}

// FatalfCtx .
func FatalfCtx(c context.Context, format string, args ...interface{}) {
	logEntry := withContext(logrus.WithFields(logrus.Fields{}), c)
	logEntry.Fatalf(format, args...)
}

// FatalWithData .
func FatalWithData(message interface{}, data interface{}) {
	writeLog(SeverityCritical, message, data)
}

// FatalWithDataCtx .
func FatalWithDataCtx(c context.Context, message interface{}, data interface{}) {
	writeLogCtx(c, SeverityCritical, message, data)
}

// Define logrus alias
var (
	Debugf     = logrus.Debugf
	Infof      = logrus.Infof
	Warnf      = logrus.Warnf
	Errorf     = logrus.Errorf
	Fatalf     = logrus.Fatalf
	Panicf     = logrus.Panicf
	Printf     = logrus.Printf
	Info       = logrus.Info
	Debug      = logrus.Debug
	Error      = logrus.Error
	Warningf   = logrus.Warningf
	Warn       = logrus.Warn
	Warning    = logrus.Warning
	WithFields = logrus.WithFields
	Fatal      = logrus.Fatal
)
