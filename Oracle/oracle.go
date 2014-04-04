// Stefan Nilsson 2013-03-13

// This program implements an ELIZA-like oracle (en.wikipedia.org/wiki/ELIZA).
package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"regexp"
	"strings"
	"time"
)

const (
	star                = "Pythia"
	venue               = "Delphi"
	prompt              = "> "
	starshorttermmemory = 10 //Pythia can hold ten questions in her short term memory until the drugs kick in.
)

func main() {
	fmt.Printf("Welcome to %s, the oracle at %s.\n", star, venue)
	fmt.Println("Your questions will be answered in due time.")

	oracle := Oracle()
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print(prompt)
		line, _ := reader.ReadString('\n')
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		fmt.Printf("%s heard: %s\n", star, line)
		oracle <- line // The channel doesn't block.
	}
}

// Oracle returns a channel on which you can send your questions to the oracle.
// You may send as many questions as you like on this channel, it never blocks.
// The answers arrive on stdout, but only when the oracle so decides.
// The oracle also prints sporadic prophecies to stdout even without being asked.
func Oracle() chan<- string {
	questions := make(chan string, starshorttermmemory)
	prophecies := make(chan string, starshorttermmemory)

	go generateProphecies(questions, prophecies)
	go printAnswers(prophecies)

	return questions
}

func printAnswers(answers <-chan string) {
	for answer := range answers {
		fmt.Print("\r")
		for _, char := range answer {
			time.Sleep(time.Duration(20+rand.Intn(10)) * time.Millisecond)
			fmt.Print(string(char))
		}
		fmt.Print("\n", prompt)
	}
}

func generateProphecies(questions chan string, prophecies chan string) {
	for {
		// Keep them waiting. Pythia, the original oracle at Delphi,
		// only gave prophecies on the seventh day of each month.
		time.Sleep(time.Duration(10+rand.Intn(5)) * time.Second)
		select {
		case question := <-questions:
			generateProphecy(question, prophecies)
		default:
			generateProphecy("", prophecies)
		}
	}
}

func generateProphecy(question string, prophecies chan<- string) {

	if question == "" {
		lonely := []string{
			"A shy man is man and the path he walks must be that of his own. Does that sound too cheesy?",
			"GET ME MY WATER! This is SO BORIIING!",
			"I can't believe they all fall for my nonsense. Hahaha",
			"WHERE ARE MY PILLS?! Aldora, you took them didn't you, you filthy theif!",
		}
		prophecies <- "...\nOk, no one's here, right?\n" + lonely[rand.Intn(len(lonely))]
	} else {
		prophecy(question, prophecies)
	}
}

// This is the oracle's secret algorithm.
// It waits for a while and then sends a message on the answer channel.
// TODO: make it better.
func prophecy(question string, prophecies chan<- string) {

	disrespect, _ := regexp.MatchString("(?i)(fuck|you|damn|bitch|whore)", question)
	if disrespect {
		prophecies <- "How DARE YOU!?!?! Guards!!!!"
		return
	}
	humble, _ := regexp.MatchString("(?i)help", question)
	properAddress, _ := regexp.MatchString("(?i)("+star+"|mighty|highness|master|queen|wise|beautiful)", question)
	humble = humble && properAddress
	if !humble {
		prophecies <- "Show respect for the oracle! Peasant!"
		return
	}
	// Find the longest word.
	longestWord := ""
	words := strings.Fields(question) // Fields extracts the words into a slice.
	for _, w := range words {
		if len(w) > len(longestWord) {
			longestWord = w
		}
	}

	// Cook up some pointless nonsense.
	nonsense := []string{
		"The moon is dark.",
		"The sun is bright.",
		"If you look then you see.",
		"The fog is deepening...",
		"The grass is growing every spring.",
		"Snow is both white and cold.",
	}
	prophecies <- nonsense[rand.Intn(len(nonsense))]
}

func init() { // Functions called "init" are executed before the main function.
	// Use new pseudo random numbers every time.
	rand.Seed(time.Now().Unix())
}
