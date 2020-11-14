package printer

import (
	"fmt"

	"github.com/fatih/color"
)

var (
	infoStr           = "✔  INFO:"
	warnStr           = "❗ WARN:"
	errStr            = "✗  ERROR:"
	fmtStr            = "%s %s\n"
	noticeColorPrefix = color.New(color.Bold, color.FgCyan).SprintFunc()
	noticeColorMsg    = color.New(color.FgCyan).SprintFunc()

	errorColorPrefix = color.New(color.Bold, color.FgRed).SprintFunc()
	errorColorMsg    = color.New(color.FgRed).SprintFunc()

	warnColorPrefix = color.New(color.Bold, color.FgYellow).SprintFunc()
	warnColorMsg    = color.New(color.FgYellow).SprintFunc()
)

// Noticef print an informational message
func Noticef(s string, args ...interface{}) {
	out := fmt.Sprintf(s, args...)
	fmt.Printf(fmtStr, noticeColorPrefix(infoStr), noticeColorMsg(out))
}

// Noticeln print an informational message
func Noticeln(s string) {
	fmt.Printf(fmtStr, noticeColorPrefix(infoStr), noticeColorMsg(s))
}

// Warnf print an informational message
func Warnf(s string, args ...interface{}) {
	out := fmt.Sprintf(s, args...)
	fmt.Printf(fmtStr, warnColorPrefix(warnStr), warnColorMsg(out))
}

// Warnln print an informational message
func Warnln(s string) {
	fmt.Printf(fmtStr, warnColorPrefix(warnStr), warnColorMsg(s))
}

// Errorf print an informational message
func Errorf(s string, args ...interface{}) {
	out := fmt.Sprintf(s, args...)
	fmt.Printf(fmtStr, errorColorPrefix(errStr), errorColorMsg(out))
}

// Errorln print an informational message
func Errorln(s string) {
	fmt.Printf(fmtStr, errorColorPrefix(errStr), errorColorMsg(s))
}
