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
	LevelDebug = iota
	LevelInfo
	LevelWarn
	LevelError
	LevelOff
)

const (
	Ldate         = 1 << iota     //日期示例： 2009/01/23
	Ltime                         //时间示例: 01:23:23
	Lmicroseconds                 //毫秒示例: 01:23:23.123123
	Llongfile                     //绝对路径和行号: /a/b/c/d.go:23
	Lshortfile                    //文件和行号: d.go:23
	LUTC                          //日期时间转为0时区的
	LstdFlags     = Ldate | Ltime //Go提供的标准抬头信息
)

type Logger struct {
	dirname string
	format  string
	output  string
	level   int
	logger  *log.Logger
}

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

func (l *Logger) SetLevel(level int) {
	l.level = level
}

func (l *Logger) SetFlags(flag int) {
	l.logger.SetFlags(flag)
}

func (l *Logger) Debug(v ...interface{}) {
	if l.level <= LevelDebug {
		l.logger.SetPrefix("[Debug] ")
		l.writeLog("", v...)
	}
}

func (l *Logger) Info(v ...interface{}) {
	if l.level <= LevelInfo {
		l.logger.SetPrefix("[Info] ")
		l.writeLog("", v...)
	}
}

func (l *Logger) Warn(v ...interface{}) {
	if l.level <= LevelWarn {
		l.logger.SetPrefix("[Warn] ")
		l.writeLog("", v...)
	}
}

func (l *Logger) Error(v ...interface{}) {
	if l.level <= LevelError {
		l.logger.SetPrefix("[Error] ")
		l.writeLog("", v...)
	}
}

func (l *Logger) Debugf(format string, v ...interface{}) {
	if l.level <= LevelDebug {
		l.logger.SetPrefix("[Debug] ")
		l.writeLog(format, v...)
	}
}

func (l *Logger) Infof(format string, v ...interface{}) {
	if l.level <= LevelInfo {
		l.logger.SetPrefix("[Info] ")
		l.writeLog(format, v...)
	}
}

func (l *Logger) Warnf(format string, v ...interface{}) {
	if l.level <= LevelWarn {
		l.logger.SetPrefix("[Warn] ")
		l.writeLog(format, v...)
	}
}

func (l *Logger) Errorf(format string, v ...interface{}) {
	if l.level <= LevelError {
		l.logger.SetPrefix("[Error] ")
		l.writeLog(format, v...)
	}
}

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
