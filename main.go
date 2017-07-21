package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"

	"github.com/acsellers/temple/routing"
)

var (
	port    = flag.Int("port", 8008, "Port for Web")
	apiport = flag.Int("api", 9009, "Port for API")
	cfg     = flag.String("config", "config.json", "Config File")
)

func init() {
	flag.Parse()
}
func main() {
	go web()
	go api()
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
}

func web() {
	w := http.NewServeMux()
	w.Handle("/", routing.MainHandler())
	fmt.Println("Starting web on ", *port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", *port), w))
}

func api() {
	m := http.NewServeMux()
	m.Handle("/", routing.MainHandler())
	fmt.Println("Starting API on ", *apiport)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", *apiport), m))
}
