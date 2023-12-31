package log

import "log"

type logLevel string

const (
	// InfoLevel indicates that the printed log, is a harmless info.
	InfoLevel logLevel = "\033[32m[INFO]\033[0m"
	// WarningLevel means things are getting heavier.
	WarningLevel logLevel = "\033[33m[WARNING]\033[0m"
	// ErrorLevel means that something really bad happened.
	ErrorLevel logLevel = "\033[31m[ERROR]\033[0m"
)

func init() {
	log.SetFlags(log.LstdFlags | log.Lmsgprefix | log.LUTC)
}

// Infoln prints an info log with a new line.
func Infoln(v ...any) {
	Println(string(InfoLevel), v...)
}

// Info prints an info log.
func Info(v ...any) {
	Print(string(InfoLevel), v...)
}

// Infof prints a formatted info log.
func Infof(format string, v ...any) {
	Printf(string(InfoLevel), format, v...)
}

// Warningln prints a warning log with a new line.
func Warningln(v ...any) {
	Println(string(WarningLevel), v...)
}

// Warning prints a warning log.
func Warning(v ...any) {
	Print(string(WarningLevel), v...)
}

// Warningf prints a formatted warning log.
func Warningf(format string, v ...any) {
	Printf(string(WarningLevel), format, v...)
}

// Errorln prints an error log with a new line.
func Errorln(v ...any) {
	Println(string(ErrorLevel), v...)
}

// Error prints an error log.
func Error(v ...any) {
	Print(string(ErrorLevel), v...)
}

// Errorf prints a formatted error log.
func Errorf(format string, v ...any) {
	Printf(string(ErrorLevel), format, v...)
}

// Println prints a log with a specific prefix with a new line.
func Println(prefix string, v ...any) {
	log.SetPrefix(prefix + " ")
	log.Println(v...)
	log.SetPrefix("")
}

// Print prints a log with a specific prefix.
func Print(prefix string, v ...any) {
	log.SetPrefix(prefix + " ")
	log.Print(v...)
	log.SetPrefix("")
}

// Printf prints a formatted log with a specific prefix.
func Printf(prefix string, format string, v ...any) {
	log.SetPrefix(prefix + " ")
	log.Printf(format, v...)
	log.SetPrefix("")
}

// Fatalln prints a log with a specific prefix with a new line,
// and terminates the application with an error code (1).
func Fatalln(prefix string, v ...any) {
	log.SetPrefix(prefix + " ")
	log.Fatalln(v...)
	log.SetPrefix("")
}
