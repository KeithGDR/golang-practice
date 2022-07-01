package server

import (
	"fmt"
	"log"
	"net/http"
)

func onRequest(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.Error(w, "404 not found.", http.StatusNotFound)
		return
	}

	
}

func main() {
	http.HandleFunc("/", onRequest)

	fmt.Printf("Starting web server...\n")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
