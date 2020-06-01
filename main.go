package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"log"
	"os"
)

func main() {
	file, err := os.Open("questions.csv")
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

	_, score := startQuiz(questionMap, 30)

	fmt.Printf("You scored %v/%v!\n", score, len(questions))

}

func startQuiz(q map[string]string, n int) (map[string][]string, int) {
	result := map[string][]string{}
	score := 0

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
