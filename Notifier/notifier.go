package main

import (
	"fmt"
	"time"
)

func Remind(text string, pause time.Duration) {
	tickChannel := time.Tick(pause)
	for {
		select {
		case now := <-tickChannel:
			fmt.Printf(text, now.Format("15:04"))
		}
	}
}

func main() {
	go Remind("Klockan är %s: Dags att äta\n", 3*time.Hour)
	go Remind("Klockan är %s: Dags att arbeta\n", 8*time.Hour)
	go Remind("Klockan är %s: Dags att sova\n", 24*time.Hour)
	select {}
}
