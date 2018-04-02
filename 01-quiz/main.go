package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"time"
)

type quiz struct {
	Questions      []question
	TimeLimit      time.Duration
	CorrectAnswers int
}

type question struct {
	Prompt string
	Answer string
}

func main() {
	csvFile := flag.String("file", "problems.csv", "CSV file containing the quiz questions")
	timeLimit := flag.Int("limit", 30, "Time limit for quiz completion")
	// shuffle := flag.Bool("shuffle", false, "Shuffle problem ordering")
	flag.Parse()

	f, err := os.Open(*csvFile)
	if err != nil {
		panic(fmt.Sprintf("Unable to open CSV file: %s", err))
	}
	defer f.Close()

	questions, err := parseFile(f)
	if err != nil {
		panic(fmt.Sprintf("Unable to parse CSV file: %s", err))
	}

	quiz := &quiz{
		Questions: questions,
		TimeLimit: time.Duration(*timeLimit) * time.Second,
	}
	quiz.Run()
}

func parseFile(f *os.File) ([]question, error) {
	r := csv.NewReader(f)
	lines, err := r.ReadAll()
	if err != nil {
		return nil, err
	}

	questions := make([]question, len(lines))
	for i, line := range lines {
		questions[i] = question{
			Prompt: line[0],
			Answer: line[1],
		}
	}
	return questions, nil
}

func (q *quiz) Run() {
	timer := time.NewTimer(q.TimeLimit)

loop:
	for i, question := range q.Questions {
		fmt.Printf("#%d: %s\n", i+1, question.Prompt)
		answerCh := make(chan string)
		go func() {
			var answer string
			fmt.Scanf("%s\n", &answer)
			answerCh <- answer
		}()

		select {
		case <-timer.C:
			break loop
		case answer := <-answerCh:
			if answer == question.Answer {
				q.CorrectAnswers++
			}
		}
	}

	fmt.Printf("You scored %d out of %d.\n", q.CorrectAnswers, len(q.Questions))
}
