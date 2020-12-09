package printer

import (
	"fmt"

	"github.com/fatih/color"
)

var (
	regStr              = ""
	infoStr             = "✔  INFO:"
	warnStr             = "❗ WARN:"
	errStr              = "✗  ERROR:"
	fmtStr              = "%s %s\n"
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
)

// Printf print an informational message
func Printf(s string, args ...interface{}) {
	out := fmt.Sprintf(s, args...)
	fmt.Printf(fmtStr, regColorPrefix(regStr), regColorMsg(out))
}

// Println print an informational message
func Println(s string) {
	fmt.Printf(fmtStr, regColorPrefix(regStr), regColorMsg(s))
}

// Noticef print an informational message
func Noticef(s string, args ...interface{}) {
	out := fmt.Sprintf(s, args...)
	fmt.Printf(fmtStr, noticeColorPrefix(infoStr), noticeColorMsg(out))
}

// Noticeln print an informational message
func Noticeln(s string) {
	fmt.Printf(fmtStr, noticeColorPrefix(infoStr), noticeColorMsg(s))
}

// NoticeHif print an informational message
func NoticeHif(s string, args ...interface{}) {
	out := fmt.Sprintf(s, args...)
	fmt.Printf(fmtStr, noticeHiColorPrefix(infoStr), noticeHiColorMsg(out))
}

// NoticeHiln print an informational message
func NoticeHiln(s string) {
	fmt.Printf(fmtStr, noticeHiColorPrefix(infoStr), noticeHiColorMsg(s))
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
