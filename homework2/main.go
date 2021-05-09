package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"golang.org/x/sync/errgroup"
)

func serve(addr string, handler http.Handler, stop <-chan struct{} ) error {
	srv := http.Server{
		Addr: addr,
		Handler: handler,
	}

	go func() {
		<- stop
		fmt.Println("\nClosing server!")
		srv.Close()
	}()

	// Start http server
	return srv.ListenAndServe()
}

func serveApp(stop <-chan struct{}) error{
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(resp http.ResponseWriter, req *http.Request) {
		fmt.Fprintln(resp, "Hello World!")
	})
	fmt.Println("Starting app server at 127.0.0.1:8080")
	return serve("127.0.0.1:8080", mux, stop)
}

func serveDebug(stop <-chan struct{}) error{
	fmt.Println("Starting debug server at 127.0.0.1:8001")
	return serve("127.0.0.1:8001", http.DefaultServeMux, stop)
}

func main() {
	g, ctx := errgroup.WithContext(context.Background())
	stop_chan := make(chan struct{})

	sig_chan := make(chan os.Signal, 1)
	signals := []os.Signal{syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGINT}
	signal.Notify(sig_chan, signals...)

	// init application server
	g.Go(func() error {
		return serveApp(stop_chan)
	})
	// init debug server
	g.Go(func() error {
		return serveDebug(stop_chan)
	})
	// go routine that handles termination signals
	g.Go(func() error {
		for {
			select {
				case <-ctx.Done():
					return ctx.Err()
				case sig := <-sig_chan:
					fmt.Printf("\n\nSignal received: %v \n", sig)
					close(stop_chan)
			}
		}
	})

	// Wait for all HTTP fetches to complete.
	time.Sleep(time.Second * 1)
	fmt.Println("awaiting signal")
	if err := g.Wait(); err != nil && err != http.ErrServerClosed{
		log.Fatal(err)
	}

	fmt.Println("All servers have been properly closed.")
}