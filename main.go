package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"os"
	"strings"
	"text/template"
)

type Contents struct {
	Title string `json:"title"`
	Worth string `json:"worth"`
}

type Review struct {
	Title string
	Text  string
}

func main() {
	http.HandleFunc("/", serveHTTP)
	http.HandleFunc("/static/", serveStatic)
	http.HandleFunc("/images/", serveImage)
	http.HandleFunc("/api/upload", upload)
	http.HandleFunc("/api/auth", auth)
	http.ListenAndServe(":80", nil)
}

func serveHTTP(w http.ResponseWriter, req *http.Request) {
	path := req.URL.Path

	if path == "/" {
		serveIndex(w, req)
	} else {
		serveContents(w, req)
	}
}

func serveIndex(w http.ResponseWriter, req *http.Request) {
	contentsList := getContentsList()
	indexContents := map[string][]Contents{}
	for _, v := range contentsList {
		indexContents[v.Worth] = append(indexContents[v.Worth], v)
	}

	index, err := template.ParseFiles("templates/index.tmpl")
	if err != nil {
		http.Error(w, err.Error(), 404)
		return
	}
	if err = index.Execute(w, indexContents); err != nil {
		http.Error(w, err.Error(), 404)
	}
}

func serveContents(w http.ResponseWriter, req *http.Request) {
	title := strings.Replace(req.URL.Path, "/", "", 1)
	contents, err := getContents(title)
	if err != nil {
		http.Error(w, err.Error(), 404)
		return
	}
	review, err := getReview(title)
	if err != nil {
		http.Error(w, err.Error(), 404)
		return
	}

	contentsData := map[string]string{
		"Title":  title,
		"Worth":  contents.Worth,
		"Review": review.Text,
	}

	page, err := template.ParseFiles("templates/contents.tmpl")
	if err != nil {
		http.Error(w, err.Error(), 404)
		return
	}
	if err = page.Execute(w, contentsData); err != nil {
		http.Error(w, err.Error(), 404)
	}
}

func getContentsList() []Contents {
	result := []Contents{}
	file, err := os.ReadFile("contentsdata/contents.json")
	if err != nil {
		return result
	}
	err = json.Unmarshal(file, &result)
	if err != nil {
		return result
	}
	return result
}

func getContents(title string) (Contents, error) {
	contentsList := getContentsList()
	for _, v := range contentsList {
		if v.Title == title {
			return v, nil
		}
	}
	return Contents{}, errors.New("404 Not Found")
}

func getReview(title string) (Review, error) {
	text, err := os.ReadFile("contentsdata/review/" + title)
	if err != nil {
		return Review{}, err
	}
	return Review{Title: title, Text: string(text)}, nil
}

func serveStatic(w http.ResponseWriter, req *http.Request) {
	http.StripPrefix("/static/", http.FileServer(http.Dir("static"))).ServeHTTP(w, req)
}

func serveImage(w http.ResponseWriter, req *http.Request) {
	http.StripPrefix("/images/", http.FileServer(http.Dir("contentsdata/images"))).ServeHTTP(w, req)
}

func upload(w http.ResponseWriter, req *http.Request) {}

func auth(w http.ResponseWriter, req *http.Request) {}
