package app

import (
	"encoding/json"
	"errors"
	"web-game-translation-tools/module"
)

type ConfigDataType struct {
	//debug
	Debug bool `json:"Debug"`
	//原始文件路径
	DataPath string `json:"DataPath"`
	//输出路径
	OutputPath string `json:"OutputPath"`
}

func setConfig() (ConfigDataType,error) {
	res := ConfigDataType{}
	//读取配置
	by, err := module.File.LoadFile("." + Sep + "config" + Sep + "config.json")
	if err != nil {
		return res,errors.New("cannot load config.json file, " + err.Error())
	}
	//解析数据
	err = json.Unmarshal(by, &res)
	if err != nil {
		return res,errors.New("cannot read config , error : " + err.Error())
	}
	return res,nil
}

//写入配置
func saveConfig(configData ConfigDataType) error {
	by,err := json.Marshal(configData)
	if err != nil{
		return errors.New("cannot save config.json file, " + err.Error())
	}
	return module.File.WriteFile("." + Sep + "config" + Sep + "config.json", by)
}