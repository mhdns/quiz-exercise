package main

import (
	"bufio"
	"encoding/csv"
	"flag"
	"fmt"
	"log"
	"os"
	"time"
)

func main() {
	var filename string
	var timeLimit int
	flag.StringVar(&filename, "filename", "questions.csv", "Declare the csv file containing the questions")
	flag.IntVar(&timeLimit, "time", 10, "Declare test time limit in seconds")
	flag.Parse()

	questionMap, questionsCount := parseQuiz(filename)

	_, score := startQuiz(questionMap, questionsCount, timeLimit)

	fmt.Printf("You scored %v/%v!\n", score, questionsCount)

}

func parseQuiz(filename string) (map[string]string, int) {
	// fmt.Println(filename, " being parsed...")
	file, err := os.Open(filename)
	if err != nil {
		log.Fatalln(err)
	}
	defer file.Close()

	data := csv.NewReader(file)

	questions, err := data.ReadAll()
	if err != nil {
		log.Fatalln(err)
	}

	questionMap := map[string]string{}
	for _, v := range questions {
		questionMap[v[0]] = v[1]
	}

	return questionMap, len(questions)
}

func startQuiz(q map[string]string, questionCount, timeLimit int) (map[string][]string, int) {
	result := map[string][]string{}
	score := 0
	scoreChan := make(chan int, questionCount+1)
	totalChan := make(chan int, 1)
	closeChan := make(chan int)

	fmt.Print("Press 'Enter' to continue...")
	bufio.NewReader(os.Stdin).ReadBytes('\n')

	totalChan <- questionCount
	go timer(timeLimit, scoreChan, totalChan, closeChan)

	go func(closeChan chan int) {
		select {
		case <-closeChan:
			close(scoreChan)
			close(totalChan)
			return
		}
	}(closeChan)

	for k, v := range q {

		fmt.Println("Question: ", k)

		reader := bufio.NewReader(os.Stdin)
		fmt.Print("Answer: ")
		ans, _ := reader.ReadString('\n')
		ans = ans[0 : len(ans)-1]

		if ans == v {
			result[k] = []string{v, ans, "correct"}
			fmt.Println("Correct!")
			score++
			scoreChan <- 1
		} else {
			fmt.Println("Wrong!")
			result[k] = []string{v, ans, "wrong"}
		}
	}
	close(scoreChan)
	close(totalChan)
	close(closeChan)
	return result, score
}

func timer(t int, scoreChan, totalChan, closeChan chan int) {
	time.Sleep(time.Duration(t) * time.Second)

	questionCount := <-totalChan
	finalScore := 0
	closeChan <- 0
	close(closeChan)

	for v := range scoreChan {
		finalScore += v
	}

	fmt.Printf("\nTimes up! You scored %v/%v!\n", finalScore, questionCount)
	os.Exit(0)
}
