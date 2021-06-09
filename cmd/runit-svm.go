package main

import (
	runit_svm "github.com/isbm/runit-svm"
)

func main() {
	svm := runit_svm.NewSVM()
	if err := svm.Init(); err != nil { // Look for /etc/runit/rc.d and load the map
		panic(err)
	}

	svm.Run()
}
