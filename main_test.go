package main

import "testing"

func TestBy(t *testing.T) {
	opml := OPML{}
	opml.Body = append(opml.Body, Outline{
		Title: "z",
	})
	opml.Body = append(opml.Body, Outline{
		Title: "b",
	})
	opml.Body = append(opml.Body, Outline{
		Title: "a",
	})

	By(TitleSorter).Sort(opml.Body)

	if opml.Body[0].Title != "a" {
		t.Fatalf("expected %s, got %s", "a", opml.Body[0].Title)
	}
}

func TestParseOPML(t *testing.T) {
	input := []byte(`
	<opml>
		<head>
			<title>foo</title>
		</head>
		<body>
			<outline title="bar"></outline>
			<outline title="baz"></outline>
		</body>
	</opml>
	`)

	output, _ := ParseOPML(input)

	if output.Head.Title != "foo" {
		t.Fatalf("expected %s, got %s", "foo", output.Head.Title)
	}

	if output.Body[0].Title != "bar" {
		t.Fatalf("expected %s, got %s", "bar", output.Body[0].Title)
	}

	if output.Body[1].Title != "baz" {
		t.Fatalf("expected %s, got %s", "baz", output.Body[1].Title)
	}
}

func TestParseFeed(t *testing.T) {
	input := []byte(`
	<rss>
		<channel>
			<description>...</description>
			<itunes:image href="image.jpg"></itunes:image>
		</channel>
	</rss>
	`)

	output, _ := ParseFeed(input)

	if output.Image.HREF != "image.jpg" {
		t.Fatalf("expected %s, got %s", "image.jpg", output.Image.HREF)
	}
}
