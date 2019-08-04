package app

import (
	"os"
	"web-game-translation-tools/module"
)

var(
	//UI
	ui UIStruct
	//文件路径分隔符
	Sep = string(os.PathSeparator)
	//全局配置
	configData ConfigDataType
)

//启动APP
func App(){
	var err error
	//读取配置信息
	configData, err = setConfig()
	if err != nil{
		module.Log.Error(err)
		return
	}
	//启动日志
	module.Log.SetConfig(configData.Debug)
	//设定Config
	ui.Title = "FTM翻译工具"
	//启动窗口
	err = ui.Init()
	if err != nil{
		module.Log.Error(err)
		return
	}
}