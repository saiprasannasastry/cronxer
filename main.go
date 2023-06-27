package main

import (
	"fmt"
	crx "github.com/segmentio/cron-parser/cronxer"
	"os"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run main.go <cron-string>")
		return
	}

	cronString := os.Args[1]
	parser := crx.New()
	output, err := parser.Parse(cronString)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Println(output)

	parser.GetNextCronJobs(cronString)
}
//5 next times.
// "* 0 1,15 * 1-5" "5"