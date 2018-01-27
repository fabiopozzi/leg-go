package main

import (
	"fmt"
	"github.com/nsf/termbox-go"
	"github.com/ungerik/go-rss"
)

func main() {
	term_err := termbox.Init()
	if term_err != nil {
		panic(term_err)
	}
	defer termbox.Close()

	channel, err := rss.Read("http://www.ansa.it/sito/ansait_rss.xml")
	if err != nil {
		fmt.Println(err)
	}

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
