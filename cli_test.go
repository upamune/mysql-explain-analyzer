package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
	"testing"

	"github.com/tenntenn/golden"
)

var flagUpdate bool

func init() {
	flag.BoolVar(&flagUpdate, "update", false, "update golden files")
}

func Test_CLI(t *testing.T) {
	f, err := os.Open("testdata/sample_explain.json")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { f.Close() })

	rawJSON, err := io.ReadAll(f)
	if err != nil {
		t.Fatal(err)
	}

	cases := map[string]struct {
		args    string
		in      string
		wantErr bool
	}{
		"file":  {"mea ./testdata/sample_explain.json", "", false},
		"stdin": {"mea", string(rawJSON), false},
	}

	for name, tt := range cases {
		name, tt := name, tt
		t.Run(name, func(t *testing.T) {
			t.Setenv("NO_COLOR", "true")

			var stdout, stderr bytes.Buffer
			cli := &CLI{
				Stdout: &stdout,
				Stderr: &stderr,
				Stdin:  io.NopCloser(strings.NewReader(tt.in)),
			}

			args := strings.Split(tt.args, " ")
			err := cli.Run(args)
			if (err != nil) != tt.wantErr {
				t.Errorf("unexpected error: %v", err)
				return
			}

			c := golden.New(t, flagUpdate, "testdata", fmt.Sprintf("%s_%s", t.Name(), name))

			if diff := c.Check("_stdout", &stdout); diff != "" {
				t.Error("stdout\n", diff)
			}

			if diff := c.Check("_stderr", &stderr); diff != "" {
				t.Error("stderr\n", diff)
			}
		})
	}
}
