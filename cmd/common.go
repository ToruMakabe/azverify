package cmd

import (
	"fmt"
	"os"
)

func checkErr(msg interface{}) {
	if msg != nil {
		fmt.Fprintln(os.Stderr, "[ERROR] ", msg)
		os.Exit(1)
	}
}
