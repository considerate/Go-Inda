// Stefan Nilsson 2013-03-13

// This is a testbed to help you understand channels better.
package main

import (
	"fmt"
	"math/rand"
	"strconv"
	"sync"
	"time"
)

func main() {
	// Use different random numbers each time this program is executed.
	rand.Seed(time.Now().Unix())

	const strings = 32
	const producers = 4
	/*
		Vad händer om man ökar antalet konsumenter från 2 till 4?
		Den använder 4 goroutines för att läsa av kanalen. Detta kan teoretiskt
		medföra att processen kan paralleliseras bättre men det innebär också att
		mer minne används.
	*/
	const consumers = 2

	before := time.Now()
	ch := make(chan string)
	wgp := new(sync.WaitGroup)
	wgp.Add(producers)
	for i := 0; i < producers; i++ {
		go Produce("p"+strconv.Itoa(i), strings/producers, ch, wgp)
	}
	for i := 0; i < consumers; i++ {
		go Consume("c"+strconv.Itoa(i), ch)
	}
	wgp.Wait() // Wait for all producers to finish.
	close(ch)
	/*
		Vad händer om man byter plats på satserna wgp.Wait() och close(ch) i slutet av main-funktionen?
		Kanalen stängs direkt och producenterna kan inte producera något.
	*/
	/*
		Vad händer om man tar bort satsen close(ch) helt och hållet?
		Då stängs inte kanalen och alla konsumenter fortsätter att vara aktiva.
	*/
	fmt.Println("time:", time.Now().Sub(before))
	/*
		Kan man vara säker på att alla strängar blir utskrivna innan programmet stannar?
		Ja.
		Från golang spec:
		"A send on an unbuffered channel can proceed if a receiver is ready."
		Detta innebär att en konsument alltid kan läsa av meddelandet innan exekvering fortsätter.
	*/
}

// Produce sends n different strings on the channel and notifies wg when done.
func Produce(id string, n int, ch chan<- string, wg *sync.WaitGroup) {
	for i := 0; i < n; i++ {
		//RandomSleep(100) // Simulate time to produce data.
		ch <- id + ":" + strconv.Itoa(i)
	}
	wg.Done()
	/*
		Vad händer om man flyttar close(ch) från main-funktionen och i stället stänger kanalen i slutet av funktionen Produce?
		När en producent är klar kommer kanalen att stängas och övriga kan inte längre skriva till kanalen.
	*/
}

// Consume prints strings received from the channel until the channel is closed.
func Consume(id string, ch <-chan string) {
	for s := range ch {
		fmt.Println(id, "received", s)
		RandomSleep(100) // Simulate time to consume data.
	}
}

// RandomSleep waits for x ms, where x is a random number, 0 ≤ x < n,
// and then returns.
func RandomSleep(n int) {
	time.Sleep(time.Duration(rand.Intn(n)) * time.Millisecond)
}
