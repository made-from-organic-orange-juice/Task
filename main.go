package main

import (
	"fmt"
	"log"

	"github.com/made-from-organic-orange-juice/task/wpsapi"
)

func main() {

	wModules, err := wpsapi.SystemSnapshot{}.New()
	if err != nil {
		log.Fatalf("error: %s", err)
	}

	instancesMap := wModules.CountInstances()

	for key, val := range instancesMap {
		fmt.Printf("Process: %s, Instances: %d\n", key, val)
	}

}
