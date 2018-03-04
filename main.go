package main

import (
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

	go func() {
		for range time.Tick(time.Second) {
			err := cmd.Run()

			if err != nil {
				log.Fatal(err)
			}
		}
	}()

	for {
		n, _ := out.ReadString('\n')
		fmt.Println(n)
	}

	fmt.Printf("in all caps: %q\n", out.String())
}
