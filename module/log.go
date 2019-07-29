package module

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"os"
	"time"
)

type LogType struct {
	//日志存储目录
	LogDir string
	//debug模式
	DebugOn bool
	//gin日志对象
	GinLog *logrus.Logger
	ginFd  *os.File
	//全局日志对象
	Logger    *logrus.Logger
	globLogFd *os.File
}

var (
	Log LogType
)

//初始化日志结构
func (this *LogType) SetConfig(debugOn bool) {
	//设定debug
	this.DebugOn = debugOn
	//如果没启动log则重建目录
	this.LogDir = "." + Sep + "log"
	if this.DebugOn == false {
		//确保log目录存在
		err := File.CreateFolder(this.LogDir)
		if err != nil {
			fmt.Println("log set config, create folder is error, " + err.Error())
		}
	}
	//设定logger
	this.Logger = logrus.New()
	this.Logger.SetFormatter(&logrus.JSONFormatter{})
	//设定gin
	this.GinLog = logrus.New()
	this.GinLog.SetFormatter(&logrus.JSONFormatter{})
	//自动初始化
	this.SetData()
}

//设定1次数据
func (this *LogType) SetData() {
	var err error

	if this.DebugOn == true {
		//如果启动debug模式，则输出到控制台，而不是文件
		logrus.Info("log debug is on..")

		this.GinLog.SetOutput(os.Stdout)
		this.GinLog.SetNoLock()

		this.Logger.SetOutput(os.Stdout)
		this.Logger.SetNoLock()
	} else {
		//如果关闭，则输出到日志文件内
		ginFdSrc := this.LogDir + Sep + "gin." + time.Now().Format("2006010215") + ".log"
		this.ginFd, err = os.OpenFile(ginFdSrc, os.O_WRONLY|os.O_CREATE|os.O_APPEND, os.ModeAppend)
		if err != nil {
			fmt.Println("cannot open file by gin logger, ", err)
		} else {
			this.GinLog.SetOutput(this.ginFd)
			this.GinLog.SetNoLock()
		}

		logErrorPath := this.LogDir + Sep + "system." + time.Now().Format("2006010215") + ".log"
		this.globLogFd, err = os.OpenFile(logErrorPath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, os.ModeAppend)
		if err != nil {
			fmt.Println("cannot open file by system logger, ", err)
		} else {
			this.Logger.SetOutput(this.globLogFd)
			this.Logger.SetNoLock()
		}
	}
}

//自动维护
func (this *LogType) Run() {
	for {
		this.SetData()
		time.Sleep(time.Second * 1)
	}
}

//发送日志组合
func (this *LogType) Error(args ...interface{}) {
	this.Logger.Error(args)
}
func (this *LogType) Info(args ...interface{}) {
	this.Logger.Info(args)
}
func (this *LogType) Debug(args ...interface{}) {
	this.Logger.Debug(args)
}