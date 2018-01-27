package main

import (
	"fmt"
	"github.com/nsf/termbox-go"
	"github.com/ungerik/go-rss"
	"os"
)

func parseRss(url string) *rss.Channel {
	channel, err := rss.Read(url)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	return channel
}

func main() {
	term_err := termbox.Init()
	if term_err != nil {
		panic(term_err)
	}
	defer termbox.Close()

	channel := parseRss("http://www.ansa.it/sito/ansait_rss.xml")
	//fmt.Println(channel.Title)
	fmt.Println("Titoli:")

	for _, item := range channel.Item {
		fmt.Println(item.Title)
	}

loop:
	for {
		switch ev := termbox.PollEvent(); ev.Type {
		case termbox.EventKey:
			switch ev.Key {
			case termbox.KeyEsc:
				break loop
			}
		}
	}
}
