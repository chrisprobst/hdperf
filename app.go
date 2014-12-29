package main

import (
	"io"
	"log"
	"os"
	"time"
)

func main() {
	if len(os.Args) != 2 {
		log.Fatal("No path specified")
	}

	// Create new benchmark file
	f, err := os.Create(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	defer os.Remove(os.Args[1])

	// Create garbage
	buf := make([]byte, 32*1024)
	for i := range buf {
		buf[i] = byte(i % 255)
	}

	// Write file
	t := time.Now()
	total := 1024 * 1024 * 1024
	for i := 0; i < total; {
		l := len(buf)
		if total-i < l {
			l = total - i
		}

		n, err := f.Write(buf[:l])
		i += n
		if err != nil {
			return
		}
	}
	f.Close()
	d := time.Now().Sub(t)
	mb := total / 1024 / 1024
	rate := float64(mb) / d.Seconds()
	log.Println("Written:", mb, "mb")
	log.Println("Write rate:", rate, "mb/s")

	// Open benchmark file
	f, err = os.Open(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}

	// Read file
	t = time.Now()
	i := 0
	for {
		n, err := f.Read(buf)
		i += n
		if err == io.EOF {
			break
		} else if err != nil {
			log.Fatalln(err)
		}
	}
	f.Close()
	d = time.Now().Sub(t)
	mb = i / 1024 / 1024
	rate = float64(mb) / d.Seconds()
	log.Println("Read:", mb, "mb")
	log.Println("Read rate:", rate, "mb/s")
}
