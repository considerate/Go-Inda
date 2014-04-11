// http://www.nada.kth.se/~snilsson/concurrency/
package main

import (
	"fmt"
	"runtime"
	"sync"
)

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
}

// This programs demonstrates how a channel can be used for sending and
// receiving by any number of goroutines. It also shows how  the select
// statement can be used to choose one out of several communications.
func main() {
	people := []string{"Anna", "Bob", "Cody", "Dave", "Eva"}
	/* Vad händer om man tar bort bufferten på kanalen match?
	 * Deadlock, eftersom den både försöker läsa från och skriva till
	 * samma kanal.
	 */
	match := make(chan string, 1) // Make room for one unmatched send.
	/* Vad händer om man byter deklarationen wg := new(sync.WaitGroup) mot var wg sync.WaitGroup och parametern wg *sync.WaitGroup mot wg sync.WaitGroup?
	 * Struct:en till wg kommer att kopieras till Seek och när Done körs kommer
	 * det inte påverka wg i main() vilket resulterar i deadlock.
	 */
	wg := new(sync.WaitGroup)
	wg.Add(len(people))
	for _, name := range people {
		/* Vad händer om man tar bort go-kommandot från Seek-anropet i main-funktionen?
		 * Operationen blir icke-samtidig och deterministisk.
		 */
		go Seek(name, match, wg)
	}
	wg.Wait()
	select {
	case name := <-match:
		fmt.Printf("No one received %s’s message.\n", name)
		/* Vad händer om man tar bort default-fallet från case-satsen i main-funktionen?
		 * Deadlock ifall antalet personer är jämnt.
		 */
	default:
		fmt.Printf("No pending message.\n")
		// There was no pending send operation.
	}
}

// Seek either sends or receives, whichever possible, a name on the match
// channel and notifies the wait group when done.
func Seek(name string, match chan string, wg *sync.WaitGroup) {
	select {
	case peer := <-match:
		fmt.Printf("%s sent a message to %s.\n", peer, name)
	case match <- name:
		// Wait for someone to receive my message.
	}
	wg.Done()
}
