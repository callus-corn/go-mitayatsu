package main

import (
	"net/http"
	"strings"
	"text/template"
)

type Contents struct {
	Title     string
	Star      string
	Recommend string
	Comment   string
	Time      string
	Media     string
}

func main() {
	handleAdmin()
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	http.HandleFunc("/", contents)
	http.HandleFunc("/api/upload", upload)
	http.ListenAndServe(":80", nil)
}

func contents(w http.ResponseWriter, req *http.Request) {
	path := req.URL.Path

	if path == "/" {
		render_index(w, req)
	} else {
		render_contents(w, req)
	}
}

func render_index(w http.ResponseWriter, req *http.Request) {
	contentsTypes := []string{
		"時間作ってでもみた方がいい",
		"時間があるなら見た方がいい",
		"暇なら見た方がいい",
		"誰にでもオススメ",
		"大体の人にオススメ",
		"コアな人にオススメ",
	}

	contentsList := map[string][]Contents{}
	for _, v := range contentsTypes {
		contentsList[v] = getContentsList(v)
	}

	index, err := template.ParseFiles("templates/index.tmpl")
	if err != nil {
		http.Error(w, err.Error(), 404)
		return
	}
	if err = index.Execute(w, contentsList); err != nil {
		http.Error(w, err.Error(), 404)
	}
}

func render_contents(w http.ResponseWriter, req *http.Request) {
	title := strings.Replace(req.URL.Path, "/", "", 1)
	contents := getContents(title)
	page, err := template.ParseFiles("templates/contents.tmpl")
	if err != nil {
		http.Error(w, err.Error(), 404)
		return
	}
	if err = page.Execute(w, contents); err != nil {
		http.Error(w, err.Error(), 404)
	}
}

func getContentsList(contentsType string) []Contents {
	c1 := Contents{
		Title:     "プリパラ",
		Star:      "時間作ってでもみた方がいい",
		Recommend: "大体の人にオススメ",
		Comment:   "みろ",
		Time:      "30分×38話",
		Media:     "dアニメストア",
	}
	return []Contents{c1, c1, c1, c1, c1, c1, c1}
}

func getContents(title string) Contents {
	return Contents{
		Title:     "プリパラ",
		Star:      "時間作ってでもみた方がいい",
		Recommend: "大体の人にオススメ",
		Comment:   "みろ",
		Time:      "30分×38話",
		Media:     "dアニメストア",
	}
}
