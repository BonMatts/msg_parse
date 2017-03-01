package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

type resp struct {
	Mentions  []string `json:"mentions"`
	Emoticons []string `json:"emoticons"`
	Links     []lnk    `json:"links"`
}

type lnk struct {
	URL   string
	Title string
}

var server *httptest.Server

func init() {
	server = httptest.NewServer(Handlers())

}

func TestParseMsg(t *testing.T) {
	//just mentions
	m := unmarshRes(sendReq("@chris you around?", t))

	if !reflect.DeepEqual(m.Mentions, []string{"chris"}) {
		t.Errorf("expected mentions to equal [chris], got %v", m.Mentions)
	}

	//just emoticons
	e := unmarshRes(sendReq("Good morning! (megusta) (coffee)", t))

	if !reflect.DeepEqual(e.Emoticons, []string{"megusta", "coffee"}) {
		t.Errorf("expected emoticons to equal [megusta, coffee], got %v", e.Emoticons)
	}

	//just links
	l := unmarshRes(sendReq("Olympics are starting soon; http://www.nbcolympics.com", t))
	exp := []lnk{lnk{URL: "http://www.nbcolympics.com", Title: "2018 PyeongChang Olympic Games | NBC Olympics"}}
	if !reflect.DeepEqual(l.Links, exp) {
		t.Errorf("expected links to equal %v, got %v", exp, l.Links)
	}

	//all three

	c := unmarshRes(sendReq("@bob @john (success) such a cool feature; https://twitter.com/jdorfman/status/430511497475670016", t))
	lk := lnk{URL: "https://twitter.com/jdorfman/status/430511497475670016", Title: `Justin Dorfman on Twitter: &quot;nice @littlebigdetail from @HipChat (shows hex colors when pasted in chat). http://t.co/7cI6Gjy5pq&quot;`}
	r := resp{
		Mentions:  []string{"bob", "john"},
		Emoticons: []string{"success"},
		Links:     []lnk{lk},
	}
	if !reflect.DeepEqual(c, r) {
		t.Errorf("expected links to equal %v, got %v", r, c)
	}
}

func sendReq(input string, t *testing.T) *http.Response {
	req, err := http.NewRequest("POST", server.URL, bytes.NewBufferString(input))
	res, err := http.DefaultClient.Do(req)

	if err != nil {
		t.Error(err) //error sending request
	}
	if res.StatusCode != 200 {
		t.Errorf("Success expected: %d", res.StatusCode)
	}
	//this bit is for debugging, dumps response into stout
	//responseDump, err := httputil.DumpResponse(res, true)
	//if err != nil {
	//	fmt.Println(err)
	//}
	//fmt.Println(string(responseDump))
	return res
}

func unmarshRes(res *http.Response) resp {
	b := resp{}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println("Error: %s\n", err) //error reading body
	}
	err = json.Unmarshal(body, &b)
	if err != nil {
		fmt.Println(err) //error unmarshalling body
	}
	//print struct for debugging
	//fmt.Printf("%+v\n", b)

	return b
}

//Input: "@bob @john (success) such a cool feature; https://twitter.com/jdorfman/status/430511497475670016"
//Return:
//{
//"mentions": [
//"bob",
//"john"
//],
//"emoticons": [
//"success"
//],
//"links": [
//{
//"url": "https://twitter.com/jdorfman/status/430511497475670016",
//"title": "Justin Dorfman on Twitter: &quot;nice @littlebigdetail from @HipChat (shows hex colors when pasted in chat). http://t.co/7cI6Gjy5pq&quot;"
//}
//]
//}
