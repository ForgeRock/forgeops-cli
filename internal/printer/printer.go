package printer

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/ForgeRock/forgeops-cli/api"
)

var (
	regStr              = ""
	infoStr             = "✔  INFO:"
	warnStr             = "❗ WARN:"
	errStr              = "✗  ERROR:"
	fmtStr              = "%s %s"
	regColorPrefix      = color.New(color.Bold, color.FgWhite).SprintFunc()
	regColorMsg         = color.New(color.FgWhite).SprintFunc()
	noticeColorPrefix   = color.New(color.Bold, color.FgCyan).SprintFunc()
	noticeColorMsg      = color.New(color.FgCyan).SprintFunc()
	noticeHiColorPrefix = color.New(color.Bold, color.FgHiGreen).SprintFunc()
	noticeHiColorMsg    = color.New(color.FgHiGreen).SprintFunc()
	errorColorPrefix    = color.New(color.Bold, color.FgRed).SprintFunc()
	errorColorMsg       = color.New(color.FgRed).SprintFunc()
	warnColorPrefix     = color.New(color.Bold, color.FgYellow).SprintFunc()
	warnColorMsg        = color.New(color.FgYellow).SprintFunc()
	console             = log.Logger
	cmdResultOut        = log.Logger
	logn                = log.Logger
)

// OutType determine if we should be logging with text or json
type OutType string

var (
	// OutText output as text
	OutText OutType = "Text"
	// OutJson output as json
	OutJson OutType = "Json"
	// CommandOut output config setting
	CommandOut OutType = OutText
)

// Printf print an informational message
func Printf(s string, args ...interface{}) {
	out := fmt.Sprintf(s, args...)
	console.Printf(fmtStr, regColorPrefix(regStr), regColorMsg(out))
}

// Println print an informational message
func Println(s string) {
	console.Printf(fmtStr, regColorPrefix(regStr), regColorMsg(s))
}

// Noticef print a notice message
func Noticef(s string, args ...interface{}) {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	out := fmt.Sprintf(s, args...)
	console.Printf(fmtStr, noticeColorPrefix(infoStr), noticeColorMsg(out))
}

// Noticeln print a notice message
func Noticeln(s string) {
	console.Printf(fmtStr, noticeColorPrefix(infoStr), noticeColorMsg(s))
}

// NoticeHif print a notice message
func NoticeHif(s string, args ...interface{}) {
	out := fmt.Sprintf(s, args...)
	console.Printf(fmtStr, noticeHiColorPrefix(infoStr), noticeHiColorMsg(out))
}

// NoticeHiln print a notice message
func NoticeHiln(s string) {
	console.Printf(fmtStr, noticeHiColorPrefix(infoStr), noticeHiColorMsg(s))
}

// Warnf print a warning message
func Warnf(s string, args ...interface{}) {
	out := fmt.Sprintf(s, args...)
	console.Printf(fmtStr, warnColorPrefix(warnStr), warnColorMsg(out))
}

// Warnln print an informational message
func Warnln(s string) {
	console.Printf(fmtStr, warnColorPrefix(warnStr), warnColorMsg(s))
}

// Errorf print an error message
func Errorf(s string, args ...interface{}) {
	out := fmt.Sprintf(s, args...)
	console.Printf(fmtStr, errorColorPrefix(errStr), errorColorMsg(out))
}

// Errorln print an error message
func Errorln(s string) {
	console.Printf(fmtStr, errorColorPrefix(errStr), errorColorMsg(s))
}

// JsonResult provide a message that contains api.ForgeOpsResult
// This logger is configured to ignore log levels set by a user because it's the
// "return value" for a commmand. This is the output method to be used for scripting interfaces
func JsonResult(msg string, res *api.ForgeOpsResult) {
	eventResult := zerolog.Dict()
	for _, msg := range res.Results {
		for k, v := range msg {
			eventResult.Str(k, v)
		}
	}
	cmdResultOut.Info().
		Str("version", res.Version).
		Str("status", string(res.Status)).
		Dict("results", eventResult).Msg(msg)
}

// Logger global main logger that should be used to log errors, debug etc
func Logger() zerolog.Logger {
	return logn
}

// InitLogn configures main logger for program..
// logn with "text" mode is always debug level due to zerolog
func InitLogn(logType OutType, l zerolog.Level) {
	if l == zerolog.Disabled {
		logn = zerolog.Nop()
		CommandOut = logType
		return
	}
	switch logType {
	case OutText:
		CommandOut = OutText
		consolen := zerolog.ConsoleWriter{Out: os.Stdout}
		logn = zerolog.New(consolen).With().Timestamp().Logger()
	case OutJson:
		CommandOut = OutJson
		logn = zerolog.New(os.Stdout).With().Timestamp().Logger()
	}
	logn.Level(l)
}

func init() {
	// Global logging config
	zerolog.TimeFieldFormat = time.RFC3339

	// Setup logging to console
	consoleOutput := zerolog.ConsoleWriter{Out: os.Stdout}
	consoleOutput.FormatLevel = func(i interface{}) string {
		return ""
	}
	consoleOutput.FormatMessage = func(i interface{}) string {
		return fmt.Sprintf("%s", i)
	}
	consoleOutput.FormatFieldName = func(i interface{}) string {
		return fmt.Sprintf("%s:", i)
	}
	consoleOutput.FormatFieldValue = func(i interface{}) string {
		return strings.ToUpper(fmt.Sprintf("%s", i))
	}
	consoleOutput.PartsOrder = []string{
		zerolog.MessageFieldName,
	}
	console = zerolog.New(consoleOutput)

	// Setup a command result log, this should be used to log the result of a command
	cmdResultOut = zerolog.New(os.Stdout).
		With().
		Timestamp().
		Logger()
}
