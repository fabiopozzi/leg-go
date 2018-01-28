package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/gdamore/tcell"
	"github.com/mmcdole/gofeed"
	"github.com/rivo/tview"
)

type Slide func(nextSlide func()) (title string, content tview.Primitive)

var app = tview.NewApplication()
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

/*
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
	// TODO: refactor insieme a printNews per estrarre parte rendering.
	w, _ := termbox.Size()
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
	item := feed.Items[curRow-cfg.startLine]
	printTitle(item.Title)
	bg := termbox.AttrBold
	fg := termbox.AttrBold | termbox.ColorWhite
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

*/

func RSS(nextSlide func()) (title string, content tview.Primitive) {
	table := tview.NewTable().
		SetFixed(1, 1)

	list := tview.NewList()

	showDescription := func() {
		// TODO: mostrare il testo quando premi enter.
	}

	feed = parseRss("http://www.ansa.it/sito/ansait_rss.xml")

	list.ShowSecondaryText(false)

	// add a list element for each feed title
	shortcut := 'a'
	for _, item := range feed.Items {
		list.AddItem(item.Title, "", shortcut, showDescription)
		shortcut++
	}

	return "RSS", tview.NewFlex().
		AddItem(tview.NewFlex().
			SetDirection(tview.FlexRow).
			AddItem(list, 10, 1, true).
			AddItem(table, 0, 1, true), 0, 1, true)
}

func main() {
	// all application tabs (To be implemented)
	slides := []Slide{
		RSS,
	}

	// The bottom row has some info on where we are.
	info := tview.NewTextView().
		SetDynamicColors(true).
		SetRegions(true).
		SetWrap(false)

	// Create pages for all slides
	curSlide := 0
	info.Highlight(strconv.Itoa(curSlide))
	pages := tview.NewPages()
	previousSlide := func() {
		curSlide = (curSlide - 1 + len(slides)) % len(slides)
		info.Highlight(strconv.Itoa(curSlide))
		pages.SwitchToPage(strconv.Itoa(curSlide))
	}
	nextSlide := func() {
		curSlide = (curSlide + 1) % len(slides)
		info.Highlight(strconv.Itoa(curSlide))
		pages.SwitchToPage(strconv.Itoa(curSlide))
	}

	for index, slide := range slides {
		title, primitive := slide(nextSlide)
		pages.AddPage(strconv.Itoa(index), primitive, true, index == curSlide)
		fmt.Fprintf(info, `%d ["%d"][darkcyan]%s[white][""]  `, index+1, index, title)
	}

	layout := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(pages, 0, 1, true).
		AddItem(info, 1, 1, false)

	// keybindings to move between slides and exit.
	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyCtrlN {
			nextSlide()
		} else if event.Key() == tcell.KeyCtrlP {
			previousSlide()
		} else if event.Key() == tcell.KeyEsc {
			app.Stop()
		}
		return event
	})

	if err := app.SetRoot(layout, true).SetFocus(layout).Run(); err != nil {
		panic(err)
	}
}
