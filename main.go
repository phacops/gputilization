package main

import (
	"bytes"
	"os/exec"
	"strconv"
	"strings"

	"github.com/gizak/termui"
)

const (
	averageOverInSeconds = 5
)

func main() {
	err := termui.Init()

	if err != nil {
		panic(err)
	}

	defer termui.Close()

	var over []float64

	p := termui.NewPar("Press q to quit")
	p.Height = 3
	p.TextFgColor = termui.ColorWhite
	p.BorderLabel = "Instructions"
	p.BorderFg = termui.ColorWhite

	g := termui.NewGauge()
	g.Height = 10
	g.BorderLabel = "GPU utilization"
	g.PercentColor = termui.ColorBlue | termui.AttrBold
	g.BarColor = termui.ColorYellow
	g.BorderFg = termui.ColorWhite
	g.BorderLabelFg = termui.ColorCyan
	g.LabelAlign = termui.AlignRight
	g.Label = "{{percent}}%"

	lc := termui.NewLineChart()
	lc.Mode = "dot"
	lc.Data = func() []float64 { return over }()
	lc.AxesColor = termui.ColorWhite
	lc.LineColor = termui.ColorCyan | termui.AttrBold
	lc.X = 300
	lc.Y = 100

	termui.Body.AddRows(
		termui.NewRow(
			termui.NewCol(12, 0, p),
		),
		termui.NewRow(
			termui.NewCol(12, 0, g),
		),
		termui.NewRow(
			termui.NewCol(12, 0, lc),
		),
	)

	termui.Body.Align()
	termui.Render(termui.Body)

	termui.Handle("/timer/1s", func(e termui.Event) {
		out, err := utilization()

		if err != nil {
			panic(err)
		}

		over = append(over, averageFromString(out))
		count := len(over)

		var from int

		if count >= averageOverInSeconds {
			from = count - averageOverInSeconds
		}

		g.Percent = int(averageFromFloats(over[from:]))

		termui.Body.Align()
		termui.Render(termui.Body)
	})

	termui.Handle("/sys/kbd/q", func(termui.Event) {
		termui.StopLoop()
	})

	termui.Handle("/sys/wnd/resize", func(e termui.Event) {
		termui.Body.Width = termui.TermWidth()
		termui.Body.Align()
		termui.Clear()
		termui.Render(termui.Body)
	})

	termui.Loop()
}

func averageFromFloats(f []float64) float64 {
	var o float64

	for _, v := range f {
		o += v
	}

	return o / float64(len(f))
}

func averageFromString(s string) float64 {
	v := strings.Split(s, "\n")

	var o int

	for _, n := range v {
		i, _ := strconv.Atoi(n)
		o += i
	}

	return float64(o / len(v))
}

func utilization() (string, error) {
	var out bytes.Buffer

	cmd := exec.Command("nvidia-smi", "--query-gpu=utilization.gpu", "--format=csv,noheader,nounits")
	cmd.Stdout = &out

	err := cmd.Run()

	if err != nil {
		return "", err
	}

	return strings.TrimSpace(out.String()), nil
}
