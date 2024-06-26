// config
package floger

import (
	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"runtime"
)

//--------------------
// LOG LEVEL
//--------------------

// Log levels to control the logging output.
const (
	LevelTrace = iota
	LevelDebug1
	LevelDebug2
	LevelDebug3
	LevelDebug4
	LevelDebug5
	LevelInfo
	LevelWarning
	LevelError
	LevelCritical

	LogFileSizeMB = 1024 * 1024
)

// logLevel controls the global log level used by the logger.
type FileLoger struct {
	level int //= LevelInfo

	// logger references the used application logger.
	fileLogger    *log.Logger //= log.New(os.Stdout, "", log.Ldate|log.Lmicroseconds)
	errFileLogger *log.Logger //= log.New(os.Stderr, "", log.Ldate|log.Lmicroseconds)
	consoleLogger *log.Logger

	logFileHandle  *os.File
	strLogFileName string

	errLogFileHandle  *os.File
	strErrLogFileName string

	enableConsole   bool //= false
	maxLogFileSize  int64
	isFileLoger     bool
	localCalldepth  int
	LogFilter       string
	regexpLogFilter *regexp.Regexp
}

var (
	defaultFileLoger = NewFileLoger("", LevelDebug2, 10)
)

const (
	FLOGER_CALL_DEPTH = 4
)

func NewFileLoger(fileName string, nLevel int, maxFileSize int64) *FileLoger {
	fLoger := &FileLoger{}
	fLoger.level = nLevel
	fLoger.maxLogFileSize = maxFileSize
	fLoger.enableConsole = false
	fLoger.fileLogger = log.New(os.Stdout, "", log.Ldate|log.Lmicroseconds|log.Lshortfile)
	fLoger.errFileLogger = log.New(os.Stderr, "", log.Ldate|log.Lmicroseconds|log.Lshortfile)
	fLoger.consoleLogger = log.New(os.Stdout, "", log.Ldate|log.Lmicroseconds|log.Lshortfile)
	fLoger.isFileLoger = false
	fLoger.localCalldepth = 3
	if len(fileName) > 0 {
		os.MkdirAll(path.Dir(fileName), 0755)
		fLoger.strErrLogFileName = fileName + ".err"
		fLoger.strLogFileName = fileName

		err := fLoger.initErrLogger()
		if err != nil {
			log.Printf("%v", err)
			return nil
		}

		err = fLoger.initLogger()
		if err != nil {
			log.Printf("%v", err)
			return nil
		}
	}

	fLoger.SetLevel(nLevel)

	return fLoger
}
func (fl *FileLoger) flushLogger() {
	if fl.logFileHandle != nil {
		fl.logFileHandle.Sync()
	}

	if fl.errLogFileHandle != nil {
		fl.errLogFileHandle.Sync()
	}
}

func (fl *FileLoger) initLogger() error {

	var err error
	fl.logFileHandle, err = os.OpenFile(fl.strLogFileName, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0644)
	if err != nil {
		return fmt.Errorf("Open Log file error:%v", err)
	}
	fl.fileLogger = log.New(fl.logFileHandle, "", log.Ldate|log.Lmicroseconds|log.Lshortfile)
	fl.isFileLoger = true
	return nil
}

func (fl *FileLoger) initErrLogger() error {

	var err error
	fl.errLogFileHandle, err = os.OpenFile(fl.strErrLogFileName, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0644)
	if err != nil {
		return fmt.Errorf("Open Log file error:%v", err)
	}
	fl.errFileLogger = log.New(fl.errLogFileHandle, "", log.Ldate|log.Lmicroseconds|log.Lshortfile)
	fl.isFileLoger = true
	return nil
}

func InitLogger(fileName string) error {
	defaultFileLoger = nil
	defaultFileLoger = NewFileLoger(fileName, LevelInfo, 10)
	return nil
}

func Logger() *FileLoger {
	return defaultFileLoger
}

func (fl *FileLoger) finitLogger() {
	fl.logFileHandle.Close()
	fl.fileLogger = log.New(os.Stdout, "", log.Ldate|log.Lmicroseconds|log.Lshortfile)
}

func (fl *FileLoger) finitErrLogger() {
	fl.errLogFileHandle.Close()
	fl.errFileLogger = log.New(os.Stderr, "", log.Ldate|log.Lmicroseconds|log.Lshortfile)
}

func FinitLogger() {
	defaultFileLoger.FinitLogger()
}

func (fl *FileLoger) FinitLogger() {
	fl.finitLogger()
	fl.finitErrLogger()
}

func (fl *FileLoger) checkMaxLogFileSize(maxSize int64) {
	if fl.getLogFileLength() > maxSize {
		fl.finitLogger()
		bakLogFileName := fmt.Sprintf("%s.bak", fl.strLogFileName)
		os.Remove(bakLogFileName)
		os.Rename(fl.strLogFileName, bakLogFileName)
		fl.initLogger()
	}
}

func (fl *FileLoger) checkMaxErrLogFileSize(maxSize int64) {
	if fl.getErrLogFileLength() > maxSize {
		fl.finitErrLogger()
		bakErrLogFileName := fmt.Sprintf("%s.bak", fl.strErrLogFileName)
		os.Remove(bakErrLogFileName)
		os.Rename(fl.strErrLogFileName, bakErrLogFileName)
		fl.initErrLogger()
	}
}

func (fl *FileLoger) CheckMaxFileSize(maxSize int64) {
	fl.flushLogger()
	fl.checkMaxLogFileSize(maxSize)
	fl.checkMaxErrLogFileSize(maxSize)
}

func CheckMaxFileSize(maxSize int64) {
	defaultFileLoger.CheckMaxFileSize(maxSize)
}

func (fl *FileLoger) getLogFileLength() int64 {
	f, err := os.OpenFile(fl.strLogFileName, os.O_RDONLY, 0666)
	if err != nil {
		return 0
	}

	defer f.Close()
	stat, err := f.Stat()
	if err != nil {
		return 0
	}
	return stat.Size()
}

func (fl *FileLoger) getErrLogFileLength() int64 {
	f, err := os.OpenFile(fl.strErrLogFileName, os.O_RDONLY, 0666)
	if err != nil {
		return 0
	}

	defer f.Close()
	stat, err := f.Stat()
	if err != nil {
		return 0
	}
	return stat.Size()
}

// LogLevel returns the global log level and can be used in
// own implementations of the logger interface.
func (fl *FileLoger) Level() int {
	return fl.level
}

func Level() int {
	return defaultFileLoger.Level()
}

// SetLogLevel sets the global log level used by the simple
// logger.
func (fl *FileLoger) SetLevel(l int) {
	fl.level = l
	if fl.level <= LevelDebug3 {
		fl.localCalldepth = 5
	} else if fl.level <= LevelDebug4 {
		fl.localCalldepth = 3
	} else if fl.level <= LevelDebug5 {
		fl.localCalldepth = 2
	} else {
		fl.localCalldepth = 1
	}
}

func SetLevel(l int) {
	defaultFileLoger.SetLevel(l)
}

func (fl *FileLoger) SetLogFilter(logFilter string) {
	if len(logFilter) == 0 || logFilter == "none" {
		fl.LogFilter = ""
		fl.regexpLogFilter = nil
	}
	fl.LogFilter = logFilter
	fl.regexpLogFilter = regexp.MustCompile(logFilter)
}

func SetLogFilter(logFilter string) {
	defaultFileLoger.SetLogFilter(logFilter)
}

func (fl *FileLoger) SetEnableConsole(enable bool) {
	fl.enableConsole = enable
}

func SetEnableConsole(enable bool) {
	defaultFileLoger.SetEnableConsole(enable)
}

// SetLogger sets a new logger.

func (fl *FileLoger) SetLogger(l *log.Logger) {
	fl.fileLogger = l
}

func SetLogger(l *log.Logger) {
	defaultFileLoger.SetLogger(l)
}

func (fl *FileLoger) SetErrLogger(l *log.Logger) {
	fl.errFileLogger = l
}

func SetErrLogger(l *log.Logger) {
	defaultFileLoger.SetErrLogger(l)
}

// Trace logs a message at trace level.
func Trace(v ...interface{}) {
	defaultFileLoger.Trace(v...)
}
func (fl *FileLoger) Trace(v ...interface{}) {
	if fl.level <= LevelTrace {
		fl.fileLoggerOutput(LevelTrace, FLOGER_CALL_DEPTH, fmt.Sprintf("[T] %+v\n", v))

		if fl.isFileLoger && fl.enableConsole {
			fl.consoleLoggerOutput(LevelTrace, FLOGER_CALL_DEPTH, fmt.Sprintf("\033[32m[T] %+v\033[0m\n", v))
		}
	}
}

// Debug logs a message at debug level.
func Debug(v ...interface{}) {
	defaultFileLoger.Debug(v...)
}
func (fl *FileLoger) Debug(v ...interface{}) {
	if fl.level <= LevelDebug1 {
		fl.fileLoggerOutput(LevelTrace, FLOGER_CALL_DEPTH, fmt.Sprintf("[D0] %+v\n", v))

		if fl.isFileLoger && fl.enableConsole {
			fl.consoleLoggerOutput(LevelTrace, FLOGER_CALL_DEPTH, fmt.Sprintf("\033[32m[D0] %+v\033[0m\n", v))
		}
	}
}
func Debug1(v ...interface{}) {
	defaultFileLoger.Debug1(v...)
}
func (fl *FileLoger) Debug1(v ...interface{}) {
	if fl.level <= LevelDebug1 {
		fl.fileLoggerOutput(LevelDebug1, FLOGER_CALL_DEPTH, fmt.Sprintf("[D1] %+v\n", v))

		if fl.isFileLoger && fl.enableConsole {
			fl.consoleLoggerOutput(LevelDebug1, FLOGER_CALL_DEPTH, fmt.Sprintf("\033[32m[D1] %+v\033[0m\n", v))
		}
	}
}
func Debug2(v ...interface{}) {
	defaultFileLoger.Debug2(v...)
}
func (fl *FileLoger) Debug2(v ...interface{}) {
	if fl.level <= LevelDebug2 {
		fl.fileLoggerOutput(LevelDebug2, FLOGER_CALL_DEPTH, fmt.Sprintf("[D2] %+v\n", v))

		if fl.isFileLoger && fl.enableConsole {
			fl.consoleLoggerOutput(LevelDebug2, FLOGER_CALL_DEPTH, fmt.Sprintf("\033[32m[D2] %+v\033[0m\n", v))
		}
	}
}
func Debug3(v ...interface{}) {
	defaultFileLoger.Debug3(v...)
}
func (fl *FileLoger) Debug3(v ...interface{}) {
	if fl.level <= LevelDebug3 {
		fl.fileLoggerOutput(LevelDebug3, FLOGER_CALL_DEPTH, fmt.Sprintf("[D3] %+v\n", v))

		if fl.isFileLoger && fl.enableConsole {
			fl.consoleLoggerOutput(LevelDebug3, FLOGER_CALL_DEPTH, fmt.Sprintf("\033[32m[D3] %+v\033[0m\n", v))
		}
	}
}
func Debug4(v ...interface{}) {
	defaultFileLoger.Debug4(v...)
}
func (fl *FileLoger) Debug4(v ...interface{}) {
	if fl.level <= LevelDebug4 {
		fl.fileLoggerOutput(LevelDebug4, FLOGER_CALL_DEPTH, fmt.Sprintf("[D4] %+v\n", v))

		if fl.isFileLoger && fl.enableConsole {
			fl.consoleLoggerOutput(LevelDebug4, FLOGER_CALL_DEPTH, fmt.Sprintf("[D4] %+v\n", v))
		}
	}
}
func Debug5(v ...interface{}) {
	defaultFileLoger.Debug5(v...)
}
func (fl *FileLoger) Debug5(v ...interface{}) {
	if fl.level <= LevelDebug5 {
		fl.fileLoggerOutput(LevelDebug5, FLOGER_CALL_DEPTH, fmt.Sprintf("[D5] %+v\n", v))

		if fl.isFileLoger && fl.enableConsole {
			fl.consoleLoggerOutput(LevelDebug5, FLOGER_CALL_DEPTH, fmt.Sprintf("[D5] %+v\n", v))
		}
	}
}

func Debug5f(sfmt string, v ...interface{}) {
	defaultFileLoger.Debug5f(sfmt, v...)
}

func (fl *FileLoger) Debug5f(sfmt string, v ...interface{}) {
	if fl.level <= LevelDebug5 {
		fl.fileLoggerOutput(LevelDebug5, FLOGER_CALL_DEPTH, fmt.Sprintf("[D5] %s\n", fmt.Sprintf(sfmt, v...)))

		if fl.isFileLoger && fl.enableConsole {
			fl.consoleLoggerOutput(LevelDebug5, FLOGER_CALL_DEPTH, fmt.Sprintf("\033[36m[D5] %s\033[0m\n", fmt.Sprintf(sfmt, v...)))
		}
	}
}

func Debug4f(sfmt string, v ...interface{}) {
	defaultFileLoger.Debug4f(sfmt, v...)
}

func (fl *FileLoger) Debug4f(sfmt string, v ...interface{}) {
	if fl.level <= LevelDebug4 {
		fl.fileLoggerOutput(LevelDebug4, FLOGER_CALL_DEPTH, fmt.Sprintf("[D4] %s\n", fmt.Sprintf(sfmt, v...)))

		if fl.isFileLoger && fl.enableConsole {
			fl.consoleLoggerOutput(LevelDebug4, FLOGER_CALL_DEPTH, fmt.Sprintf("\033[36m[D4] %s\033[0m\n", fmt.Sprintf(sfmt, v...)))
		}
	}
}

func Debug3f(sfmt string, v ...interface{}) {
	defaultFileLoger.Debug3f(sfmt, v...)
}

func (fl *FileLoger) Debug3f(sfmt string, v ...interface{}) {
	if fl.level <= LevelDebug3 {
		fl.fileLoggerOutput(LevelDebug3, FLOGER_CALL_DEPTH, fmt.Sprintf("[D3] %s\n", fmt.Sprintf(sfmt, v...)))

		if fl.isFileLoger && fl.enableConsole {
			fl.consoleLoggerOutput(LevelDebug3, FLOGER_CALL_DEPTH, fmt.Sprintf("\033[36m[D3] %s\033[0m\n", fmt.Sprintf(sfmt, v...)))
		}
	}
}

func Debug2f(sfmt string, v ...interface{}) {
	defaultFileLoger.Debug2f(sfmt, v...)
}

func (fl *FileLoger) Debug2f(sfmt string, v ...interface{}) {
	if fl.level <= LevelDebug2 {
		fl.fileLoggerOutput(LevelDebug2, FLOGER_CALL_DEPTH, fmt.Sprintf("[D2] %s\n", fmt.Sprintf(sfmt, v...)))

		if fl.isFileLoger && fl.enableConsole {
			fl.consoleLoggerOutput(LevelDebug2, FLOGER_CALL_DEPTH, fmt.Sprintf("\033[36m[D2] %s\033[0m\n", fmt.Sprintf(sfmt, v...)))
		}
	}
}

// Info logs a message at info level.
func Info(v ...interface{}) {
	defaultFileLoger.Info(v...)
}

func Infof(sfmt string, v ...interface{}) {
	defaultFileLoger.Infof(sfmt, v...)
}

func (fl *FileLoger) Print(v ...interface{}) {
	fl.Debug(v...)
}

func (fl *FileLoger) Printf(sfmt string, v ...interface{}) {
	fl.Debug5f(sfmt, v...)
}

func (fl *FileLoger) Log(v ...interface{}) {
	fl.Info(v...)
}

func (fl *FileLoger) Logf(sfmt string, v ...interface{}) {
	fl.Infof(sfmt, v...)
}

func (fl *FileLoger) Infof(sfmt string, v ...interface{}) {
	if fl.level <= LevelInfo {
		fl.fileLoggerOutput(LevelInfo, FLOGER_CALL_DEPTH, fmt.Sprintf("[I] %s\n", fmt.Sprintf(sfmt, v...)))

		if fl.isFileLoger && fl.enableConsole {
			fl.consoleLoggerOutput(LevelInfo, FLOGER_CALL_DEPTH, fmt.Sprintf("\033[36m[I] %s\033[0m\n", fmt.Sprintf(sfmt, v...)))
		}
	}
}

func (fl *FileLoger) Info(v ...interface{}) {
	if fl.level <= LevelInfo {
		fl.fileLoggerOutput(LevelInfo, FLOGER_CALL_DEPTH, fmt.Sprintf("[I] %+v\n", v))

		if fl.isFileLoger && fl.enableConsole {
			fl.consoleLoggerOutput(LevelInfo, FLOGER_CALL_DEPTH, fmt.Sprintf("\033[36m[I] %+v\033[0m\n", v))
		}
	}
}

// Warning logs a message at warning level.
func Warn(v ...interface{}) {
	defaultFileLoger.Warn(v...)
}
func Warnf(sfmt string, v ...interface{}) {
	defaultFileLoger.Warnf(sfmt, v...)
}

func (fl *FileLoger) Warnf(sfmt string, v ...interface{}) {
	if fl.level <= LevelWarning {
		fl.fileLoggerOutput(LevelWarning, FLOGER_CALL_DEPTH, fmt.Sprintf("[W] %s\n", fmt.Sprintf(sfmt, v...)))

		if fl.isFileLoger && fl.enableConsole {
			fl.consoleLoggerOutput(LevelWarning, FLOGER_CALL_DEPTH, fmt.Sprintf("\033[33m[W] %s\033[0m\n", fmt.Sprintf(sfmt, v...)))
		}
	}
}

func (fl *FileLoger) Warn(v ...interface{}) {
	if fl.level <= LevelWarning {
		fl.fileLoggerOutput(LevelWarning, FLOGER_CALL_DEPTH, fmt.Sprintf("[W] %+v\n", v))

		if fl.isFileLoger && fl.enableConsole {
			fl.consoleLoggerOutput(LevelWarning, FLOGER_CALL_DEPTH, fmt.Sprintf("\033[33m[W] %+v\033[0m\n", v))
		}
	}
}

// Error logs a message at error level.
func Error(v ...interface{}) {
	defaultFileLoger.Error(v...)
}

func Errorf(sfmt string, v ...interface{}) {
	defaultFileLoger.Errorf(sfmt, v...)
}

func PanicError(v ...interface{}) {
	defaultFileLoger.PanicError(v...)
}

func PanicErrorf(sfmt string, v ...interface{}) {
	defaultFileLoger.PanicErrorf(sfmt, v...)
}

func (fl *FileLoger) Errorf(sfmt string, v ...interface{}) {
	if fl.level <= LevelError {
		fl.fileLoggerOutput(LevelError, FLOGER_CALL_DEPTH, fmt.Sprintf("[E] %s\n", fmt.Sprintf(sfmt, v...)))
		fl.errFileLoggerOutput(FLOGER_CALL_DEPTH, fmt.Sprintf("[E] %s\n", fmt.Sprintf(sfmt, v...)))

		if fl.isFileLoger && fl.enableConsole {
			fl.consoleLoggerOutput(LevelError, FLOGER_CALL_DEPTH, fmt.Sprintf("\033[31m[E] %s\033[0m\n", fmt.Sprintf(sfmt, v...)))
		}
	}
}

func (fl *FileLoger) PanicErrorf(sfmt string, v ...interface{}) {
	if fl.level <= LevelError {
		fl.fileLoggerOutput(LevelError, FLOGER_CALL_DEPTH, fmt.Sprintf("[E] %s\n", fmt.Sprintf(sfmt, v...)))
		fl.errFileLoggerOutput(FLOGER_CALL_DEPTH, fmt.Sprintf("[E] %s\n", fmt.Sprintf(sfmt, v...)))

		if fl.isFileLoger && fl.enableConsole {
			fl.consoleLoggerOutput(LevelError, FLOGER_CALL_DEPTH, fmt.Sprintf("\033[31m[E] %s\033[0m\n", fmt.Sprintf(sfmt, v...)))
		}

		for i := 1; ; i += 1 {
			pc, file, line, ok := runtime.Caller(i)
			if !ok {
				break
			}
			fn := runtime.FuncForPC(pc)
			fl.fileLoggerOutput(LevelError, FLOGER_CALL_DEPTH, fmt.Sprintf("  %s:%s:%d\r\n", fn.Name(), file, line))
			fl.errFileLoggerOutput(FLOGER_CALL_DEPTH, fmt.Sprintf("  %s:%s:%d\r\n", fn.Name(), file, line))

			if fl.isFileLoger && fl.enableConsole {
				fl.consoleLoggerOutput(LevelError, FLOGER_CALL_DEPTH, fmt.Sprintf("\033[31m[E]  %s:%s:%d\033[0m\n", fn.Name(), file, line))
			}
		}
	}
}

func (fl *FileLoger) FilterLogString(mesgLevel int, s string) bool {
	if mesgLevel < LevelDebug5 && fl.regexpLogFilter != nil {
		matched := fl.regexpLogFilter.FindAllStringSubmatch(s, -1)
		if len(matched) == 0 {
			return true
		}
	}

	return false
}

func (fl *FileLoger) fileLoggerOutput(mesgLevel int, calldepth int, s string) {
	if fl.FilterLogString(mesgLevel, s) {
		return
	}
	caller := ""
	for i := 1; i < fl.localCalldepth; i++ {
		_, file, line, ok := runtime.Caller(calldepth + i)
		if !ok {
			break
			file = "???"
			line = 0
		}
		fileName := filepath.Base(file)
		if "asm_amd64.s" == fileName {
			break
		}
		caller = caller + ">" + fmt.Sprintf("%s:%d", fileName, line)
	}

	fl.fileLogger.Output(calldepth, caller+s)
}

func (fl *FileLoger) consoleLoggerOutput(mesgLevel int, calldepth int, s string) {
	if fl.FilterLogString(mesgLevel, s) {
		return
	}
	caller := ""
	for i := 1; i < fl.localCalldepth; i++ {
		_, file, line, ok := runtime.Caller(calldepth + i)

		if !ok {
			break
			file = "???"
			line = 0
			break
		}
		fileName := filepath.Base(file)
		if "asm_amd64.s" == fileName ||
			"main.go" == fileName {
			break
		}
		caller = caller + ">" + fmt.Sprintf("%s:%d", fileName, line)
	}

	fl.consoleLogger.Output(calldepth, caller+s)
}

func (fl *FileLoger) errFileLoggerOutput(calldepth int, s string) {
	caller := ""
	for i := 1; i < fl.localCalldepth; i++ {
		_, file, line, ok := runtime.Caller(calldepth + i)
		if !ok {
			break
			file = "???"
			line = 0
		}
		fileName := filepath.Base(file)
		if "asm_amd64.s" == fileName {
			break
		}
		caller = caller + ">" + fmt.Sprintf("%s:%d", fileName, line)
	}

	fl.errFileLogger.Output(calldepth, caller+s)
}

func (fl *FileLoger) Error(v ...interface{}) {
	if fl.level <= LevelError {
		fl.fileLoggerOutput(LevelError, FLOGER_CALL_DEPTH, fmt.Sprintf("[E] %+v\n", v))
		fl.errFileLoggerOutput(FLOGER_CALL_DEPTH, fmt.Sprintf("[E] %+v\n", v))

		if fl.isFileLoger && fl.enableConsole {
			fl.consoleLoggerOutput(LevelError, FLOGER_CALL_DEPTH, fmt.Sprintf("\033[31m[E] %+v\033[0m\n", v))
		}
	}
}

func (fl *FileLoger) PanicError(v ...interface{}) {
	if fl.level <= LevelError {
		fl.fileLoggerOutput(LevelError, FLOGER_CALL_DEPTH, fmt.Sprintf("[E] %+v\n", v))
		fl.errFileLoggerOutput(FLOGER_CALL_DEPTH, fmt.Sprintf("[E] %+v\n", v))

		if fl.isFileLoger && fl.enableConsole {
			fl.consoleLoggerOutput(LevelError, FLOGER_CALL_DEPTH, fmt.Sprintf("\033[31m[E] %+v\033[0m\n", v))
		}

		for i := 2; ; i += 1 {
			pc, file, line, ok := runtime.Caller(i)
			if !ok {
				break
			}
			fn := runtime.FuncForPC(pc)
			fl.fileLoggerOutput(LevelError, FLOGER_CALL_DEPTH, fmt.Sprintf("!  %s:%s:%d\r\n", fn.Name(), file, line))
			fl.errFileLoggerOutput(FLOGER_CALL_DEPTH, fmt.Sprintf("!  %s:%s:%d\r\n", fn.Name(), file, line))

			if fl.isFileLoger && fl.enableConsole {
				fl.consoleLoggerOutput(LevelError, FLOGER_CALL_DEPTH, fmt.Sprintf("!  %s:%s:%d\r\n", fn.Name(), file, line))
			}
		}
	}
}

// Critical logs a message at critical level.
func Critical(v ...interface{}) {
	defaultFileLoger.Critical(v...)
}

func Criticalf(sfmt string, v ...interface{}) {
	defaultFileLoger.Criticalf(sfmt, v...)
}

func (fl *FileLoger) Criticalf(sfmt string, v ...interface{}) {
	if fl.level <= LevelCritical {
		fl.fileLoggerOutput(LevelCritical, FLOGER_CALL_DEPTH, fmt.Sprintf("[C] %s\n", fmt.Sprintf(sfmt, v...)))
		fl.errFileLoggerOutput(FLOGER_CALL_DEPTH, fmt.Sprintf("[C] %s\n", fmt.Sprintf(sfmt, v...)))

		if fl.isFileLoger && fl.enableConsole {
			fl.consoleLoggerOutput(LevelCritical, FLOGER_CALL_DEPTH, fmt.Sprintf("[C] %s\n", fmt.Sprintf(sfmt, v...)))
		}
	}
}

func (fl *FileLoger) Critical(v ...interface{}) {
	if fl.level <= LevelCritical {
		fl.fileLoggerOutput(LevelCritical, FLOGER_CALL_DEPTH, fmt.Sprintf("[C] %+v\n", v))
		fl.errFileLoggerOutput(FLOGER_CALL_DEPTH, fmt.Sprintf("[C] %+v\n", v))

		if fl.isFileLoger && fl.enableConsole {
			fl.consoleLoggerOutput(LevelCritical, FLOGER_CALL_DEPTH, fmt.Sprintf("[C] %+v\n", v))
		}
	}
}
