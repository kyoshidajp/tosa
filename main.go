package main

import "os"

// Version of tosa
const Version string = "v0.0.3"

func main() {
	cli := &CLI{outStream: os.Stdout, errStream: os.Stderr}
	os.Exit(cli.Run(os.Args))
}
