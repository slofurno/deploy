package main

import (
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

func main() {
	args := os.Args[1:]
	method := args[0]
	uri := args[1]

	r, _ := http.NewRequest(method, uri, nil)
	client := http.Client{}
	res, err := client.Do(r)

	if err != nil {
		log.Fatal(err)
	}

	b := make([]byte, 8)
	for {
		n, err := res.Body.Read(b)
		if err == io.EOF {
			break
		} else if err != nil || n == 0 {
			log.Fatal(err)
		}
		log.Printf("read %d bytes\n", n)
		time.Sleep(100 * time.Millisecond)
	}
}
