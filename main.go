package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

var VERSION = "v0.5.1"

type readyHandler struct {
	mu    sync.Mutex
	ready bool
	alive bool
}

func (s *readyHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Printf("got req %s on %s\n", r.URL.Path, os.Getenv("HOSTNAME"))
	switch r.URL.Path {
	case "/healthz":
		if s.alive {
			w.WriteHeader(200)
			w.Write([]byte(VERSION))
		} else {
			w.WriteHeader(500)
			w.Write([]byte(VERSION))
		}
	case "/readyz":
		if s.ready {
			w.WriteHeader(200)
			w.Write([]byte(VERSION))
		} else {
			w.WriteHeader(500)
			w.Write([]byte(VERSION))
		}
	case "/random":
		f, err := os.Open("/dev/random")
		if err != nil {
			fmt.Println(err)
		}
		n, err := io.CopyN(w, f, 1<<20)
		if err != nil {
			fmt.Println(err)
		}

		fmt.Printf("sent %d random bytes\n", n)
	case "/kill":
		s.mu.Lock()
		defer s.mu.Unlock()
		s.alive = false
		w.Write([]byte("ok"))
	case "/sleep":
		time.Sleep(30 * time.Second)
		w.Write([]byte("ok"))

	default:
		s.mu.Lock()
		defer s.mu.Unlock()
		res := &struct {
			Version  string
			Ready    bool
			Url      string
			Hostname string
		}{
			Version:  VERSION,
			Ready:    s.ready,
			Url:      r.URL.Path,
			Hostname: os.Getenv("HOSTNAME"),
		}

		b, _ := json.Marshal(res)
		w.Header().Set("Content-Type", "application/json")
		w.Write(b)
	}
}

func main() {
	fmt.Println("vim-go")
	handler := &readyHandler{ready: true, alive: true}

	c := make(chan os.Signal, 2)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)

	h := &http.Server{Addr: ":3008", Handler: handler}

	done := make(chan struct{}, 1)

	shutdown := func() {
		close(done)
	}

	go func() {
		for s := range c {
			switch s {
			case syscall.SIGINT:
				os.Exit(0)
			case syscall.SIGTERM:
				handler.mu.Lock()
				handler.ready = false
				handler.mu.Unlock()
				shutdown()
			}
		}
	}()

	go func() {
		switch err := h.ListenAndServe(); err {
		case http.ErrServerClosed:
			log.Print("listener closed")
		default:
			log.Fatal(err)
		}
	}()

	<-done
	t0 := time.Now()
	fmt.Println("starting shutdown...")
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	if err := h.Shutdown(ctx); err != nil {
		fmt.Printf("error during shutdown: %v\n", err)
	}
	fmt.Printf("shutdown after %d seconds\n", time.Now().Sub(t0).Seconds())
}
