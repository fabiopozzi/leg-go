package main

import (
	"fmt"
	"github.com/ungerik/go-rss"
)

func main() {
	channel, err := rss.Read("http://www.ansa.it/sito/ansait_rss.xml")
	if err != nil {
		fmt.Println(err)
	}

	//fmt.Println(channel.Title)
	fmt.Println("Titoli:")

	for _, item := range channel.Item {
		fmt.Println(item.Title)
	}
}
