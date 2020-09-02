package golog

import (
	"fmt"
	"log"
	"os"
	"path"
	"strings"
	"time"
)

const (
	//LevelDebug debug
	LevelDebug = iota
	//LevelInfo info
	LevelInfo
	//LevelWarn warn
	LevelWarn
	//LevelError error
	LevelError
	//LevelOff off
	//LevelOff
)

//Logger 结构体
type Logger struct {
	dirname     string
	format      string
	output      string
	level       int
	logger      *log.Logger
	logChan     chan string
	logChanStat bool
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
	logger := log.New(os.Stdout, "", log.LstdFlags|log.Lshortfile)
	return &Logger{
		dirname: dirname,
		format:  format,
		output:  output,
		level:   level,
		logger:  logger,
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

//SetFlags 标准抬头信息
func (l *Logger) SetFlags(flag int) {
	l.logger.SetFlags(flag)
}

//Debug debug日志输出
func (l *Logger) Debug(v ...interface{}) {
	if l.level <= LevelDebug {
		l.logger.SetPrefix("[Debug] ")
		l.writeLog("", v...)
	}
}

//Info info日志输出
func (l *Logger) Info(v ...interface{}) {
	if l.level <= LevelInfo {
		l.logger.SetPrefix("[Info] ")
		l.writeLog("", v...)
	}
}

//Warn warn日志输出
func (l *Logger) Warn(v ...interface{}) {
	if l.level <= LevelWarn {
		l.logger.SetPrefix("[Warn] ")
		l.writeLog("", v...)
	}
}

//Error error日志输出
func (l *Logger) Error(v ...interface{}) {
	if l.level <= LevelError {
		l.logger.SetPrefix("[Error] ")
		l.writeLog("", v...)
	}
}

//Debugf debug格式化日志输出
func (l *Logger) Debugf(format string, v ...interface{}) {
	if l.level <= LevelDebug {
		l.logger.SetPrefix("[Debug] ")
		l.writeLog(format, v...)
	}
}

//Infof info格式化日志输出
func (l *Logger) Infof(format string, v ...interface{}) {
	if l.level <= LevelInfo {
		l.logger.SetPrefix("[Info] ")
		l.writeLog(format, v...)
	}
}

//Warnf warn格式化日志输出
func (l *Logger) Warnf(format string, v ...interface{}) {
	if l.level <= LevelWarn {
		l.logger.SetPrefix("[Warn] ")
		l.writeLog(format, v...)
	}
}

//Errorf error格式化日志输出
func (l *Logger) Errorf(format string, v ...interface{}) {
	if l.level <= LevelError {
		l.logger.SetPrefix("[Error] ")
		l.writeLog(format, v...)
	}
}

//writeLog 写日志操作
func (l *Logger) writeLog(format string, v ...interface{}) {
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
		localFile := strings.TrimRight(l.dirname, "/") + "/" + time.Now().Format(l.getFormat()) + ".log"
		localDir := path.Dir(localFile)
		err := os.MkdirAll(localDir, 0666)
		if err != nil {
			log.Printf("%v\n", err)
			return
		}
		file, err := os.OpenFile(localFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			log.Printf("%v\n", err)
			return
		}
		defer file.Close()
		l.logger.SetOutput(file)
	}
	_ = l.logger.Output(3, o)
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
