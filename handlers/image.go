package handlers

import (
	"fmt"
	"html/template"
	"net/http"
)

func ImageUploadHandler(w http.ResponseWriter, r *http.Request) {
	file, err := checkFileTypeAndSize(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer file.Close()
	fmt.Println("File is valid. Processing the file...")
}

func ImageHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("templates/image.html"))
	tmpl.Execute(w, nil)
}