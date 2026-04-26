package main

import (
	"fmt"
	"net/http"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Selamat Datang di API Pantau Gizi! Server berjalan dengan baik.")
	})

	fmt.Println("Server berjalan di port 8080...")
	http.ListenAndServe(":8080", nil)
}
