/*go**************************************************************************
 File    : Logs.go
 Time    : 2018/9/11 15:04
 Author  : yanue
 Desc    : 日志功能

 Copyright (c) Shenzhen BB Team.
**************************************************************************go*/

package common

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"os"
	"path/filepath"
	"runtime"
	"runtime/debug"
	"strings"

	"go-esport-account/common/logs"
)

func init() {
	if runtime.GOOS != "linux" {
		Logger.SetLogger("console", "")
	}
}

type CLogger struct {
	*logs.BeeLogger
}

func (this *CLogger) Start(filename string) {
	if len(filename) == 0 {
		filename = filepath.Base(os.Args[0])
		if filepath.Ext(filename) != "" {
			po := strings.LastIndex(filename, ".")
			filename = filename[:po]
		}
	}

	this.setFile(LogPath, filename)
	this.Info("start " + filename + " ...")
}

/*
 *@param deep可选参数,兼容无参调用.默认为4,定位调用者代码位置,业务层有再次包装可继续加大深度
 *@note 让日志记录文件名和行号
 */
func (this *CLogger) EnableFileLine(deep ...int) {
	this.EnableFuncCallDepth(true)
	deepLevel := 4
	if len(deep) > 0 {
		deepLevel = deep[0]
	}
	this.SetLogFuncCallDepth(deepLevel)
}

func (this *CLogger) setFile(path, filename string) {
	if _, err := os.Stat(path); err != nil {
		os.Mkdir(path, 0755)
	}

	format := `{"perm":"0666","maxsize":20000000,"maxDays":60,"daily":false,"filename":"%s/%s.log"}`
	config := fmt.Sprintf(format, path, filename)
	Logger.SetLogger("file", config)
}

var Logger *CLogger = &CLogger{BeeLogger: logs.NewLogger(10000)}

func Debug(format string, v ...interface{}) {
	Logger.Debug(format, v...)
}

func Info(format string, v ...interface{}) {
	Logger.Info(format, v...)
}

func Notice(format string, v ...interface{}) {
	Logger.Notice(format, v...)
}

func Warning(format string, v ...interface{}) {
	Logger.Warning(format, v...)
}

func Error(format string, v ...interface{}) {
	// 把stack打印出来，略掉前2行
	stackInfo := string(debug.Stack())
	slice1 := stackInfo[strings.Index(stackInfo, "\n")+1:]
	stackInfoOut := "\n" +
		strings.Repeat("*", 70) +
		"\n" +
		slice1[strings.Index(slice1, "\n")+1:] +
		strings.Repeat("*", 70)

	Logger.Error(format+stackInfoOut, v...)
}

func EchoDebug(c gin.Context, format string, v ...interface{}) {
	Debug(fmt.Sprintf("clt:%s, %s", c.Request.RemoteAddr, format), v...)
}

func EchoInfo(c gin.Context, format string, v ...interface{}) {
	Info(fmt.Sprintf("clt:%s, %s", c.Request.RemoteAddr, format), v...)
}

func EchoNotice(c gin.Context, format string, v ...interface{}) {
	Notice(fmt.Sprintf("clt:%s, %s", c.Request.RemoteAddr, format), v...)
}

func EchoWarning(c gin.Context, format string, v ...interface{}) {
	Warning(fmt.Sprintf("clt:%s, %s", c.Request.RemoteAddr, format), v...)
}

func EchoError(c gin.Context, format string, v ...interface{}) {
	Error(fmt.Sprintf("clt:%s, %s", c.Request.RemoteAddr, format), v...)
}
