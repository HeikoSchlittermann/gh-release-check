package main

import (
	"io"
	"log"
	"os"
	"os/exec"
)

func jq(b io.ReadCloser) {
	jq := exec.Command("jq", ".")
	jq.Stdout = os.Stdout
	pipe, err := jq.StdinPipe()
	if err != nil {
		log.Fatal(err)
	}
	if err := jq.Start(); err != nil {
		log.Println("the attempt to start jq returned:", err)
		pipe = os.Stdout
	}
	io.Copy(pipe, b)
	pipe.Close()
	_ = jq.Wait()
	return
}
