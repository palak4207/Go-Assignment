// Package logger is a custom package that will abstract away log functionality
package logger

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

type level int

// If you leave logLevel without any default value then its value is 0 or TRACE
var logLevel level

// Do we even need these many levels? Whats the purpose? Think about it! Lets keep things simple
const (
	// TRACE level logging ...
	TRACE level = iota
	// DEBUG level logging ...
	DEBUG
	// INFO level logging ...
	INFO
	// WARN level logging ...
	WARN
	// ERROR level logging ...
	ERROR
	// FATAL level logging will stop execution
	FATAL
)

// SetLogLevel ...
func SetLogLevel(l level) {
	logLevel = l
}

// LType ...
type LType int

var logType LType

var logFlag = log.Ldate | log.Ltime

var (
	logger = log.New(os.Stdout, "LOG: ", logFlag)
)

func getLogWriter(l level, wc []io.WriteCloser) (out io.WriteCloser) {

	out = os.Stdout

	if len(wc) > 0 {
		out = wc[0]
	}

	switch l {
	case ERROR:
		if len(wc) == 0 {
			out = os.Stderr
		}
	case FATAL:
		if len(wc) == 0 {
			out = os.Stderr
		}
	}

	return
}

// New logger
func New(l level, wc ...io.WriteCloser) *log.Logger {
	SetLogLevel(l)

	out := getLogWriter(l, wc)

	fn := func() {
		if out != nil {
			out.Close()
		}
	}
	cleanUpFuncs = append(cleanUpFuncs, fn)

	logger = log.New(out, "LOG", logFlag)

	return logger
}

// location where our logs will be stored
var logFile string

type cleanUpFn func()

var cleanUpFuncs = make([]cleanUpFn, 0)

// SetLogType ..
func SetLogType(logT LType) {
	logType = logT
}

// CleanUp will try to clean up any used resources for logging like files, directories etc
// You must call this when you done setting up the log
func CleanUp() {
	for i, fn := range cleanUpFuncs {
		D("Cleaning up >> ", i)
		fn()
	}
}

// SetFileLogger will set path where logs will get stored
func SetFileLogger(fileName string, l level, fileDir ...string) {
	var dir string

	if len(fileDir) == 0 {
		T("No fileDir provided")
		cwd, err := os.Getwd()
		if err != nil {
			log.Fatal(">> ", err)
		}
		dir = filepath.Join(cwd, "logs")
		if err = os.Mkdir(dir, 0777); os.IsNotExist(err) {
			E("Err in mkdir")
			log.Fatal(err)
		}
	} else {
		dir = fileDir[0]
	}

	logFile = filepath.Join(dir, fileName)
	T("Dir >> ", logFile)

	// WSL // file permissions
	// Windows: perm has no effect
	// 1 1 0 >> R W X >> 6 6 0
	wc, err := os.OpenFile(logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)

	if err != nil {
		P(err)
		os.Exit(2)
	}

	logger = New(l, wc)
}

// T is a trace logger
func T(v ...interface{}) {
	if int(logLevel) <= int(TRACE) {
		logger.SetPrefix("TRACE: ")
		printStackTrace(2, v...)
	}
}

// D is a debug logger
func D(v ...interface{}) {
	if int(logLevel) <= int(DEBUG) {
		logger.SetPrefix("DEBUG: ")
		printStackTrace(2, v...)
	}
}

// I is a Info logger
func I(v ...interface{}) {
	if int(logLevel) <= int(INFO) {
		logger.SetPrefix("INFO: ")
		printStackTrace(2, v...)
	}
}

// W is a warn logger
func W(v ...interface{}) {
	if int(logLevel) <= int(WARN) {
		logger.SetPrefix("WARNING: ")
		printStackTrace(2, v...)
	}
}

// E is an error logger
func E(v ...interface{}) {
	if int(logLevel) <= int(ERROR) {
		logger.SetPrefix("ERROR: ")
		printStackTrace(10, v...)
	}
}

// F is a fatal logger
func F(v ...interface{}) {
	if int(logLevel) <= int(FATAL) {
		logger.SetPrefix("FATAL: ")
		printStackTrace(10, v...)
		os.Exit(99)
	}
}

// P will log Println
func P(v ...interface{}) {
	log.Println(v...)
}

// PF will log Printf
func PF(format string, v ...interface{}) {
	log.Printf(format, v...)
}

// SL is log with stack trace
func SL(v ...interface{}) {
	printStackTrace(10)
	log.Println(v...)
}

func printStackTrace(maxStackLength int, v ...interface{}) {
	stackBuf := make([]uintptr, maxStackLength)
	length := runtime.Callers(3, stackBuf[:])
	stack := stackBuf[:length]

	trace := ""
	frames := runtime.CallersFrames(stack)
	// debug.PrintStack()
	for {
		frame, more := frames.Next()
		if !strings.Contains(frame.File, "runtime/") {
			trace = trace + fmt.Sprintf("\n:: %s:%d :: %s", frame.File, frame.Line, frame.Function)
		}

		if !more {
			break
		}
	}
	logger.Println(trace, v)

}
