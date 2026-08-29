package main

import (
	"fmt"
	"net/http"
	"os"
)

var version = "dev"

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, DevOps! version=%s\n", version)
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	fmt.Println("listening on :" + port)
	http.ListenAndServe(":"+port, nil)
}
