package main

import (
	"encoding/xml"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"sort"
	"strings"
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
	Feed    Feed
}

// OPML is...
type OPML struct {
	XMLName xml.Name  `xml:"opml"`
	Head    Head      `xml:"head"`
	Body    []Outline `xml:"body>outline"`
}

// Image is...
type Image struct {
	URL   string `xml:"url"`
	Title string `xml:"title"`
	Link  string `xml:"link"`
	HREF  string `xml:"href,attr"`
}

// Feed is...
type Feed struct {
	XMLName     xml.Name      `xml:"rss"`
	Description template.HTML `xml:"channel>description"`
	Image       Image         `xml:"channel>image"`
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}

// By does...
type By func(o1, o2 *Outline) bool

// Sort does...
func (by By) Sort(outlines []Outline) {
	os := &outlineSorter{
		outlines: outlines,
		by:       by,
	}
	sort.Sort(os)
}

// Swap is part of sort.Interface.
func (s *outlineSorter) Swap(i, j int) {
	s.outlines[i], s.outlines[j] = s.outlines[j], s.outlines[i]
}

// Less is part of sort.Interface. It is implemented by calling the "by" closure in the sorter.
func (s *outlineSorter) Less(i, j int) bool {
	return s.by(&s.outlines[i], &s.outlines[j])
}

func (s *outlineSorter) Len() int {
	return len(s.outlines)
}

type outlineSorter struct {
	outlines []Outline
	by       func(o1, o2 *Outline) bool // Closure used in the Less method.
}

// ParseOPML does...
func ParseOPML(input []byte) (OPML, error) {
	opml := OPML{}
	err := xml.Unmarshal(input, &opml)
	if err != nil {
		log.Println(opml)
	}
	if len(opml.Body) == 1 {
		// Some feed exporters (like Pocket Casts) have an extra <outline text="feeds"> element
		type Feeds struct {
			Feeds []Outline `xml:"body>outline>outline"`
		}
		feeds := Feeds{}
		err := xml.Unmarshal(input, &feeds)
		if err == nil {
			opml.Body = feeds.Feeds
		}
	}
	return opml, err
}

// ParseFeed does...
func ParseFeed(input []byte) (Feed, error) {
	feed := Feed{}
	err := xml.Unmarshal(input, &feed)
	if err != nil {
		log.Println(feed)
	}
	return feed, err
}

// TitleSorter does...
func TitleSorter(o1, o2 *Outline) bool {
	return strings.ToLower(o1.Title) < strings.ToLower(o2.Title)
}

func main() {
	args := os.Args[1:]

	if len(args) == 0 {
		log.Fatal("missing required arguments")
	}

	filename := args[0]

	file, err := ioutil.ReadFile(filename)
	must(err)

	opml, err := ParseOPML(file)
	must(err)

	By(TitleSorter).Sort(opml.Body)
	log.Printf("%d feeds to process", len(opml.Body))
	for index, outline := range opml.Body {
		title := outline.Title
		if title == "" {
			title = outline.Text
		}
		if outline.XMLURL == "" {
			log.Printf("No feed URL for %s\n", title)
			continue
		}
		log.Printf("Get feed %d: %s - %s\n", index+1, title, outline.XMLURL)
		res, err := http.Get(outline.XMLURL)
		must(err)

		body, err := ioutil.ReadAll(res.Body)
		must(err)
		res.Body.Close()

		feed, err := ParseFeed(body)
		if err != nil {
			log.Println("Error parsing feed. Continuing.")
			continue
		}

		opml.Body[index].Feed = feed
	}

	tmpl := `
	<!doctype html>
	<html class="no-js" lang="">
	<head>
		<meta charset="utf-8">
		<meta http-equiv="x-ua-compatible" content="ie=edge">
		<title>{{ .Head.Title }}</title>
		<meta name="description" content="">
		<meta name="viewport" content="width=device-width, initial-scale=1, shrink-to-fit=no">
		<link rel="stylesheet" href="https://maxcdn.bootstrapcdn.com/bootstrap/4.0.0-beta.3/css/bootstrap.min.css" integrity="sha384-Zug+QiDoJOrZ5t4lssLdxGhVrurbmBWopoEl+M6BdEfwnCJZtKxi1KgxUyJq13dy" crossorigin="anonymous">
		<style>
			.truncate {
				overflow: hidden;
				white-space: nowrap;
				text-overflow: ellipsis;
			}
		</style>
	</head>
	<body>
		<div class="container pt-5 pb-5">
			<h3>
				{{ .Head.Title }}
			</h3>
			<hr>
			{{ with .Body }}
				{{ range . }}
					<div class="media">
						<img width="100" class="img-fluid mr-3" src="{{ .Feed.Image.HREF }}" alt="{{ .Feed.Image.Title }}">
						<div class="media-body">
							<h5 class="mt-0">{{ .Title }}</h5>
							<p>{{ .Feed.Description }}</p>
							<small class="text-muted">
							{{ with .HTMLURL }}
								<a target="_blank" href="{{ . }}">Web</a> ·
							{{ end }}
								<a target="_blank" href="{{ .XMLURL }}">Feed</a>
							</small>
						</div>
					</div>
					<br>
				{{ end }}
			{{ end }}
		</div>
	</body>
	</html>
	`

	t, err := template.New("opml").Parse(tmpl)
	must(err)
	err = t.Execute(os.Stdout, opml)
	must(err)
}
