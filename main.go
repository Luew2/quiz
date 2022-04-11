package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"
)

var csvFile *string
var timer *int
var problems []problem

func init() {
	csvFile = flag.String("csv", "problems.csv",
		"a csv file in the format 'queston,answer' for your quiz!")
	timer = flag.Int("timer", 60, "an integer input in terms of seconds for time limiting the quiz.")
	flag.Parse()
}

func main() {
	file := loadProblem(*csvFile)
	parse(file)
	quiz(*timer)
}

func loadProblem(csv string) *os.File {
	file, err := os.Open(csv)
	if err != nil {
		exit(fmt.Sprintf("Failed to open the CSV file: %s", csv))
	}
	return file

}

func parse(f *os.File) {
	r := csv.NewReader(f)
	lines, err := r.ReadAll()
	if err != nil {
		exit("Failed to parse CSV file")
	}
	problems = parseLines(lines)
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

func quiz(t int) {
	fmt.Printf("You have %v seconds to finish this quiz:\n", t)
	correct := 0
	stopWatch := time.NewTimer(time.Duration(t) * time.Second)

problemLoop:
	for i, p := range problems {
		fmt.Printf("problem #%d: %s\n", i+1, p.q)
		answerCh := make(chan string)
		go func() {
			var answer string
			fmt.Scanf("%s\n", &answer)
			answerCh <- answer
		}()
		select {
		case <-stopWatch.C:
			fmt.Println()
			break problemLoop
		case answer := <-answerCh:
			if answer == p.a {
				correct++
				fmt.Printf("Correct! Total number correct: %d\n", correct)
			} else {
				fmt.Printf("Incorrect! Total number correct: %d\n", correct)
			}
		}
	}
	fmt.Printf("\nYou scored %d out of %d.\n", correct, len(problems))
}

type problem struct {
	q string
	a string
}

func exit(msg string) {
	fmt.Println(msg)
	os.Exit(1)
}
