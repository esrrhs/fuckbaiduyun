package main

import (
	"fmt"
	"github.com/therecipe/qt/widgets"
	"os"
	"path/filepath"
	"time"
)

func main() {
	widgets.NewQApplication(len(os.Args), os.Args)

	var window = widgets.NewQMainWindow(nil, 0)
	var centralWidget = widgets.NewQWidget(window, 0)

	echoGroup := widgets.NewQGroupBox2("", nil)

	inputLabel := widgets.NewQLabel2("输入：", nil, 0)
	input := widgets.NewQLineEdit(nil)
	input.SetText(filepath.Dir(os.Args[0]))

	outputLabel := widgets.NewQLabel2("输出：", nil, 0)
	output := widgets.NewQLineEdit(nil)

	inputButton := widgets.NewQPushButton2("选择", nil)
	outputButton := widgets.NewQPushButton2("选择", nil)

	inputButton.ConnectClicked(func(checked bool) {
		w := widgets.NewQFileDialog2(nil, "选择输入目录", "", "")
		w.SetFileMode(widgets.QFileDialog__DirectoryOnly)
		input.SetText(w.GetExistingDirectory(nil, "选择输入目录", "", 0))
	})
	outputButton.ConnectClicked(func(checked bool) {
		w := widgets.NewQFileDialog2(nil, "选择输出目录", "", "")
		w.SetFileMode(widgets.QFileDialog__DirectoryOnly)
		output.SetText(w.GetExistingDirectory(nil, "选择输出目录", "", 0))
	})

	passLabel := widgets.NewQLabel2("密码：", nil, 0)
	pass := widgets.NewQLineEdit(nil)
	pass.SetText("123456")

	swapButton := widgets.NewQPushButton2("交换", nil)
	swapButton.ConnectClicked(func(checked bool) {
		tmp := input.Text()
		input.SetText(output.Text())
		output.SetText(tmp)
	})

	split := widgets.NewQComboBox(nil)
	split.AddItems([]string{"1G", "4G", "10G"})
	do := widgets.NewQComboBox(nil)
	do.AddItems([]string{"加密", "解密"})
	fuckButton := widgets.NewQPushButton2("GO", nil)

	cur := widgets.NewQProgressBar(nil)

	fuckButton.ConnectClicked(func(checked bool) {
		defer func() {
			if err := recover(); err != nil {
				a := widgets.NewQMessageBox(nil)
				a.SetText(fmt.Sprintf("%v", err))
				a.Show()
				//os.Exit(1)
			}
		}()

		cur.SetValue(0)

		var jobtotal int32
		var jobdone int32

		dojob(&jobtotal, &jobdone, input.Text(), output.Text(), do.CurrentText() == "加密", pass.Text(),
			1024*1024)

		for jobtotal > jobdone {
			cur.SetValue(int(jobdone * 100 / jobtotal))
			time.Sleep(time.Duration(100) * time.Millisecond)
		}

		cur.SetValue(100)
	})

	var echoLayout = widgets.NewQGridLayout2()
	echoLayout.AddWidget(inputLabel, 0, 0, 0)
	echoLayout.AddWidget(input, 0, 1, 0)
	echoLayout.AddWidget(inputButton, 0, 2, 0)
	echoLayout.AddWidget(outputLabel, 1, 0, 0)
	echoLayout.AddWidget(output, 1, 1, 0)
	echoLayout.AddWidget(outputButton, 1, 2, 0)

	echoLayout.AddWidget(passLabel, 2, 0, 0)
	echoLayout.AddWidget(pass, 2, 1, 0)
	echoLayout.AddWidget(swapButton, 2, 2, 0)

	echoLayout.AddWidget(split, 3, 0, 0)
	echoLayout.AddWidget(do, 3, 1, 0)
	echoLayout.AddWidget(fuckButton, 3, 2, 0)
	echoLayout.AddWidget3(cur, 4, 0, 3, 3, 0)
	echoGroup.SetLayout(echoLayout)

	var layout = widgets.NewQGridLayout2()
	layout.AddWidget(echoGroup, 0, 0, 0)

	centralWidget.SetLayout(layout)
	window.SetCentralWidget(centralWidget)
	window.SetMinimumWidth(500)
	window.SetWindowTitle("fuck baiduyun")
	window.Show()

	widgets.QApplication_Exec()
}
