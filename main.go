package main

import (
	"bytes"
	"net/http"
	"text/template"
)

type Contents struct {
	Title   string
	Image   string
	Star    string
	Comment string
}

func main() {
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	http.HandleFunc("/", render)
	http.ListenAndServe(":80", nil)
}

func render(w http.ResponseWriter, req *http.Request) {
	path := req.URL.Path
	if path != "/" {
		http.Error(w, "404 not found", 404)
		return
	}

	index, err := template.ParseFiles("templates/index.tmpl")
	if err != nil {
		http.Error(w, err.Error(), 404)
		return
	}

	data := map[string]string{}
	contentsTypes := []string{
		"GodContents",
		"GoodContents",
		"NormalContents",
		"PopularContents",
		"AlmostAllContents",
		"CoreContents",
	}

	for _, v := range contentsTypes {
		buf := bytes.NewBufferString("")
		tmpl, err := template.ParseFiles("templates/contentsList.tmpl")
		if err != nil {
			http.Error(w, err.Error(), 404)
			return
		}
		contentsList := getContentsList(v)
		if err := tmpl.Execute(buf, contentsList); err != nil {
			http.Error(w, err.Error(), 404)
			return
		}
		data[v] = buf.String()
	}

	if err := index.Execute(w, data); err != nil {
		http.Error(w, err.Error(), 404)
	}
}

func getContentsList(contentsType string) []Contents {
	c1 := Contents{
		Title:   "テスト",
		Image:   "",
		Star:    "",
		Comment: "",
	}
	return []Contents{c1}
}
