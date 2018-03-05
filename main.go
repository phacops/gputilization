package main

import (
	"bytes"
	"fmt"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

func main() {
	var i int
	var over [5]float64

	for range time.Tick(time.Second) {
		out, err := utilization()

		if err != nil {
			panic(err)
		}

		over[i] = averageFromString(out)

		if i == 4 {
			i = 0
		}

		fmt.Println("utilization:", averageFromFloats(over))
	}
}

func averageFromFloats(f [5]float64) float64 {
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
