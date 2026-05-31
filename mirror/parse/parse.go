// Package parse provides utilities for extracting data from HTML documents.
package parse

import (
	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

// ExtractHTML recursively traverses an HTML node tree and extracts URLs from relevant tags.
// Discovered URLs are appended to the provided links slice.
func ExtractHTML(n *html.Node, links *[]string) {
	if n.Type == html.ElementNode {
		switch n.DataAtom {
		case atom.A, atom.Link:
			for _, a := range n.Attr {
				if a.Key == "href" {
					*links = append(*links, a.Val)
				}
			}
		case atom.Img, atom.Script, atom.Video, atom.Audio, atom.Source, atom.Iframe:
			for _, a := range n.Attr {
				if a.Key == "src" {
					*links = append(*links, a.Val)
				}
			}
		}
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		ExtractHTML(c, links)
	}
}
