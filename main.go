package main

import (
	"net/http"

	f "groupie-tracker/functions"
)

// Define a struct that matches the JSON structure

func main() {
	f.GetArtistData()
	f.GetRelationData()
	f.GetLocationData()
	http.HandleFunc("/styles/", f.ServeStyle)
	http.HandleFunc("/", f.FirstPage)
	http.HandleFunc("/artist", f.OtherPages)
	http.HandleFunc("/search", f.SearchPage)
	http.ListenAndServe(":9991", nil)
}
