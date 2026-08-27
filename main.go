package main

import (
	"fmt"
	"net/http"
)

func Hello(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	if name == "" {
		fmt.Fprintln(w, "Hello input name  /Hello?name=Lexa", name)
		return
	}
	fmt.Fprintln(w, "Hello world ", name)
}

func main() {
	http.HandleFunc("/Hello", Hello)

	err := http.ListenAndServe(":8000", nil)
	if err != nil {
		fmt.Println(err)
	}
}
