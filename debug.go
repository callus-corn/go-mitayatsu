//go:build debug

package main

import "net/http"

func upload(w http.ResponseWriter, req *http.Request) {

}

func handleAdmin() {
	http.Handle("/admin/", http.StripPrefix("/admin/", http.FileServer(http.Dir("admin"))))
}
