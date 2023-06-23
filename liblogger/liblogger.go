package liblogger

import (
	"fmt"

	"github.com/helloferdie/golib/libtime"

	"io"
	"os"
	"runtime"
	"strconv"
	"strings"

	"github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
)

var initialize = false
var lg = logrus.New()

// loadConfig -
func loadConfig() {
	if !initialize {
		f := os.Getenv("log_file")
		if f == "" {
			f = "log.log"
		}

		wr := &lumberjack.Logger{
			Filename:   os.Getenv("dir_log") + "/" + f,
			MaxSize:    500, // megabytes
			MaxBackups: 3,
			MaxAge:     28,   //days
			Compress:   true, // disabled by default
		}

		silent := os.Getenv("log_silent")
		if silent == "1" {
			lg.SetOutput(wr) // Log to disk only
		} else {
			lg.SetOutput(io.MultiWriter(os.Stdout, wr)) // Log to stdout & disk
		}
		lg.SetReportCaller(true)

		initialize = true
	}
}

// trace - Backtrace log
func trace(stack int) string {
	pc := make([]uintptr, 10) // at least 1 entry needed
	runtime.Callers(stack, pc)
	f := runtime.FuncForPC(pc[0])
	file, line := f.FileLine(pc[0])
	dir := os.Getenv("dir_root") + "/"
	file = strings.Replace(file, dir, "", 1)
	lineStr := strconv.Itoa(line)
	return file + ":" + lineStr + " " + f.Name() + "\n"
}

// Log - Log error
func Log(v map[string]interface{}, doTrace bool) *logrus.Entry {
	loadConfig()
	if v == nil {
		v = map[string]interface{}{}
	}

	v["at"] = libtime.NowToString()
	if doTrace {
		v["trace3"] = trace(3)
		v["trace4"] = trace(4)
	}
	return lg.WithFields(v)
}

// Printf - Print format with prepend app timestamp
func Printf(format string, a ...any) (n int, err error) {
	a = append([]interface{}{libtime.NowToString()}, a...)
	return fmt.Printf("%s "+format, a...)
}

// Println - Print format with prepend app timestamp with new line
func Println(format string, a ...any) (n int, err error) {
	return Printf(format+"\n", a...)
}
