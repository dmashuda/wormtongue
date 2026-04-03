package main

import (
	"os"

	"github.com/dmashuda/wormtongue/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
