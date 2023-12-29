package log

import "log"

type logLevel string

const (
	InfoLevel    logLevel = "\033[32m[INFO]\033[0m"
	WarningLevel logLevel = "\033[33m[WARNING]\033[0m"
	ErrorLevel   logLevel = "\033[31m[ERROR]\033[0m"
)

func init() {
	log.SetFlags(log.LstdFlags | log.Lmsgprefix | log.LUTC)
}

func Infoln(v ...any) {
	Println(string(InfoLevel), v...)
}

func Info(v ...any) {
	Print(string(InfoLevel), v...)
}

func Infof(format string, v ...any) {
	Printf(string(InfoLevel), format, v...)
}

func Warningln(v ...any) {
	Println(string(WarningLevel), v...)
}

func Warning(v ...any) {
	Print(string(WarningLevel), v...)
}

func Warningf(format string, v ...any) {
	Printf(string(WarningLevel), format, v...)
}

func Errorln(v ...any) {
	Println(string(ErrorLevel), v...)
}

func Error(v ...any) {
	Print(string(ErrorLevel), v...)
}

func Errorf(format string, v ...any) {
	Printf(string(ErrorLevel), format, v...)
}

func Println(prefix string, v ...any) {
	log.SetPrefix(prefix + " ")
	log.Println(v...)
	log.SetPrefix("")
}

func Print(prefix string, v ...any) {
	log.SetPrefix(prefix + " ")
	log.Print(v...)
	log.SetPrefix("")
}

func Printf(prefix string, format string, v ...any) {
	log.SetPrefix(prefix + " ")
	log.Printf(format, v...)
	log.SetPrefix("")
}

func Fatalln(prefix string, v ...any) {
	log.SetPrefix(prefix + " ")
	log.Fatalln(v...)
	log.SetPrefix("")
}
