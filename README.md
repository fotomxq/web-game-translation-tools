# web游戏翻译工具

本脚本用于快速翻译web或任意文本结构为主的游戏，本工具主要针对于Twine/SugarCube工具开发的文字类游戏。当然，不出意外其他所有文本类游戏，都是支持的。

# 使用方法

# 代码部署和编译

1、安装第三方库

go get github.com/lxn/walk

go get github.com/akavel/rsrc

rsrc -manifest test.manifest -o rsrc.syso

2、编译exe文件

go build -ldflags="-H windowsgui"

# FAQ

# 未来计划

1、操作界面。

2、词汇自动拆分重组功能，避免短词汇翻译影响到长句翻译。

3、自动翻译，以及自动减少句子中穿插代码被一起翻译的问题。

4、自动提取文本。

# 协议

本项目采用Apache2.0协议。