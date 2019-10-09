package main

import (
	"encoding/json"
	"fmt"
	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/widgets"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
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
		input.SetText(w.GetExistingDirectory(nil, "选择输入目录", "", 0) + "/")
	})
	outputButton.ConnectClicked(func(checked bool) {
		w := widgets.NewQFileDialog2(nil, "选择输出目录", "", "")
		w.SetFileMode(widgets.QFileDialog__DirectoryOnly)
		output.SetText(w.GetExistingDirectory(nil, "选择输出目录", "", 0) + "/")
	})

	passLabel := widgets.NewQLabel2("密码：", nil, 0)
	pass := widgets.NewQLineEdit(nil)
	pass.SetText("123456")

	split := widgets.NewQComboBox(nil)
	split.AddItems([]string{"1G", "4G", "10G", "20G"})
	do := widgets.NewQComboBox(nil)
	do.AddItems([]string{"加密", "解密"})
	fuckButton := widgets.NewQPushButton2("GO", nil)

	cur := widgets.NewQProgressBar(nil)
	curf := widgets.NewQProgressBar(nil)
	swapButton := widgets.NewQPushButton2("交换", nil)
	swapButton.ConnectClicked(func(checked bool) {
		tmp := input.Text()
		input.SetText(output.Text())
		output.SetText(tmp)
		if do.CurrentText() == "加密" {
			do.SetCurrentText("解密")
		} else {
			do.SetCurrentText("加密")
		}
	})

	fuckButton.ConnectClicked(func(checked bool) {

		errstr := new(string)
		jobtotal := new(int32)
		jobdone := new(int32)
		filetotal := new(int64)
		filedone := new(int64)

		input.SetDisabled(true)
		output.SetDisabled(true)
		pass.SetDisabled(true)
		inputButton.SetDisabled(true)
		outputButton.SetDisabled(true)
		split.SetDisabled(true)
		do.SetDisabled(true)
		fuckButton.SetDisabled(true)
		swapButton.SetDisabled(true)

		t := core.NewQTimer(nil)
		t.ConnectEvent(func(e *core.QEvent) bool {

			if len(*errstr) > 0 {
				t.DisconnectEvent()
				a := widgets.NewQMessageBox(nil)
				a.SetText(*errstr)
				a.Show()

				input.SetDisabled(false)
				output.SetDisabled(false)
				pass.SetDisabled(false)
				inputButton.SetDisabled(false)
				outputButton.SetDisabled(false)
				split.SetDisabled(false)
				do.SetDisabled(false)
				fuckButton.SetDisabled(false)
				swapButton.SetDisabled(false)

				return true
			}
			if *jobtotal != 0 {
				cur.SetValue(int(*jobdone * 100 / *jobtotal))
			}
			if *filetotal != 0 {
				curf.SetValue(int(*filedone * 100 / *filetotal))
			}
			return true

		})
		t.Start(100)

		go func() {

			defer func() {
				if err := recover(); err != nil {
					*errstr = fmt.Sprintf("%v", err)
				}
			}()

			gConfig.Split = split.CurrentText()
			gConfig.Do = do.CurrentText()
			gConfig.Input = input.Text()
			gConfig.Output = output.Text()
			gConfig.Pass = pass.Text()
			saveJson(gConfig)

			cur.SetValue(0)
			curf.SetValue(0)

			split, _ := strconv.Atoi(strings.TrimRight(split.CurrentText(), "G"))

			dojob(jobtotal, jobdone, filetotal, filedone, input.Text(), output.Text(), do.CurrentText() == "加密", pass.Text(),
				1000*1000*1000*split)

			cur.SetValue(100)
			curf.SetValue(100)

			*errstr = "ok"
		}()

	})

	var echoLayout = widgets.NewQGridLayout2()
	echoLayout.AddWidget2(inputLabel, 0, 0, 0)
	echoLayout.AddWidget2(input, 0, 1, 0)
	echoLayout.AddWidget2(inputButton, 0, 2, 0)
	echoLayout.AddWidget2(outputLabel, 1, 0, 0)
	echoLayout.AddWidget2(output, 1, 1, 0)
	echoLayout.AddWidget2(outputButton, 1, 2, 0)

	echoLayout.AddWidget2(passLabel, 2, 0, 0)
	echoLayout.AddWidget2(pass, 2, 1, 0)
	echoLayout.AddWidget2(swapButton, 2, 2, 0)

	echoLayout.AddWidget2(split, 3, 0, 0)
	echoLayout.AddWidget2(do, 3, 1, 0)
	echoLayout.AddWidget2(fuckButton, 3, 2, 0)
	echoLayout.AddWidget3(cur, 4, 0, 3, 3, 0)
	echoLayout.AddWidget3(curf, 7, 0, 3, 3, 0)
	echoGroup.SetLayout(echoLayout)

	var layout = widgets.NewQGridLayout2()
	layout.AddWidget2(echoGroup, 0, 0, 0)

	lg := loadJson()
	if lg != nil {
		gConfig = *lg
		do.SetCurrentText(gConfig.Do)
		split.SetCurrentText(gConfig.Split)
		input.SetText(gConfig.Input)
		output.SetText(gConfig.Output)
		pass.SetText(gConfig.Pass)
	}

	centralWidget.SetLayout(layout)
	window.SetCentralWidget(centralWidget)
	window.SetMinimumWidth(500)
	window.SetWindowTitle("fuck baiduyun")
	window.Show()

	widgets.QApplication_Exec()
}

type Config struct {
	Input  string `json:"Input"`
	Output string `json:"Output"`
	Do     string `json:"Do"`
	Split  string `json:"Split"`
	Pass   string `json:"Pass"`
}

var gConfig Config

func saveJson(c Config) {
	jsonFile, err := os.OpenFile(".fuckbaiduyun.json",
		os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		return
	}
	defer jsonFile.Close()

	str, err := json.Marshal(&c)
	if err != nil {
		return
	}
	jsonFile.Write(str)

}
func loadJson() *Config {
	jsonFile, err := os.Open(".fuckbaiduyun.json")
	if err != nil {
		return nil
	}
	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)
	var c Config

	err = json.Unmarshal(byteValue, &c)
	if err != nil {
		return nil
	}

	return &c
}
