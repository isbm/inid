package main

import (
	"fmt"
	"os"
	"path"

	"github.com/isbm/inid"
	"github.com/isbm/inid/sysctl"
)

func main() {
	mode := path.Base(os.Args[0])
	if mode == "inidctl" {
		if err := sysctl.NewSvmCtl().Run(); err != nil {
			fmt.Printf("Error: %s\n", err.Error())
		}
	} else {
		// Init system mode
		svm := inid.NewSVM()
		if err := svm.Init(); err != nil { // Look for /etc/runit/rc.d and load the map
			panic(err)
		}

		svm.Run()
	}
}
