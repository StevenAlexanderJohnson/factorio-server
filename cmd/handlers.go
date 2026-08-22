package main

import (
	"net/http"
)

func handleHealthz(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("healthy"))
}
