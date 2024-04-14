package main

import (
	"fmt"
	"os"
)

const (
	exitOK  = 0
	exitErr = 1
)

var (
	Version   string
	CommitSHA string
)

func main() {
	os.Exit(realMain())
}

func realMain() int {
	cli := &CLI{
		Stdin:  os.Stdin,
		Stdout: os.Stdout,
	}
	if err := cli.Run(os.Args); err != nil {
		fmt.Fprintln(os.Stderr, err)
		return exitErr
	}
	return exitOK
}
