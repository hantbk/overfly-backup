package logger

import (
	"fmt"
	"log"
	"os"
)

var (
	logFlag = log.Ldate | log.Ltime | log.LUTC
	myLog   = log.New(os.Stdout, "", logFlag)
)

func init() {

}

// Print log
func Print(v ...interface{}) {
	myLog.Print(v...)
}

// Println log
func Println(v ...interface{}) {
	myLog.Println(v...)
}

// Debug log
func Debug(v ...interface{}) {
	myLog.Println("[debug]", fmt.Sprint(v...))
}

// Info log
func Info(v ...interface{}) {
	myLog.Println(v...)
}

// Warn log
func Warn(v ...interface{}) {
	myLog.Println("[warn]", fmt.Sprint(v...))
}

// Error log
func Error(v ...interface{}) {
	myLog.Println("[error]", fmt.Sprint(v...))
}
