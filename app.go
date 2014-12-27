package main

import (
	"log"
	"os"
	"time"
)

func main() {
	if len(os.Args) != 2 {
		log.Fatal("No path specified")
	}

	f, err := os.Create(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	defer os.Remove(os.Args[1])

	seconds := 3
	duration := time.Duration(seconds) * time.Second
	written := make(chan int)

	go func() {
		buf := make([]byte, 32*1024)
		total := 0
		defer func() {
			written <- total
		}()

		for {
			n, err := f.Write(buf)
			total += n
			if err != nil {
				return
			}
		}
	}()

	time.Sleep(duration)
	f.Close()
	total := <-written

	amount := float64(total) / 1024 / 1024
	rate := amount / float64(seconds)
	log.Println("Written:", amount, "mb")
	log.Println("Write rate:", rate, "mb/s")
}
