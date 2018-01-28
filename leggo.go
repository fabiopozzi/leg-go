package main

import (
	"fmt"
	"github.com/nsf/termbox-go"
	"github.com/ungerik/go-rss"
	"os"
)

const headline = "Titoli:"

var channel *rss.Channel
var curRow int

func parseRss(url string) *rss.Channel {
	c, err := rss.Read(url)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	return c
}

func printTitle() {
	bg := termbox.AttrBold | termbox.ColorRed
	for i := 0; i < len(headline); i++ {
		fg := termbox.AttrBold | termbox.ColorWhite
		termbox.SetCell(i+1, 0, rune(headline[i]), fg, bg)
	}
}

func rowIncrement() {
	// TODO: make starting line a const (remove magic number)
	numRows := len(channel.Item) + 2
	if curRow < numRows {
		curRow++
	}
	drawAll()
}

func rowDecrement() {
	// TODO: make starting line a const (remove magic number)
	if curRow > 2 {
		curRow--
	}
	drawAll()
}

func printNews() {
	bg := termbox.ColorDefault
	y := 2
	for _, item := range channel.Item {
		runes := []rune(item.Title)
		for i := 0; i < len(runes); i++ {
			if y == curRow {
				bg = termbox.ColorRed | termbox.AttrBold
			} else {
				bg = termbox.AttrBold
			}
			fg := termbox.AttrBold | termbox.ColorGreen
			termbox.SetCell(i+1, y, runes[i], fg, bg)
		}
		y++
	}
}

func drawAll() {
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)

	printTitle()
	printNews()

	termbox.Flush()
}

func main() {
	err := termbox.Init()
	if err != nil {
		panic(err)
	}
	defer termbox.Close()

	channel = parseRss("http://www.ansa.it/sito/ansait_rss.xml")
	curRow = 2

	drawAll()
loop:
	for {
		switch ev := termbox.PollEvent(); ev.Type {
		case termbox.EventKey:
			switch ev.Key {
			case termbox.KeyEsc:
				break loop
			case termbox.KeyArrowDown:
				rowIncrement()
			case termbox.KeyArrowUp:
				rowDecrement()
			}
		case termbox.EventResize:
			drawAll()
		}
	}
}
