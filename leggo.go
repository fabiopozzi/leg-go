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

func RSS(nextSlide func()) (title string, content tview.Primitive) {
	list := tview.NewList()

	feed = parseRss("http://www.ansa.it/sito/ansait_rss.xml")

	list.ShowSecondaryText(false)
	list.SetBorder(true).SetTitle("Elenco notizie")

	// add a list element for each feed title
	shortcut := 'a'
	for _, item := range feed.Items {
		list.AddItem(item.Title, "", shortcut, nil)
		shortcut++
	}

	newsView := tview.NewTextView().
		SetChangedFunc(func() {
			app.Draw()
		})
	newsView.SetBorder(true).SetTitle("Testo notizie")

	list.SetChangedFunc(func(idx int, maintxt string, secondTxt string, shortcut rune) {
		newsView.Clear()
		fmt.Fprintf(newsView, feed.Items[idx].Description)
	})

	return "RSS", tview.NewFlex().
		AddItem(list, 0, 1, true).
		AddItem(newsView, 0, 1, false)
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
