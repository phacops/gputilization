package main

import (
	"bytes"
	"os/exec"
	"strconv"
	"strings"
	"time"

	termbox "github.com/nsf/termbox-go"
)

const (
	averageOver     = 5
	utilizationText = "Utilization:"
	percent         = "%"
)

func main() {
	err := termbox.Init()

	if err != nil {
		panic(err)
	}

	defer termbox.Close()

	var over []float64

	for i := range time.Tick(time.Second) {
		out, err := utilization()

		if err != nil {
			panic(err)
		}

		over = append(over, averageFromString(out))
		err = termbox.Clear(termbox.ColorWhite, termbox.ColorBlack)

		if err != nil {
			panic(err)
		}

		average := averageFromFloats(over)
		t := strconv.FormatFloat(average, 'g', 2, 64)

		var p int

		for _, c := range utilizationText {
			termbox.SetCell(p, 0, rune(c), termbox.ColorWhite, termbox.ColorBlack)
			p += 1
		}

		for _, c := range t {
			termbox.SetCell(p, 0, rune(c), termbox.ColorWhite, termbox.ColorBlack)
			p += 1
		}

		termbox.SetCell(p, 0, '%', termbox.ColorWhite, termbox.ColorBlack)

		err = termbox.Flush()

		if err != nil {
			panic(err)
		}

		if i.Second()%averageOver == 0 {
			over = []float64{}
		}
	}
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
