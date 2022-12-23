package golog

import (
	"fmt"
	"log"
	"net"
	"os"
	"path"
	"strings"
	"time"
)

const (
	LevelDebug = iota
	LevelInfo
	LevelWarn
	LevelError
)

var (
	logDebug *log.Logger
	logInfo  *log.Logger
	logWarn  *log.Logger
	logError *log.Logger
)

//Logger 结构体
type Logger struct {
	dirname     string
	format      string
	output      string
	level       int
	logFile     *os.File
	logName     string
	logChan     chan string
	logChanStat bool
}

func init() {
	logDebug = log.New(os.Stdout, "[Debug] ", log.LstdFlags|log.Lshortfile)
	logInfo = log.New(os.Stdout, "[Info] ", log.LstdFlags|log.Lshortfile)
	logWarn = log.New(os.Stdout, "[Warn] ", log.LstdFlags|log.Lshortfile)
	logError = log.New(os.Stdout, "[Error] ", log.LstdFlags|log.Lshortfile)
}

//NewLogger 实例化
func NewLogger(options map[string]interface{}) *Logger {
	var dirname = "."
	if v, ok := options["dirname"].(string); ok && v != "" {
		dirname = v
	}
	var format string
	if v, ok := options["format"].(string); ok && v != "" {
		format = v
	}
	var output string
	if v, ok := options["output"].(string); ok && v != "" {
		output = v
	}
	var level int
	if v, ok := options["level"].(int); ok && v > 0 {
		level = v
	}
	return &Logger{
		dirname: dirname,
		format:  format,
		output:  output,
		level:   level,
	}
}

//SetLevel 设置级别
func (l *Logger) SetLevel(level int) {
	l.level = level
}

//SetChan 设置log输出通道
func (l *Logger) SetChan(c chan string) {
	l.logChanStat = true
	l.logChan = c
}

//CancelChan 取消log输出通道
func (l *Logger) CancelChan() {
	l.logChanStat = false
}

//Debug debug日志输出
func (l *Logger) Debug(v ...interface{}) {
	if l.level <= LevelDebug {
		l.writeLog(logDebug, "", v...)
	}
}

//Info info日志输出
func (l *Logger) Info(v ...interface{}) {
	if l.level <= LevelInfo {
		l.writeLog(logInfo, "", v...)
	}
}

//Warn warn日志输出
func (l *Logger) Warn(v ...interface{}) {
	if l.level <= LevelWarn {
		l.writeLog(logWarn, "", v...)
	}
}

//Error error日志输出
func (l *Logger) Error(v ...interface{}) {
	if l.level <= LevelError {
		l.writeLog(logError, "", v...)
	}
}

//Debugf debug格式化日志输出
func (l *Logger) Debugf(format string, v ...interface{}) {
	if l.level <= LevelDebug {
		l.writeLog(logDebug, format, v...)
	}
}

//Infof info格式化日志输出
func (l *Logger) Infof(format string, v ...interface{}) {
	if l.level <= LevelInfo {
		l.writeLog(logInfo, format, v...)
	}
}

//Warnf warn格式化日志输出
func (l *Logger) Warnf(format string, v ...interface{}) {
	if l.level <= LevelWarn {
		l.writeLog(logWarn, format, v...)
	}
}

//Errorf error格式化日志输出
func (l *Logger) Errorf(format string, v ...interface{}) {
	if l.level <= LevelError {
		l.writeLog(logError, format, v...)
	}
}

//writeLog 写日志操作
func (l *Logger) writeLog(logger *log.Logger, format string, v ...interface{}) {
	o := ""
	if format == "" {
		o = fmt.Sprintln(v...)
	} else {
		o = fmt.Sprintf(format, v...)
	}
	go func() {
		if !l.logChanStat {
			return
		}
		l.logChan <- o
	}()
	if l.output == "file" {
		logName := strings.TrimRight(l.dirname, "/") + "/" + time.Now().Format(l.getFormat())
		var pid, ip string
		if strings.Contains(logName, "PID") {
			pid = fmt.Sprintf("%d", os.Getpid())
		}
		logName = strings.Replace(logName, "PID", pid, -1)
		if strings.Contains(logName, "IP") {
			ip = strings.Replace(getIpAddr(), ".", "_", -1)
		}
		logName = strings.Replace(logName, "IP", ip, -1)
		if logName != l.logName {
			localDir := path.Dir(logName)
			err := os.MkdirAll(localDir, 0666)
			if err != nil {
				log.Printf("%v\n", err)
				return
			}
			if l.logFile != nil {
				l.logFile.Close()
			}
			l.logFile, err = os.OpenFile(logName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
			if err != nil {
				log.Printf("%v\n", err)
				return
			}
			l.logName = logName
		}
		logger.SetOutput(l.logFile)
	}
	_ = logger.Output(3, o)
}

//getFormat 日志文件名格式化
func (l *Logger) getFormat() string {
	format := l.format
	format = strings.Replace(format, "yyyy", "2006", -1)
	format = strings.Replace(format, "MM", "01", -1)
	format = strings.Replace(format, "dd", "02", -1)
	format = strings.Replace(format, "HH", "15", -1)
	format = strings.Replace(format, "mm", "04", -1)
	format = strings.Replace(format, "ss", "05", -1)
	return format
}

func getIpAddr() string {
	//通过UDP方式获取
	local := "127.0.0.1"
	conn, err := net.Dial("udp", "8.8.8.8:8")
	if err != nil {
		return local
	}
	defer conn.Close()
	if err != nil {
		return local
	}
	return conn.LocalAddr().(*net.UDPAddr).IP.String()
}
