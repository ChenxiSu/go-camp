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
	errors "github.com/pkg/errors"
)

func serve(ctx context.Context, addr string, handler http.Handler) error {
	srv := http.Server{
		Addr: addr,
		Handler: handler,
	}

	go func() {
		select {
		case <-ctx.Done():
			fmt.Println("\nClosing server!")
			srv.Close()
		}
	}()

	// Start http server
	return srv.ListenAndServe()
}

func serveApp(ctx context.Context) error{
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(resp http.ResponseWriter, req *http.Request) {
		fmt.Fprintln(resp, "Hello World!")
	})
	fmt.Println("Starting app server at 127.0.0.1:8080")
	return serve(ctx, "127.0.0.1:8080", mux)
}

func serveDebug(ctx context.Context) error{
	fmt.Println("Starting debug server at 127.0.0.1:8001")
	return serve(ctx, "127.0.0.1:8001", http.DefaultServeMux)
}

func main() {
	var TerminateErr = errors.New("Terminated by signal")
	g, ctx := errgroup.WithContext(context.Background())

	sig_chan := make(chan os.Signal, 1)
	signals := []os.Signal{syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGINT}
	signal.Notify(sig_chan, signals...)

	// init application server
	g.Go(func() error {
		return serveApp(ctx)
	})
	// init debug server
	g.Go(func() error {
		return serveDebug(ctx)
	})
	// go routine that handles termination signals
	g.Go(func() error {
		for {
			select {
				case <-ctx.Done():
					return ctx.Err()
				case sig := <-sig_chan:
					fmt.Printf("\n\nSignal received: %v \n", sig)
					message := fmt.Sprintf("Signal received: %v", sig)
					return errors.Wrapf(TerminateErr, message)
			}
		}
	})

	// Wait for all HTTP fetches to complete.
	time.Sleep(time.Second * 1)
	fmt.Println("awaiting signal")
	if err := g.Wait(); err != nil && !errors.Is(err, TerminateErr) {
		log.Fatal(err)
	}

	fmt.Println("All servers have been properly closed.")
}