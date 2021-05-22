package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"
)

func main() {
	csvFilename := flag.String("csv", "problems.csv", "a csv file -> 'question,answer'")
	timeLimit := flag.Int("limit", 30, "this is the time limit for quiz in sec")
	flag.Parse()

	file, err := os.Open(*csvFilename)
	if err != nil {
		exit(fmt.Sprintf("Failed to open csv file %s", *csvFilename))
	}
	r := csv.NewReader(file)
	line, err := r.ReadAll()
	if err != nil {
		exit("Failed to parse the provided csv file")
	}
	problems := parseLines(line)

	timer := time.NewTimer(time.Duration(*timeLimit) * time.Second)

	correct := 0
	for i, p := range problems {
		fmt.Printf("Problem #%d: %s = \n", i+1, p.q)
		answerChan := make(chan string)
		go func() {
			var answer string
			fmt.Scanf("%s\n", &answer)
			answerChan <- answer

		}()

		select {
		case <-timer.C:
			fmt.Printf("Problem #%d: %s = \n", i+1, p.q)
			return
		case answer := <-answerChan:
			if answer == p.a {
				correct++
			}

		}

	}

	fmt.Printf("You scored %d out of %d \n", correct, len(problems))
}

func parseLines(lines [][]string) []problem {
	ret := make([]problem, len(lines))
	for i, line := range lines {
		ret[i] = problem{
			q: line[0],
			a: strings.TrimSpace(line[1]),
		}
	}

	return ret
}

type problem struct {
	q string
	a string
}

func exit(msg string) {
	fmt.Println(msg)
	os.Exit(1)
}
