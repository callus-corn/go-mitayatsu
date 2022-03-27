package main

import (
	"archive/zip"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"text/template"
	"time"
)

type Contents struct {
	Title string `json:"title"`
	Worth string `json:"worth"`
}

type Review struct {
	Title string
	Text  string
}

var authToken string
var authTimeLimit time.Time

func main() {
	if len(os.Args) < 2 {
		fmt.Println("no args")
		return
	}
	authToken = os.Args[1]
	authTimeLimit = time.Now().Add(30 * time.Minute)

	http.HandleFunc("/", serveHTTP)
	http.HandleFunc("/static/", serveStatic)
	http.HandleFunc("/images/", serveImage)
	http.HandleFunc("/api/upload", upload)
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
	fileData, err := os.ReadFile("contentsdata/contents.json")
	if err != nil {
		return result
	}
	err = json.Unmarshal(fileData, &result)
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
	file, err := os.ReadFile("contentsdata/review/" + title)
	if err != nil {
		return Review{}, err
	}
	text := strings.Replace(string(file), "\n", "<br>", -1)
	return Review{Title: title, Text: text}, nil
}

func serveStatic(w http.ResponseWriter, req *http.Request) {
	http.StripPrefix("/static/", http.FileServer(http.Dir("static"))).ServeHTTP(w, req)
}

func serveImage(w http.ResponseWriter, req *http.Request) {
	http.StripPrefix("/images/", http.FileServer(http.Dir("contentsdata/images"))).ServeHTTP(w, req)
}

func upload(w http.ResponseWriter, req *http.Request) {
	if req.Method != "POST" {
		return
	}
	if err := req.ParseMultipartForm(1 << 30); err != nil {
		return
	}
	if authTimeLimit.Before(time.Now()) {
		return
	}
	if req.FormValue("token") != authToken {
		return
	}

	file, header, err := req.FormFile("contents")
	if err != nil {
		return
	}
	zipFile, err := os.Create(header.Filename)
	if err != nil {
		return
	}
	defer zipFile.Close()
	_, err = io.Copy(zipFile, file)
	if err != nil {
		return
	}

	reader, err := zip.OpenReader(header.Filename)
	if err != nil {
		return
	}
	defer reader.Close()

	if err = os.Mkdir("contentsdata", os.ModeDir); err != nil {
		return
	}
	for _, f := range reader.File {
		if f.Mode().IsDir() {
			continue
		}
		destPath := filepath.Join("contentsdata", f.Name)
		if err = os.MkdirAll(filepath.Dir(destPath), f.Mode()); err != nil {
			return
		}
		readCloser, err := f.Open()
		if err != nil {
			return
		}
		defer readCloser.Close()
		destFile, err := os.Create(destPath)
		if err != nil {
			return
		}
		defer destFile.Close()
		if _, err = io.Copy(destFile, readCloser); err != nil {
			return
		}
	}
}
