package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"
)

type Question struct {
	Prompt string
	Answer string
}

func main() {
	var problemFile string
	flag.StringVar(&problemFile, "file", "problems.csv", "CSV file containing the quiz questions")
	flag.Parse()

	f, err := os.Open(problemFile)
	check(err)
	defer f.Close()

	s := bufio.NewScanner(f)

	var questions []Question
	for s.Scan() {
		row := string(s.Bytes())
		question := parseQuestionFromCSV(row)
		questions = append(questions, question)
	}
	check(s.Err())

	runQuiz(questions)
}

func parseQuestionFromCSV(row string) Question {
	parts := strings.Split(row, ",")
	return Question{
		Prompt: parts[0],
		Answer: parts[1],
	}
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func runQuiz(questions []Question) {
	rightAnswers := 0
	reader := bufio.NewReader(os.Stdin)
	for i, question := range questions {
		fmt.Printf("#%d: %s\n", i+1, question.Prompt)
		answer, err := reader.ReadString('\n')
		check(err)
		if strings.TrimSuffix(answer, "\n") == question.Answer {
			rightAnswers++
		}
	}
	fmt.Printf("You scored %d out of %d.\n", rightAnswers, len(questions))
}
