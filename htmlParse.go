package main

import (
	"fmt"
	"os"
	"strings"

	"golang.org/x/net/html"
)

type label struct {
	text     string
	forInput string
}

type printable struct {
	paragraphs []string
	labels     []label
}

func getPrintableContent(htmlFileName string) (*printable, error) {
	f, err := os.Open(htmlFileName)
	if err != nil {
		return nil, fmt.Errorf("fail to open %v: %v", f, err)
	}
	defer f.Close()
	doc, err := html.Parse(f)
	if err != nil {
		return nil, fmt.Errorf("fail to parse %v: %v", f, err)
	}
	p := newPrintable()
	extractIntoPrintable(doc, p)
	return p, nil
}

func newPrintable() *printable {
	return &printable{
		paragraphs: make([]string, 0),
		labels:     make([]label, 0),
	}
}

func extractIntoPrintable(n *html.Node, p *printable) {
	if isParagraph(n) {
		p.paragraphs = append(p.paragraphs, strings.TrimSpace(buildParagraph(n)))
	} else if isLabel(n) {
		p.labels = append(p.labels, buildLabel(n))
	} else {
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			extractIntoPrintable(c, p)
		}
	}
}

func isParagraph(n *html.Node) bool {
	return n.Type == html.ElementNode && n.Data == "p"
}

func buildParagraph(n *html.Node) string {
	var p string
	if n.Type == html.TextNode {
		trimmed := make([]rune, 0)
		var prev rune
		for _, c := range n.Data {
			if !(c == ' ' && prev == ' ' || c == '\n') {
				trimmed = append(trimmed, c)
			}
			prev = c
		}
		p += string(trimmed)
	}
	if n.Type == html.ElementNode {
		if n.Data != "p" {
			p += "<" + n.Data + ">"
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			p += buildParagraph(c)
		}
		if n.Data != "p" {
			p += "</" + n.Data + ">"
		}
	}

	return p
}

func isLabel(n *html.Node) bool {
	return n.Type == html.ElementNode && n.Data == "label"
}

func buildLabel(n *html.Node) label {
	var l label
	for _, atr := range n.Attr {
		if atr.Key == "for" {
			l.forInput = atr.Val
		}
	}

	if n.Type == html.ElementNode {
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			if c.Type == html.TextNode {
				l.text = c.Data
				break
			}
		}
	}

	return l
}
