package main

import (
	"fmt"
	"log"

	"github.com/antchfx/htmlquery"
)

func main() {
	doc, err := htmlquery.LoadURL("https://github.com/RocketChat/Rocket.Chat.Electron/releases")

	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	list := htmlquery.Find(doc, "//a")

	fmt.Printf("%v\n", list)

}
