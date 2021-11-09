package main

import (
	"fmt"

	"github.com/made-from-organic-orange-juice/task/wpsapi"
)

func main() {
	NextSystemSnapShot := wpsapi.SystemSnapShot{}.SystemSnapShotIterator()

	for {

		sysSnap, err := NextSystemSnapShot()
		if err != nil {
			fmt.Printf("%s", err)
			break
		}

		fmt.Printf("Process: %s \n", sysSnap.Process.Name)
		fmt.Printf("----> Modules: \n")

		nextModule := sysSnap.Modules.Iterator()
		if err != nil {
			fmt.Printf("%s", err)
			break
		}

		for m, err := nextModule(); err == nil; m, err = nextModule() {
			fmt.Printf("-------> %s %s\n", m.BaseName, m.Path)
		}

	}

}
