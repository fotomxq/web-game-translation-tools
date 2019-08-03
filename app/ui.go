package app

import (
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
	"time"
)

type UIStruct struct {
	//窗口
	Win MainWindow
	//标题
	Title string
	//状态栏内容
	statusTipContent string
	//日志组件
	inputLog *walk.TextEdit
	//状态栏结构体
	statusBar []StatusBarItem
	//总体日志数据
	logData string
	logDataList []LogType
}

//初始化UI组件
func (this *UIStruct) Init() error {
	var err error
	//初始化ui
	this.Win = MainWindow{
		Title:   this.Title,
		MinSize: Size{600,600},
		StatusBarItems: this.statusBar,
		Layout:  VBox{},
		Children: []Widget{
			LinkLabel{
				Text:   "日志 ",
			},
			TextEdit{AssignTo: &this.inputLog},
			HSplitter{
				MaxSize:Size{0, 20},
				Children: []Widget{
					PushButton{
						Text: "导入和分析",
						OnClicked: func() {
						},
					},
					PushButton{
						Text: "使用本地词库翻译",
						OnClicked: func() {
						},
					},
					PushButton{
						Text: "开始在线翻译",
						OnClicked: func() {
						},
					},
					PushButton{
						Text: "开始翻译",
						OnClicked: func() {
						},
					},
				},
			},
		},
	}
	//更新状态信息
	this.SetStatusTip("窗口启动完成...")
	go this.setConfig()
	//启动窗口
	_,err = this.Win.Run()
	//反馈
	return err
}

//延迟设定参数
func (this *UIStruct) setConfig() {
	time.Sleep(time.Second * 1)
	//装载基本配置
	this.AppendLog("配置装载完成.")
	//构建提示信息
	this.AppendLog("----------------------------------------------------------------------")
	this.AppendLog("您可以修改导入词条规则、正则识别规则等。")
	this.AppendLog("在config目录下，您可以修改任意文件，可对导入和导出进行高级定制。")
	this.AppendLog("---------------------------高级说明-------------------------------")
	this.AppendLog("----------------------------------------------------------------------")
	this.AppendLog("最后点击“导出”按钮，翻译的文本将生成新的游戏文件（也可以覆盖原始文件，但不建议这么做）。")
	this.AppendLog("根据分析出的excel，按照列头说明进行翻译。")
	this.AppendLog("然后点击“分析”，程序将自动分析词条，并导出excel到“导出路径”内。")
	this.AppendLog("请先将游戏数据，放到“导入路径”下，或重新设定导入路径。")
	this.AppendLog("---------------------------使用说明-------------------------------")
}