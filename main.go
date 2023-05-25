package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"sync"
	"time"
)

type problem struct {
	q string
	a string
}

func parseLines(lines [][]string) []problem {
	ret := make([]problem, len(lines))

	for i, line := range lines {
		ret[i] = problem{
			q: line[0],
			a: line[1],
		}
	}
	return ret
}

func say(text string, wg *sync.WaitGroup) {
	fmt.Println(text)
}

func exit(msg string) {
	fmt.Println("Error al abrir: ", msg)
	os.Exit(1)
}

func startQuizz(problems []problem, timeLimit int) {

	timer := time.NewTimer(time.Duration(timeLimit) * time.Second)

	for _, record := range problems {

		answerCh := make(chan string)

		go func() {
			var answer string
			fmt.Scanln(&answer)
			answerCh <- answer
		}()

		select {
		case <-timer.C:
			fmt.Println("Time is out!!!")
			return
		case answer := <-answerCh:
			if answer == record.a {
				fmt.Println("Correcto!!!")
			} else {
				fmt.Println("Error!!!")
			}
		default:

		}
	}
}

func main() {

	csvFileName := flag.String("csv", "quizzgo.csv", "a csv file in the format of 'question, answer'")
	timeLimit := flag.Int("time", 30, "Time limit for the quizz")

	flag.Parse()

	file, err := os.Open(*csvFileName)

	if err != nil {
		exit(*csvFileName)
		return
	}

	defer file.Close()

	reader := csv.NewReader(file)

	records, err := reader.ReadAll()

	if err != nil {
		fmt.Println("Error al abrir: ", *csvFileName)
		os.Exit(1)
		return
	}

	problems := parseLines(records)

	fmt.Println(problems)

	startQuizz(problems, *timeLimit)

}
