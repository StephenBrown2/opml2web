package main

import (
	"encoding/xml"
	"io/ioutil"
	"log"
	"os"

	"github.com/davecgh/go-spew/spew"
)

// Head is...
type Head struct {
	Title string `xml:"title"`
}

// Outline is...
type Outline struct {
	Text    string `xml:"text,attr"`
	Title   string `xml:"title,attr"`
	XMLURL  string `xml:"xmlUrl,attr"`
	HTMLURL string `xml:"htmlUrl,attr"`
}

// Document is...
type Document struct {
	XMLName xml.Name  `xml:"opml"`
	Head    Head      `xml:"head"`
	Body    []Outline `xml:"body>outline"`
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	args := os.Args[1:]

	if len(args) == 0 {
		log.Fatal("missing required arguments")
	}

	filename := args[0]

	file, err := ioutil.ReadFile(filename)
	must(err)

	doc := Document{}

	err = xml.Unmarshal(file, &doc)
	must(err)

	spew.Dump(doc)
}
