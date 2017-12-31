package main

import (
	"encoding/xml"
	"html/template"
	"io/ioutil"
	"log"
	"os"
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

	tmpl := `
<!doctype html>
<html class="no-js" lang="">

<head>
    <meta charset="utf-8">
    <meta http-equiv="x-ua-compatible" content="ie=edge">
    <title>{{ .Head.Title }}</title>
    <meta name="description" content="">
    <meta name="viewport" content="width=device-width, initial-scale=1, shrink-to-fit=no">

    <link rel="stylesheet" href="https://maxcdn.bootstrapcdn.com/bootstrap/4.0.0-beta.3/css/bootstrap.min.css" integrity="sha384-Zug+QiDoJOrZ5t4lssLdxGhVrurbmBWopoEl+M6BdEfwnCJZtKxi1KgxUyJq13dy"
        crossorigin="anonymous">
</head>

<body>

</body>

</html>
`

	t, err := template.New("opml").Parse(tmpl)
	must(err)
	err = t.Execute(os.Stdout, doc)
	must(err)
}
