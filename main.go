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

	fmt.Print("Press 'Enter' to continue...")
	bufio.NewReader(os.Stdin).ReadBytes('\n')
	questionMap, questionsCount := parseQuiz(filename)

	_, score := startQuiz(questionMap, timeLimit)

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

func startQuiz(q map[string]string, timeLimit int) (map[string][]string, int) {
	result := map[string][]string{}
	score := 0

	go timer(timeLimit)

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
		} else {
			fmt.Println("Wrong!")
			result[k] = []string{v, ans, "wrong"}
		}
	}
	return result, score
}

func timer(t int) {
	time.Sleep(time.Duration(t) * time.Second)
	fmt.Println("\nTimes up!")
	os.Exit(0)
}
