package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"math/rand"
	"os"
	"time"
)

func main() {
	csvPtr := flag.String("csv", "problems.csv", "path of csv file in the format of (question,answer). (default \"problems.csv\")")
	limitPtr := flag.Int("limit", 30, "time limit in seconds (default 30)")
	shufflePtr := flag.Bool("shuffle", false, "shuffles the questions each run")

	flag.Parse()

	problems := parseProblems(*csvPtr, *shufflePtr)
	timeout := time.After(time.Duration(*limitPtr) * time.Second)

	correct := 0
eventLoop:
	for i, problem := range problems {
		fmt.Printf("Problem #%d: %s = ", i+1, problem.question)
		answerCh := answerGenerator()

		select {
		case <-timeout:
			fmt.Println()
			fmt.Println("Times up!")
			break eventLoop
		case answer := <-answerCh:
			if answer == problem.answer {
				correct++
			}
		}
	}

	fmt.Printf("You scored %d out of %d\n", correct, len(problems))
}

func answerGenerator() <-chan string {
	answerCh := make(chan string)
	go func() {
		var input string
		_, _ = fmt.Scanln(&input)
		answerCh <- input
	}()
	return answerCh
}

func parseProblems(csvName string, shuffle bool) []problem {
	csvFile, _ := os.Open(csvName)
	lines, _ := csv.NewReader(csvFile).ReadAll()
	problems := make([]problem, len(lines))
	for i, line := range lines {
		problems[i] = problem{
			question: line[0],
			answer:   line[1],
		}
	}
	if shuffle {
		shuffledProblems := make([]problem, len(lines))
		r := rand.New(rand.NewSource(time.Now().Unix()))
		for i, r := range r.Perm(len(problems)) {
			shuffledProblems[i] = problems[r]
		}
		return shuffledProblems
	}
	return problems
}

type problem struct {
	question string
	answer   string
}
