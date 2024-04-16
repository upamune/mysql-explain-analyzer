package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	"mysql-explain-analyzer/model/input"
	"mysql-explain-analyzer/model/output"
)

type CLI struct {
	Stdin  io.ReadCloser
	Stderr io.Writer
	Stdout io.Writer
}

func (c *CLI) Run(args []string) error {
	var filename string

	if len(args) > 1 {
		filename = args[1]
	}

	r, err := c.read(filename)
	if err != nil {
		return fmt.Errorf("failed to read %s: %v", args[0], err)
	}
	defer r.Close()

	var input input.Explain
	if err := json.NewDecoder(r).Decode(&input); err != nil {
		return fmt.Errorf("failed to decode JSON: %v", err)
	}

	res := convert(input)
	if err := write(
		c.Stdout,
		res,
	); err != nil {
		return fmt.Errorf("failed to write output: %v", err)
	}

	return nil
}

func (cli *CLI) read(filename string) (io.ReadCloser, error) {
	if filename == "" || filename == "-" {
		return cli.Stdin, nil
	}
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	return f, nil
}

func write(w io.Writer, output output.Result) error {
	printExplainTable(w, output)
	printFullTableScanComments(w, output)
	printFullIndexScanComments(w, output)
	printAnythingElseComments(w, output)
	return nil
}
