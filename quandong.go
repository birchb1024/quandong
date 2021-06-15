package main

import (
	"fmt"
	"os"
)

var Version = "none"
var printVersion bool

func main() {

	if printVersion {
		fmt.Println(Version)
		os.Exit(0)
	}

}
