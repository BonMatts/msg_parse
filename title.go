package main

// adapted from https://siongui.github.io/2016/05/10/go-get-html-title-via-net-html/
import (
	"errors"
	"io"

	"golang.org/x/net/html"
)

func isTitleElement(n *html.Node) bool {
	return n.Type == html.ElementNode && n.Data == "title"
}

func traverse(n *html.Node) (string, bool) {
	if isTitleElement(n) {
		return n.FirstChild.Data, true
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		result, ok := traverse(c)
		if ok {
			return result, ok
		}
	}

	return "", false
}

//GetHTMLTitle pulls the contents of a pages <title> tag
func GetHTMLTitle(r io.Reader) (string, error) {
	doc, err := html.Parse(r)
	if err != nil {
		return "", err
	}

	if title, ok := traverse(doc); ok {
		return title, nil
	}
	return "", errors.New("Could not parse HTML")

}
