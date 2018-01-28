package main

import (
	"fmt"
	"github.com/mmcdole/gofeed"
	"github.com/nsf/termbox-go"
	"os"
)

type configT struct {
	startLine int
	numRows   int
	headline  string
}

var cfg configT

var feed *gofeed.Feed
var curRow int

func parseRss(url string) *gofeed.Feed {
	fp := gofeed.NewParser()
	feed, err := fp.ParseURL(url)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	return feed
}

func printTitle(title string) {
	runes := []rune(title)
	bg := termbox.AttrBold | termbox.ColorRed
	for i := 0; i < len(runes); i++ {
		fg := termbox.AttrBold | termbox.ColorWhite
		termbox.SetCell(i+1, 0, runes[i], fg, bg)
	}
}

func rowIncrement() {
	if curRow < cfg.numRows {
		curRow++
	}
	drawAll()
}

func rowDecrement() {
	if curRow > cfg.startLine {
		curRow--
	}
	drawAll()
}

func showArticle() {
	w, _ := termbox.Size()
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
	item := feed.Items[curRow-cfg.startLine]
	printTitle(item.Title)
	bg := termbox.AttrBold
	fg := termbox.AttrBold | termbox.ColorGreen
	y := cfg.startLine
	runes := []rune(item.Description)
	i := 0
	for x := 0; x < len(runes); x++ {
		if i == (w - 5) {
			y++
			i = 0
		}
		termbox.SetCell(i+1, y, runes[x], fg, bg)
		i++
	}
	termbox.Flush()
}

func printNews() {
	bg := termbox.ColorDefault
	y := cfg.startLine
	// TODO: use method to generate output string from item.Title
	for _, item := range feed.Items {
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

	printTitle(cfg.headline)
	printNews()

	termbox.Flush()
}

func main() {
	err := termbox.Init()
	if err != nil {
		panic(err)
	}
	defer termbox.Close()

	feed = parseRss("http://www.ansa.it/sito/ansait_rss.xml")

	// init cfg
	cfg.headline = "Titoli:"
	cfg.startLine = 2
	cfg.numRows = len(feed.Items) + cfg.startLine

	curRow = cfg.startLine
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
			case termbox.KeyArrowRight:
				showArticle()
			case termbox.KeyArrowLeft:
				drawAll()
			}
		case termbox.EventResize:
			drawAll()
		}
	}
}
