package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	port := flag.String("p", "8081", "server listening port")
	directory := flag.String("d", ".", "files directory")
	flag.Parse()

	fs := http.FileServer(http.Dir(*directory))

	http.Handle("/", fs)

	errs := make(chan error, 2)

	go func() {
		log.Printf("Serving %s on HTTP port: %s\n", *directory, *port)
		log.Fatal(http.ListenAndServe(":"+*port, nil))
	}()

	go func() {
		c := make(chan os.Signal)
		signal.Notify(c, syscall.SIGINT)
		errs <- fmt.Errorf("%s", <-c)
	}()

	err := <-errs
	log.Printf(fmt.Sprintf("file server terminated: %s", err))


}
