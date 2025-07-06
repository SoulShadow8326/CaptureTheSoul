package main

import (
	"fmt"
	"net/http"
	"os"
)

func main() {
	http.HandleFunc("/read", func(w http.ResponseWriter, r *http.Request) {
		file := r.URL.Query().Get("file")
		if file == "" {
			http.Error(w, "file param required", http.StatusBadRequest)
			return
		}

		data, err := os.ReadFile(file)
		if err != nil {
			http.Error(w, "could not read file", http.StatusInternalServerError)
			return
		}

		w.Write(data)
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Welcome to your vulnerable service. Try /read?file=/flag.txt")
	})

	http.ListenAndServe(":1337", nil)
}
