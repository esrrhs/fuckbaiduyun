//source: http://doc.qt.io/qt-5/qtwidgets-widgets-lineedits-example.html

package main

import (
	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/gui"
	"github.com/therecipe/qt/widgets"
	"os"
	"path/filepath"
)

func main() {
	widgets.NewQApplication(len(os.Args), os.Args)

	var window = widgets.NewQMainWindow(nil, 0)
	var centralWidget = widgets.NewQWidget(window, 0)

	echoGroup := widgets.NewQGroupBox2("", nil)

	inputLabel := widgets.NewQLabel2("输入", nil, 0)
	input := widgets.NewQLineEdit(nil)
	input.SetText(filepath.Dir(os.Args[0]))

	outputLabel := widgets.NewQLabel2("输出", nil, 0)
	output := widgets.NewQLineEdit(nil)

	inputButton := widgets.NewQPushButton2("选择", nil)
	outputButton := widgets.NewQPushButton2("选择", nil)

	inputButton.ConnectClicked(func(checked bool) {
		w := widgets.NewQFileDialog2(nil, "选择输入目录", "", "")
		w.SetFileMode(widgets.QFileDialog__DirectoryOnly)
		input.SetText(w.GetExistingDirectory(nil, "选择输入目录", "", 0))
	})
	outputButton.ConnectClicked(func(checked bool) {
		w := widgets.NewQFileDialog2(nil, "选择输入目录", "", "")
		w.SetFileMode(widgets.QFileDialog__DirectoryOnly)
		output.SetText(w.GetExistingDirectory(nil, "选择输入目录", "", 0))
	})

	split := widgets.NewQComboBox(nil)
	split.AddItems([]string{"1G", "4G", "10G"})
	do := widgets.NewQComboBox(nil)
	do.AddItems([]string{"加密", "解密"})
	fuckButton := widgets.NewQPushButton2("GO", nil)

	cur := widgets.NewQProgressBar(nil)
	cur.SetValue(80)

	var echoLayout = widgets.NewQGridLayout2()
	echoLayout.AddWidget(inputLabel, 0, 0, 0)
	echoLayout.AddWidget(input, 0, 1, 0)
	echoLayout.AddWidget(inputButton, 0, 2, 0)
	echoLayout.AddWidget(outputLabel, 1, 0, 0)
	echoLayout.AddWidget(output, 1, 1, 0)
	echoLayout.AddWidget(outputButton, 1, 2, 0)
	echoLayout.AddWidget(split, 2, 0, 0)
	echoLayout.AddWidget(do, 2, 1, 0)
	echoLayout.AddWidget(fuckButton, 2, 2, 0)
	echoLayout.AddWidget3(cur, 3, 0, 3, 3, 0)
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

func echoChanged(echoLineEdit *widgets.QLineEdit, index int) {
	switch index {
	case 0:
		{
			echoLineEdit.SetEchoMode(widgets.QLineEdit__Normal)
		}

	case 1:
		{
			echoLineEdit.SetEchoMode(widgets.QLineEdit__Password)
		}

	case 2:
		{
			echoLineEdit.SetEchoMode(widgets.QLineEdit__PasswordEchoOnEdit)
		}

	case 3:
		{
			echoLineEdit.SetEchoMode(widgets.QLineEdit__NoEcho)
		}
	}
}

func validatorChanged(validatorLineEdit *widgets.QLineEdit, index int) {
	switch index {
	case 0:
		{
			validatorLineEdit.SetValidator(nil)
		}

	case 1:
		{
			validatorLineEdit.SetValidator(gui.NewQIntValidator(validatorLineEdit))
		}

	case 2:
		{
			validatorLineEdit.SetValidator(gui.NewQDoubleValidator2(-999.0, 999.0, 2, validatorLineEdit))
		}
	}

	validatorLineEdit.Clear()
}

func alignmentChanged(alignmentLineEdit *widgets.QLineEdit, index int) {
	switch index {
	case 0:
		{
			alignmentLineEdit.SetAlignment(core.Qt__AlignLeft)
		}

	case 1:
		{
			alignmentLineEdit.SetAlignment(core.Qt__AlignCenter)
		}

	case 2:
		{
			alignmentLineEdit.SetAlignment(core.Qt__AlignRight)
		}
	}
}

func inputMaskChanged(inputMaskLineEdit *widgets.QLineEdit, index int) {
	switch index {
	case 0:
		{
			inputMaskLineEdit.SetInputMask("")
		}

	case 1:
		{
			inputMaskLineEdit.SetInputMask("+99 99 99 99 99;_")
		}

	case 2:
		{
			inputMaskLineEdit.SetInputMask("0000-00-00")
			inputMaskLineEdit.SetText("00000000")
			inputMaskLineEdit.SetCursorPosition(0)
		}

	case 3:
		{
			inputMaskLineEdit.SetInputMask(">AAAAA-AAAAA-AAAAA-AAAAA-AAAAA;#")
		}
	}
}

func accessChanged(accessLineEdit *widgets.QLineEdit, index int) {
	switch index {
	case 0:
		{
			accessLineEdit.SetReadOnly(false)
		}

	case 1:
		{
			accessLineEdit.SetReadOnly(true)
		}
	}
}
