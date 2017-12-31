package main

import (
	"encoding/xml"
	"html/template"
	"io/ioutil"
	"log"
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
}

func (o Outline) String() string {
	return o.Title
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

type By func(o1, o2 *Outline) bool

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

	title := func(o1, o2 *Outline) bool {
		return strings.ToLower(o1.Title) < strings.ToLower(o2.Title)
	}

	By(title).Sort(doc.Body)

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
</head>
<body>
	<div class="container pt-5 pb-5">
		<h3>{{ .Head.Title }}</h3>
		<hr>
		{{ with .Body }}
			{{ range . }}
				<p>
					<h6>{{ .Title }}</h6>
					<small>
						<a target="_blank" href="{{ .HTMLURL }}">HTML</a> ·
						<a target="_blank" href="{{ .XMLURL }}">XML</a>
					</small>
				</p>
			{{ end }}
		{{ end }}
	</div>
</body>
</html>
`

	t, err := template.New("opml").Parse(tmpl)
	must(err)
	err = t.Execute(os.Stdout, doc)
	must(err)
}
