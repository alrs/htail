package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
)

var herokuApp string

func printLoop(outCh chan io.Reader, errCh chan error) {
	for {
		scanner := bufio.NewScanner(<-outCh)
		for scanner.Scan() {
			fmt.Println(scanner.Text())
			err := scanner.Err()
			if err != nil {
				errCh <- err
				return
			}
		}
	}
}

func execLoop(outCh chan io.Reader, errCh chan error) {
	for {
		cmd := exec.Command("heroku", "logs", "--tail", "--no-color", "-a", herokuApp)
		stdOut, err := cmd.StdoutPipe()
		if err != nil {
			errCh <- err
			return
		}
		outCh <- stdOut
		err = cmd.Start()
		if err != nil {
			errCh <- err
			return
		}
	}
}

func main() {
	flag.StringVar(&herokuApp, "a", os.Getenv("HEROKU_APP"),
		"Heroku application to tail, defaults to environment HEROKU_APP")
	flag.Parse()
	if herokuApp == "" {
		log.Fatal("select a Heroku application either with -a or HEROKU_APP")
	}
	outCh := make(chan io.Reader)
	errCh := make(chan error)
	go printLoop(outCh, errCh)
	go execLoop(outCh, errCh)
	err := <-errCh
	if err != nil {
		log.Fatal(err)
	}
}
