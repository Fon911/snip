package main

import (
	"fmt"
	"net/http"
	"strconv"
)

func Hello(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	if name == "" {
		fmt.Fprintln(w, "Hello input name  /Hello?name=Lexa", name)
		return
	}
	fmt.Fprintln(w, "Hello world ", name)
}

var links = map[string]string{}
var KeyCount int

func Shorten(w http.ResponseWriter, r *http.Request) {

	val := r.URL.Query().Get("url")
	if val == "" {
		fmt.Fprintln(w, "нужно так: /shorten?url=ссылка")
		return
	}
	KeyCount++
	key := strconv.Itoa(KeyCount)
	links[key] = val
	fmt.Fprintln(w, "записал:", key)
}

func get(w http.ResponseWriter, r *http.Request) {
	key := r.URL.Query().Get("code")
	val, ok := links[key]
	if !ok {
		fmt.Fprintln(w, "нет такой записи")
		return
	}
	http.Redirect(w, r, val, http.StatusMovedPermanently)
}

func main() {
	http.HandleFunc("/Hello", Hello)
	http.HandleFunc("/shorten", Shorten)
	http.HandleFunc("/go", get)

	err := http.ListenAndServe(":8000", nil)
	if err != nil {
		fmt.Println(err)
	}
}
