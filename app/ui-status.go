package app

import (
	. "github.com/lxn/walk/declarative"
	"web-game-translation-tools/module"
)

//修改状态栏内容
func (this *UIStruct) SetStatusTip(tip string) {
	module.Log.Info("tip: ", tip)
	this.statusTipContent = tip
	this.Win.StatusBarItems = []StatusBarItem{
		StatusBarItem{
			Text:        this.statusTipContent,
		},
	}
}