package main

import (
	"flag"
	"log"
	"net/http"
	"os"

	"github.com/hooblei/ytp"
	"github.com/joho/godotenv"
)

var addr = "127.0.0.1:8001"
var proxy *ytp.Proxy

func init() {
	var err error
	var ytHost string
	var ytAuthToken string

	if err = godotenv.Load(); err != nil {
		log.Fatal(err)
	}

	if ytHost = os.Getenv("YT_HOST"); len(ytHost) == 0 {
		log.Fatal("Missing YT_HOST setting")
	}

	if ytAuthToken = os.Getenv("YT_AUTH_TOKEN"); len(ytAuthToken) == 0 {
		log.Fatal("Missing YT_AUTH_TOKEN setting")
	}

	if proxy, err = ytp.New(ytHost, ytAuthToken); err != nil {
		log.Fatal(err)
	}
}

func main() {
	flag.StringVar(&addr, "addr", addr, "The listen address of the YouTrack API proxy.")
	flag.Parse()

	log.Println(addr)
	//http.Handle("/youtrack", proxy)
	if err := http.ListenAndServe(addr, proxy); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
