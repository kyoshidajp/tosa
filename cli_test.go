package main

import (
	"bytes"
	"fmt"
	"strings"
	"testing"
)

func TestRun(t *testing.T) {
	outStream, errStream := new(bytes.Buffer), new(bytes.Buffer)
	cli := &CLI{outStream: outStream, errStream: errStream}

	command := fmt.Sprintf(
		"tosa -debug=true -url %s", "c97e6909")

	args := strings.Split(command, " ")
	if got, want := cli.Run(args), ExitCodeOK; got != want {
		t.Fatalf("%q exits %d, want %d\n\n%s", command, got, want, errStream.String())
	}
}
