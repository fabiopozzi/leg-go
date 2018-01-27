package main

import (
	"fmt"
	"github.com/nsf/termbox-go"
	"github.com/ungerik/go-rss"
	"os"
)

const headline = "Titoli:"

var channel *rss.Channel

func parseRss(url string) *rss.Channel {
	c, err := rss.Read(url)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	return c
}

func print_title() {
	bg := termbox.AttrBold | termbox.ColorRed
	for i := 0; i < len(headline); i++ {
		fg := termbox.AttrBold | termbox.ColorWhite
		termbox.SetCell(i+1, 0, rune(headline[i]), fg, bg)
	}
}

func print_news() {
	bg := termbox.AttrBold | termbox.ColorGreen
	y := 2
	for _, item := range channel.Item {
		for i := 0; i < len(item.Title); i++ {
			fg := termbox.AttrBold | termbox.ColorWhite
			termbox.SetCell(i+1, y, rune(item.Title[i]), fg, bg)
		}
		y++
	}
}

func draw_all() {
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)

	print_title()
	print_news()

	termbox.Flush()
}

func main() {
	err := termbox.Init()
	if err != nil {
		panic(err)
	}
	defer termbox.Close()

	channel = parseRss("http://www.ansa.it/sito/ansait_rss.xml")

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
}
