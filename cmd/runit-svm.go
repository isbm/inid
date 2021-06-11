package main

import (
	"fmt"
	"os"
	"path"

	runit_svm "github.com/isbm/runit-svm"
	"github.com/isbm/runit-svm/sysctl"
)

func main() {
	mode := path.Base(os.Args[0])
	if mode == "inidctl" {
		if err := sysctl.NewSvmCtl().Run(); err != nil {
			fmt.Printf("Error: %s\n", err.Error())
		}
	} else {
		// Init system mode
		svm := runit_svm.NewSVM()
		if err := svm.Init(); err != nil { // Look for /etc/runit/rc.d and load the map
			panic(err)
		}

		svm.Run()
	}
}
