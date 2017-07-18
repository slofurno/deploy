package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

var VERSION = "v3.1"

type readyHandler struct {
	mu    sync.Mutex
	ready bool
}

func (s *readyHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.mu.Lock()
	defer s.mu.Unlock()

	switch r.URL.Path {
	case "/healthz":
		w.WriteHeader(200)
		w.Write([]byte(VERSION))
	case "/readyz":
		if s.ready {
			w.WriteHeader(200)
			w.Write([]byte(VERSION))
		} else {
			w.WriteHeader(500)
			w.Write([]byte(VERSION))
		}
	default:
		res := &struct {
			Version string
			Ready   bool
			Url     string
		}{
			Version: VERSION,
			Ready:   s.ready,
			Url:     r.URL.Path,
		}

		b, _ := json.Marshal(res)
		w.Header().Set("Content-Type", "application/json")
		w.Write(b)
	}
}

func main() {
	fmt.Println("vim-go")
	handler := &readyHandler{}

	c := make(chan os.Signal, 2)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		for {
			s := <-c
			switch s {
			case syscall.SIGINT:
				os.Exit(0)
			case syscall.SIGTERM:
				handler.mu.Lock()
				handler.ready = false
				handler.mu.Unlock()
			}
		}
	}()

	go func() {
		time.Sleep(10 * time.Second)
		handler.mu.Lock()
		defer handler.mu.Unlock()
		handler.ready = true
		fmt.Println("ready")
	}()

	http.ListenAndServe(":3008", handler)
}
