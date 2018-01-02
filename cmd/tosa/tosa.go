package main

import (
	"fmt"
	"os"

	"github.com/kyoshidajp/tosa"
)

func main() {
	os.Exit(_main())
}

func _main() int {
	cli := tosa.New()
	if err := cli.Run(); err != nil {
		fmt.Fprintf(os.Stdout, "%v\n", err)
		return 1
	}
	return 0
}
