package main

import (
	"fmt"
	"net/http"
)

func main() {
	fmt.Println("Server starting; listening on port 3000")
	http.ListenAndServe(":3000", Handlers())
}
