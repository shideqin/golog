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
	LevelOff
)

const (
	//Ldate 日期示例：2009/01/23
	Ldate = 1 << iota
	//Ltime 时间示例: 01:23:23
	Ltime
	//Lmicroseconds 毫秒示例: 01:23:23.123123
	Lmicroseconds
	//Llongfile 绝对路径和行号: /a/b/c/d.go:23
	Llongfile
	//Lshortfile 文件和行号: d.go:23
	Lshortfile
	//LUTC 日期时间转为0时区的
	LUTC
	//LstdFlags Go提供的标准抬头信息
	LstdFlags = Ldate | Ltime
)

//Logger 结构体
type Logger struct {
	dirname string
	format  string
	output  string
	level   int
	logger  *log.Logger
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
	if l.output == "file" {
		localfile := strings.TrimRight(l.dirname, "/") + "/" + time.Now().Format(l.getFormat()) + ".log"
		localdir := path.Dir(localfile)
		err := os.MkdirAll(localdir, 0666)
		if err != nil {
			log.Printf("%v\n", err)
			return
		}
		file, err := os.OpenFile(localfile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		defer file.Close()
		if err != nil {
			log.Printf("%v\n", err)
			return
		}
		l.logger.SetOutput(file)
	}
	o := ""
	if format == "" {
		o = fmt.Sprintln(v...)
	} else {
		o = fmt.Sprintf(format, v...)
	}
	l.logger.Output(3, o)
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
