package lib

import (
	"github.com/antchfx/htmlquery"
	"golang.org/x/net/html"
	"regexp"
)

func FindOne(node *html.Node, xpath string) string {
	nodeOne, err := htmlquery.Query(node, xpath)
	if err != nil {
		panic("can't reconized xpath expression")
	}
	is, attr := XpathAttr(xpath)
	if is {
		return htmlquery.SelectAttr(nodeOne, attr)
	}
	return htmlquery.InnerText(nodeOne)
}

func XpathAttr(xpath string) (flag bool, attr string) {
	flag = false
	reg := regexp.MustCompile(`/@(\w+$)`)
	if reg.MatchString(xpath) {
		flag = true
		submatch := reg.FindStringSubmatch(xpath)
		attr = submatch[1]
	}
	return
}
