package main

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
	"time"
)

func main() {
	for range time.Tick(time.Second) {
		out, err := utilization()

		if err != nil {
			panic(err)
		}

		fmt.Println(strings.Split(out, string('\n')))
	}
}

func utilization() (string, error) {
	var out bytes.Buffer

	cmd := exec.Command("nvidia-smi", "--query-gpu=utilization.gpu", "--format=csv,noheader,nounits")
	cmd.Stdout = &out

	err := cmd.Run()

	if err != nil {
		return "", err
	}

	return out.String(), nil
}
