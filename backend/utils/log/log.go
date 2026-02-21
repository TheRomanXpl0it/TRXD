package log

import (
	"os"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/log"
)

type Logger struct {
	*log.Logger
}

const NOTICE = log.Level(2)
const CRITICAL = log.Level(10)

const FatalLevel = log.FatalLevel
const CriticalLevel = CRITICAL
const ErrorLevel = log.ErrorLevel
const WarnLevel = log.WarnLevel
const NoticeLevel = NOTICE
const InfoLevel = log.InfoLevel
const DebugLevel = log.DebugLevel

var LOG *Logger = &Logger{log.Default()}

var Fatal = LOG.Fatal
var Fatalf = LOG.Fatalf
var Critical = LOG.Critical
var Criticalf = LOG.Criticalf
var Error = LOG.Error
var Errorf = LOG.Errorf
var Warn = LOG.Warn
var Warnf = LOG.Warnf
var Notice = LOG.Notice
var Noticef = LOG.Noticef
var Info = LOG.Info
var Infof = LOG.Infof
var Debug = LOG.Debug
var Debugf = LOG.Debugf

var With = LOG.With
var WithPrefix = LOG.WithPrefix
var GetPrefix = LOG.GetPrefix
var SetPrefix = LOG.SetPrefix
var SetReportCaller = LOG.SetReportCaller
var GetLevel = LOG.GetLevel
var SetLevel = LOG.SetLevel
var SetLogLevel = LOG.SetLogLevel
var SetTimeFormat = LOG.SetTimeFormat
var SetStyles = LOG.SetStyles
var DefaultStyles = LOG.DefaultStyles

var defaultStyles *log.Styles

func init() {
	err := os.Setenv("TERM", "xterm-256color")
	if err != nil {
		log.Error("Failed to set TERM environment variable", "err", err)
	}

	defaultStyles = log.DefaultStyles()
	defaultStyles.Levels[CRITICAL] = lipgloss.NewStyle().
		SetString("CRITICAL").
		Bold(true).
		MaxWidth(4).
		Foreground(lipgloss.Color("201"))
	defaultStyles.Levels[NOTICE] = lipgloss.NewStyle().
		SetString("NOTICE").
		Bold(true).
		MaxWidth(4).
		Foreground(lipgloss.Color("40"))

	defaultStyles.Keys["debug"] = lipgloss.NewStyle().Foreground(lipgloss.Color("63"))
	defaultStyles.Values["debug"] = lipgloss.NewStyle().Bold(true)
	defaultStyles.Keys["info"] = lipgloss.NewStyle().Foreground(lipgloss.Color("86"))
	defaultStyles.Values["info"] = lipgloss.NewStyle().Bold(true)
	defaultStyles.Keys["notice"] = lipgloss.NewStyle().Foreground(lipgloss.Color("40"))
	defaultStyles.Values["notice"] = lipgloss.NewStyle().Bold(true)
	defaultStyles.Keys["warn"] = lipgloss.NewStyle().Foreground(lipgloss.Color("192"))
	defaultStyles.Values["warn"] = lipgloss.NewStyle().Bold(true)
	defaultStyles.Keys["err"] = lipgloss.NewStyle().Foreground(lipgloss.Color("204"))
	defaultStyles.Values["err"] = lipgloss.NewStyle().Bold(true)
	defaultStyles.Keys["crit"] = lipgloss.NewStyle().Foreground(lipgloss.Color("201"))
	defaultStyles.Values["crit"] = lipgloss.NewStyle().Bold(true)
	defaultStyles.Keys["fatal"] = lipgloss.NewStyle().Foreground(lipgloss.Color("134"))
	defaultStyles.Values["fatal"] = lipgloss.NewStyle().Bold(true)

	LOG.SetStyles(defaultStyles)
	LOG.SetTimeFormat("2006-01-02 15:04:05")
}

func (l *Logger) DefaultStyles() *log.Styles {
	return defaultStyles
}

func (l *Logger) SetLogLevel(level string) {
	var lvl log.Level
	var err error

	level = strings.ToLower(level)
	switch level {
	case "notice":
		lvl = NoticeLevel
	case "critical":
		lvl = CriticalLevel
	default:
		lvl, err = log.ParseLevel(level)
		if err != nil {
			l.Errorf("Invalid log level '%s': %v", level, err)
			return
		}
	}
	l.SetLevel(lvl)
}

func (l *Logger) With(keyvals ...interface{}) *Logger {
	return &Logger{l.Logger.With(keyvals...)}
}

func (l *Logger) WithPrefix(prefix string) *Logger {
	return &Logger{l.Logger.WithPrefix(prefix)}
}

func (l *Logger) Notice(msg interface{}, keyvals ...interface{}) {
	l.Helper()
	l.Log(NoticeLevel, msg, keyvals...)
}

func (l *Logger) Noticef(format string, keyvals ...interface{}) {
	l.Helper()
	l.Logf(NoticeLevel, format, keyvals...)
}

func (l *Logger) Critical(msg interface{}, keyvals ...interface{}) {
	l.Helper()
	l.Log(CriticalLevel, msg, keyvals...)
}

func (l *Logger) Criticalf(format string, keyvals ...interface{}) {
	l.Helper()
	l.Logf(CriticalLevel, format, keyvals...)
}
