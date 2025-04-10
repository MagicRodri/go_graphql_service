package logging

import (
	"os"
	"runtime"

	"github.com/MagicRodri/go_graphql_service/internal/config"
	"github.com/evalphobia/logrus_sentry"
	"github.com/sirupsen/logrus"
)

// Logger - основной интерфейс логгера (по-дефолту ему полностью соотвествует logrus)
type Logger interface {
	Tracef(format string, args ...interface{})
	Debugf(format string, args ...interface{})
	Infof(format string, args ...interface{})
	Warningf(format string, args ...interface{})
	Errorf(format string, args ...interface{})
	Fatalf(format string, args ...interface{})
	Panicf(format string, args ...interface{})
	Trace(args ...interface{})
	Debug(args ...interface{})
	Info(args ...interface{})
	Warning(args ...interface{})
	Error(args ...interface{})
	Fatal(args ...interface{})
	Panic(args ...interface{})
	WithFields(fields logrus.Fields) *logrus.Entry
}

var (
	logger Logger
)

func Init() {
	l := logrus.New()
	l.Formatter = &logrus.TextFormatter{FullTimestamp: true}
	l.Out = os.Stdout

	if config.Get().Logging.RavenDSN != "" {
		hook, err := logrus_sentry.NewSentryHook(config.Get().Logging.RavenDSN, []logrus.Level{
			logrus.PanicLevel,
			logrus.FatalLevel,
		})
		if err != nil {
			l.Fatalf("Fail init raven hook: %s", err)
		} else {
			hook.StacktraceConfiguration.Enable = true
			l.Hooks.Add(hook)
		}
	}

	level, err := logrus.ParseLevel(config.Get().Logging.Level)

	if err != nil {
		l.Fatalf("Invalid log level: %s", err)
	}

	l.Level = level

	if config.Get().Logging.Path != "" {
		file, err := os.OpenFile(config.Get().Logging.Path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			l.Fatalf("Error open log file %s: %s", config.Get().Logging.Path, err)
		} else {
			l.Infof("Switch log stream to %s", config.Get().Logging.Path)
			l.Out = file
		}
	}

	logger = l
}

func Get() Logger {
	return logger
}

func Sentry() {
	err := recover()

	if err != nil {
		pc, _, _, _ := runtime.Caller(1)
		Get().Fatalf("%s: %v", runtime.FuncForPC(pc).Name(), err)
	}
}
