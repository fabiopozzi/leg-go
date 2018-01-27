package main

import (
	"fmt"
	"github.com/nsf/termbox-go"
	"github.com/ungerik/go-rss"
	"os"
)

const headline = "Titoli:"

func parseRss(url string) *rss.Channel {
	channel, err := rss.Read(url)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	return channel
}

func print_title() {
	bg := termbox.AttrBold | termbox.ColorRed
	for i := 0; i < len(headline); i++ {
		fg := termbox.AttrBold | termbox.ColorWhite
		termbox.SetCell(i+1, 0, rune(headline[i]), fg, bg)
	}
}

func draw_all() {
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)

	print_title()

	termbox.Flush()
}

func main() {
	err := termbox.Init()
	if err != nil {
		panic(err)
	}
	defer termbox.Close()

	//channel := parseRss("http://www.ansa.it/sito/ansait_rss.xml")

	draw_all()
loop:
	for {
		switch ev := termbox.PollEvent(); ev.Type {
		case termbox.EventKey:
			switch ev.Key {
			case termbox.KeyEsc:
				break loop
			}
		case termbox.EventResize:
			draw_all()
		}
	}
	//	fmt.Println(channel.Title)
	//	for _, item := range channel.Item {
	//		fmt.Println(item.Title)
	//	}
}
