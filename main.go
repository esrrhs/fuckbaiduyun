package main

import (
	"encoding/json"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

func main() {
	initLog()

	a := app.NewWithID("com.esrrhs.fuckbaiduyun")
	w := a.NewWindow("fuck baiduyun v" + Version)

	exeDir := "."
	if exe, err := os.Executable(); err == nil {
		exeDir = filepath.Dir(exe)
	}

	input := widget.NewEntry()
	input.SetText(exeDir)
	output := widget.NewEntry()
	pass := widget.NewEntry()
	pass.SetText("123456")

	split := widget.NewSelect([]string{"1G", "4G", "10G", "20G"}, nil)
	split.SetSelected("1G")
	doSel := widget.NewSelect([]string{"加密", "解密"}, nil)
	doSel.SetSelected("加密")

	cur := widget.NewProgressBar()
	curf := widget.NewProgressBar()

	pickDir := func(target *widget.Entry) {
		dialog.ShowFolderOpen(func(uri fyne.ListableURI, err error) {
			if err != nil {
				dialog.ShowError(err, w)
				return
			}
			if uri == nil {
				return
			}
			target.SetText(uri.Path())
		}, w)
	}

	inputButton := widget.NewButton("选择", func() { pickDir(input) })
	outputButton := widget.NewButton("选择", func() { pickDir(output) })

	var setRunning func(bool)
	var fuckButton *widget.Button
	swapButton := widget.NewButton("交换", func() {
		tmp := input.Text
		input.SetText(output.Text)
		output.SetText(tmp)
		if doSel.Selected == "加密" {
			doSel.SetSelected("解密")
		} else {
			doSel.SetSelected("加密")
		}
	})

	controls := []*widget.Entry{input, output, pass}
	buttons := []*widget.Button{inputButton, outputButton, swapButton}
	selects := []*widget.Select{split, doSel}

	setRunning = func(running bool) {
		if running {
			for _, c := range controls {
				c.Disable()
			}
			for _, b := range buttons {
				b.Disable()
			}
			for _, s := range selects {
				s.Disable()
			}
			fuckButton.Disable()
			return
		}
		for _, c := range controls {
			c.Enable()
		}
		for _, b := range buttons {
			b.Enable()
		}
		for _, s := range selects {
			s.Enable()
		}
		fuckButton.Enable()
	}

	fuckButton = widget.NewButton("GO", func() {
		inPath := strings.TrimSpace(input.Text)
		outPath := strings.TrimSpace(output.Text)
		if inPath == "" || outPath == "" {
			dialog.ShowInformation("提示", "请选择输入和输出目录", w)
			return
		}

		setRunning(true)
		cur.SetValue(0)
		curf.SetValue(0)

		gConfig.Split = split.Selected
		gConfig.Do = doSel.Selected
		gConfig.Input = inPath
		gConfig.Output = outPath
		gConfig.Pass = pass.Text
		saveJson(gConfig)

		gb, _ := strconv.Atoi(strings.TrimSuffix(split.Selected, "G"))
		encrypting := doSel.Selected == "加密"
		key := pass.Text
		p := &Progress{}

		stop := make(chan struct{})
		go func() {
			t := time.NewTicker(100 * time.Millisecond)
			defer t.Stop()
			for {
				select {
				case <-stop:
					return
				case <-t.C:
					jd, jt := p.Jobs()
					fd, ft := p.File()
					fyne.Do(func() {
						if jt > 0 {
							cur.SetValue(float64(jd) / float64(jt))
						}
						if ft > 0 {
							curf.SetValue(float64(fd) / float64(ft))
						}
					})
				}
			}
		}()

		go func() {
			defer close(stop)
			err := dojob(p, inPath, outPath, encrypting, key, 1000*1000*1000*gb)
			fyne.Do(func() {
				if err == nil {
					cur.SetValue(1)
					curf.SetValue(1)
				}
				setRunning(false)
				if err != nil {
					dialog.ShowError(err, w)
					return
				}
				dialog.ShowInformation("完成", "ok", w)
			})
		}()
	})

	form := container.New(&colGrid{cols: 3},
		widget.NewLabel("输入："), input, inputButton,
		widget.NewLabel("输出："), output, outputButton,
		widget.NewLabel("密码："), pass, swapButton,
		split, doSel, fuckButton,
	)
	content := container.NewVBox(form, cur, curf)

	if lg := loadJson(); lg != nil {
		gConfig = *lg
		if gConfig.Do != "" {
			doSel.SetSelected(gConfig.Do)
		}
		if gConfig.Split != "" {
			split.SetSelected(gConfig.Split)
		}
		if gConfig.Input != "" {
			input.SetText(gConfig.Input)
		}
		if gConfig.Output != "" {
			output.SetText(gConfig.Output)
		}
		if gConfig.Pass != "" {
			pass.SetText(gConfig.Pass)
		}
	}

	w.SetPadded(true)
	w.SetContent(content)
	w.Resize(fyne.NewSize(520, content.MinSize().Height+theme.Padding()*2))
	w.SetFixedSize(true)
	w.ShowAndRun()
}

// colGrid 三列对齐：左右列按最宽控件定宽，中间列拉伸，同一行等高。
type colGrid struct {
	cols int
}

func (g *colGrid) minColRow(objects []fyne.CanvasObject) (colW []float32, rowH []float32) {
	cols := g.cols
	if cols < 1 {
		cols = 3
	}
	rows := (len(objects) + cols - 1) / cols
	colW = make([]float32, cols)
	rowH = make([]float32, rows)
	for i, o := range objects {
		if o == nil || !o.Visible() {
			continue
		}
		m := o.MinSize()
		c := i % cols
		r := i / cols
		if m.Width > colW[c] {
			colW[c] = m.Width
		}
		if m.Height > rowH[r] {
			rowH[r] = m.Height
		}
	}
	return
}

func (g *colGrid) MinSize(objects []fyne.CanvasObject) fyne.Size {
	colW, rowH := g.minColRow(objects)
	pad := theme.Padding()
	var w, h float32
	for i, cw := range colW {
		w += cw
		if i > 0 {
			w += pad
		}
	}
	for i, rh := range rowH {
		h += rh
		if i > 0 {
			h += pad
		}
	}
	return fyne.NewSize(w, h)
}

func (g *colGrid) Layout(objects []fyne.CanvasObject, size fyne.Size) {
	cols := g.cols
	if cols < 1 {
		cols = 3
	}
	colW, rowH := g.minColRow(objects)
	pad := theme.Padding()
	var minW float32
	for i, cw := range colW {
		minW += cw
		if i > 0 {
			minW += pad
		}
	}
	if extra := size.Width - minW; extra > 0 && len(colW) >= 2 {
		colW[1] += extra
	}
	y := float32(0)
	for r, rh := range rowH {
		x := float32(0)
		for c := 0; c < cols; c++ {
			i := r*cols + c
			if i >= len(objects) {
				break
			}
			if o := objects[i]; o != nil {
				o.Move(fyne.NewPos(x, y))
				o.Resize(fyne.NewSize(colW[c], rh))
			}
			x += colW[c] + pad
		}
		y += rh + pad
	}
}

type Config struct {
	Input  string `json:"Input"`
	Output string `json:"Output"`
	Do     string `json:"Do"`
	Split  string `json:"Split"`
	Pass   string `json:"Pass"`
}

var gConfig Config

func configPath() string {
	exe, err := os.Executable()
	if err != nil {
		return ".fuckbaiduyun.json"
	}
	return filepath.Join(filepath.Dir(exe), ".fuckbaiduyun.json")
}

func saveJson(c Config) {
	f, err := os.OpenFile(configPath(), os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		return
	}
	defer f.Close()
	_ = json.NewEncoder(f).Encode(&c)
}

func loadJson() *Config {
	f, err := os.Open(configPath())
	if err != nil {
		return nil
	}
	defer f.Close()
	data, err := io.ReadAll(f)
	if err != nil {
		return nil
	}
	var c Config
	if err := json.Unmarshal(data, &c); err != nil {
		return nil
	}
	return &c
}
