package igtoken

import (
	"errors"
	"io"

	"golang.org/x/net/html"
)

func parseForm(b io.ReadCloser) (csrftToken, action string, err error) {
	var n *html.Node
	if n, err = html.Parse(b); err != nil {
		return
	}
	if n = getFirstNodeByTag(n, "form"); n == nil {
		err = errors.New("No form found from html")
		return
	}
	action = getAttrValue("action", n)

	if n = getFirstNodeByTagAndAttr(n, "input", "name", "csrfmiddlewaretoken"); n == nil {
		err = errors.New("No csrf found from form")
		return
	}
	csrftToken = getAttrValue("value", n)

	return
}

func parseAlert(b io.ReadCloser) (alert string) {
	var err error
	var n *html.Node
	if n, err = html.Parse(b); err != nil {
		return
	}
	if n = getFirstNodeByTagAndAttr(n, "div", "id", "alerts"); n == nil {
		return
	}
	if n = getFirstNodeByTag(n, "p"); n == nil {
		return
	}
	alert = n.FirstChild.Data
	return
}

func getFirstNodeByTag(n *html.Node, tag string) *html.Node {
	if n.Type == html.ElementNode && n.Data == tag {
		return n
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		b := getFirstNodeByTag(c, tag)
		if b != nil {
			return b
		}
	}
	return nil
}

func getFirstNodeByTagAndAttr(n *html.Node, tag, key, val string) *html.Node {
	if n.Type == html.ElementNode && n.Data == tag {
		for _, a := range n.Attr {
			if a.Key == key && a.Val == val {
				return n
			}
		}
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		b := getFirstNodeByTagAndAttr(c, tag, key, val)
		if b != nil {
			return b
		}
	}
	return nil
}

func getAttrValue(key string, n *html.Node) (val string) {
	for _, a := range n.Attr {
		if a.Key == key {
			val = a.Val
			break
		}
	}
	return
}

func joinScopes(scopes []Scope) string {
	if len(scopes) == 1 {
		return string(scopes[0])
	}
	n := len(scopeSep) * (len(scopes) - 1)
	for i := 0; i < len(scopes); i++ {
		n += len(scopes[i])
	}

	b := make([]byte, n)
	bp := copy(b, scopes[0])
	for _, s := range scopes[1:] {
		bp += copy(b[bp:], scopeSep)
		bp += copy(b[bp:], s)
	}
	return string(b)
}
