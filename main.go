package main

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"os/exec"
	"time"
)

func main() {
	var out bytes.Buffer

	scanner := bufio.NewScanner(&out)

	go func() {
		for range time.Tick(time.Second) {
			cmd := exec.Command("nvidia-smi", "--query-gpu=utilization.gpu", "--format=csv,noheader,nounits")
			cmd.Stdout = &out

			err := cmd.Run()

			if err != nil {
				log.Fatal(err)
			}
		}
	}()

	for {
		for scanner.Scan() {
			fmt.Println(scanner.Text())
		}
	}
}
