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

	g := termui.NewGauge()
	g.Width = 100
	g.Height = 20
	g.Y = 6
	g.BorderLabel = "GPU utilization"
	g.PercentColor = termui.ColorYellow
	g.BarColor = termui.ColorGreen
	g.BorderFg = termui.ColorWhite
	g.BorderLabelFg = termui.ColorMagenta

	var over []float64

	termui.Handle("/timer/1s", func(e termui.Event) {
		c := e.Data.(termui.EvtTimer)
		out, err := utilization()

		if err != nil {
			panic(err)
		}

		over = append(over, averageFromString(out))
		g.Percent = int(averageFromFloats(over))

		termui.Render(g)

		if c.Count%averageOverInSeconds == 0 {
			over = []float64{}
		}
	})

	termui.Handle("/sys/kbd/q", func(termui.Event) {
		termui.StopLoop()
	})

	termui.Loop()
}

func averageOver(g *termui.Gauge) {
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
