package main

import "github.com/virteman/webview"

func main() {
	// Open wikipedia in a 800x600 resizable window
	/*
		webview.Open("Minimal webview example",
			"https://en.m.wikipedia.org/wiki/Main_Page", 800, 600, true)
	*/
	webview.Open("Minimal webview example",
    "http://10.8.230.114:9898/sambshare/UI/web20e83849-9d84-4ccb-bd21-c30d82b7a5a1.gif", 800, 600, 3)
}
