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

	cmd := exec.Command("nvidia-smi", "--query-gpu=utilization.gpu", "--format=csv,noheader,nounits")
	cmd.Stdout = &out

	scanner := bufio.NewScanner(&out)

	go func() {
		for range time.Tick(time.Second) {
			err := cmd.Run()

			if err != nil {
				log.Fatal(err)
			}
		}
	}()

	for scanner.Scan() {
		fmt.Println(scanner.Text())
	}

	fmt.Printf("in all caps: %q\n", out.String())
}
