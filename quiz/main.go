package main

import (
	"flag"
	"fmt"
	"encoding/csv"
	"os"
	"time"
)

func main() {
	csvPtr := flag.String("csv", "problems.csv", "path of csv file in the format of (question,answer). (default \"problems.csv\")")
	limitPtr := flag.Int("limit", 30, "time limit in seconds (default 30)")

	flag.Parse()

	csvFile, _ := os.Open(*csvPtr)
	lines, _ := csv.NewReader(csvFile).ReadAll()

	timer := time.NewTimer(time.Duration(*limitPtr) * time.Second)

	var correct int

eventLoop:
	for _, line := range lines {
		fmt.Print("Problem #" + "1: " + line[0] + " = ")

		answerCh := make(chan string)
		go func() {
			var input string
			_, _ = fmt.Scanln(&input)
			answerCh <- input
		}()

		select {
		case <-timer.C:
			fmt.Println()
			fmt.Println("Times up!")
			break eventLoop
		case answer := <-answerCh:
			if answer == line[1] {
				correct++
			}
		}
	}

	fmt.Printf("You scored %d out of %d\n", correct, len(lines))
}
