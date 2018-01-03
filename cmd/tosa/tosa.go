package main

import (
	"os"

	"github.com/kyoshidajp/tosa"
)

func main() {
	os.Exit(_main())
}

func _main() int {
	cli := tosa.New()
	if exitCode := cli.Run(os.Args); exitCode != tosa.ExitCodeOK {
		return exitCode
	}
	return tosa.ExitCodeOK
}
