package main

import (
	"flag"
	"fmt"
	"os"
	"sync"
	"time"
)

func main() {

	var timeout int64
	flag.Int64Var(&timeout, "t", 5, "timeout value in seconds")
	flag.Int64Var(&timeout, "timeout", 5, "timeout value in seconds")
	flag.Parse()

	dirtyNodes := flag.Args()
	nodes, invalidNodes := cleanNodes(dirtyNodes)

	if len(nodes) == 0 {
		fmt.Println("No valid nodes supplied.")
		os.Exit(0)
	}

	if len(invalidNodes) > 0 {
		printInvalidNodes(invalidNodes)
	}

	timeoutDuration := time.Duration(timeout) * time.Second
	var wg sync.WaitGroup

	for {
		// used to force each iteration to wait for the timeout
		// TODO check if this is idiomatic Go
		wg.Add(1)
		go func() {
			time.Sleep(timeoutDuration)
			wg.Done()
		}()

		printTimestamp()

		for _, node := range nodes {
			wg.Add(1)
			go processNode(node, &wg, timeoutDuration)
		}

		wg.Wait()
		fmt.Println()
	}
}
