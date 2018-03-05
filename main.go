package main

import (
	"bytes"
	"fmt"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

const (
	averageOver = 5
)

func main() {
	var over []float64

	for i := range time.Tick(time.Second) {
		out, err := utilization()

		if err != nil {
			panic(err)
		}

		over = append(over, averageFromString(out))

		if i.Second()%averageOver == 0 {
			over = []float64{}
		}

		fmt.Println("utilization:", averageFromFloats(over))
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
