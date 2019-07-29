package app

import (
	"time"
	"web-game-translation-tools/module"
)

//输出日志封装
func (this *UIStruct) AppendLog(message string) {
	module.Log.Info(message)
	this.SetLogData(message)
}

//输出日志封装
func (this *UIStruct) AppendLogError(err error) {
	module.Log.Error(err)
	this.SetLogData(err.Error())
}

//输出日志到窗口
func (this *UIStruct) SetLogData(message string) {
	var err error
	data := LogType{
		CreateTime: time.Now().Format("2006-01-02_15:04:05"),
		Message: message,
	}
	this.logData = data.CreateTime + " " + data.Message + "\r\n" + this.logData
		this.logDataList = append(this.logDataList, data)
	err = this.inputLog.SetText(this.logData)
	if err != nil{
		module.Log.Error(err)
		return
	}
	if  len(this.logDataList) >= 100{
		newList := []LogType{}
		for _,v := range this.logDataList{
			newList = append(newList, v)
		}
		this.logDataList = newList
	}
	return
}