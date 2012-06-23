package main

import (
	"flag"
	"net/http"
	"net/url"
	"fmt"
)

var msg string
var clientID string
var host string

func main() {
	flag.StringVar(&msg, "m", "Hello", "Message")
	flag.StringVar(&clientID, "c", "0", "ClientID")
	flag.StringVar(&host, "host", "gophernotify.appspot.com", "Host")
	flag.Parse()
	urlStr := fmt.Sprintf("http://%s/%s/post", host, clientID)
	http.PostForm(urlStr, url.Values{"message": {msg}})
}
