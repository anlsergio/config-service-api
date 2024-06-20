package main

import (
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Hello World"))
	})

	log.Println("Starting server on port 8080")
	// TODO: listen for syscalls to shutdown server gracefully
	log.Fatal(http.ListenAndServe(":8080", router))
}
