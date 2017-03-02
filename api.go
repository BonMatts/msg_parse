package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"
)

//Handlers well, it's the router, my friend.
func Handlers() *mux.Router {
	r := mux.NewRouter()

	r.HandleFunc("/", parseMsgHandler).Methods("POST") //there is just one endpoint!

	return r
}

func parseMsgHandler(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	msg := NewMsg(fmt.Sprintf("%s", body))

	msgjson, err := json.Marshal(msg)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Printf("Error: %s\n", err)
		return
	}

	w.Write(msgjson)
}
