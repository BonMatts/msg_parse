package main

import (
	"net/http"
	"regexp"

	"github.com/mvdan/xurls"
)

//Msg parsed message struct
type Msg struct {
	txt   string
	Ments *[]string `json:"mentions, omitempty"`
	Emos  *[]string `json:"emoticons, omitempty"`
	Lnks  *[]link   `json:"links, omitempty"`
}

type link struct {
	URL   string  `json:"url"`
	Title *string `json:"title"`
}

//NewMsg  message constructor
func NewMsg(text string) *Msg {
	msg := &Msg{txt: text}
	msg.txt = text
	msg.parseMsg()
	return msg

}

func (m *Msg) parseMsg() {
	m.Ments = m.find(`@(\w*)`)

	m.Emos = m.find(`\((\w{1,15})\)`)

	m.findLink()
}

func (m *Msg) findLink() {
	urls := xurls.Strict.FindAllString(m.txt, -1)
	lnks := []link{}
	if urls != nil {

		for _, url := range urls {
			lnk := link{URL: url}
			err := lnk.getTitle()
			if err != nil {
				//todo: log the error "could not parse html"
			}
			lnks = append(lnks, lnk)

		}
	}
	m.Lnks = &lnks
}

func (l *link) getTitle() error {
	resp, err := http.Get(l.URL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	title, err := GetHTMLTitle(resp.Body)

	if err != nil {
		return err

	}
	l.Title = &title
	return nil

}

func (m *Msg) find(pattern string) *[]string {
	re := regexp.MustCompile(pattern)
	var mtchs []string
	for _, match := range re.FindAllStringSubmatch(m.txt, -1) {
		mtchs = append(mtchs, match[1])
	}
	return &mtchs
}
