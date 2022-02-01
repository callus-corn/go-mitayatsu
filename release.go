//go:build !debug

package main

import (
	"net/http"
)

func upload(w http.ResponseWriter, req *http.Request) {
	http.NotFound(w, req)
}

func handleAdmin() {}
