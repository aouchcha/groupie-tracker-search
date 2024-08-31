package main

import (
	"fmt"
	"net/http"

	f "groupie-tracker-search/functions"
)

func main() {
	f.GetArtistData()
	f.GetRelationData()
	f.GetLocationData()
	f.GetDatesData()
	http.HandleFunc("/styles/", f.ServeStyle)
	http.HandleFunc("/", f.FirstPage)
	http.HandleFunc("/suggest", f.SuggestHandler)

	http.HandleFunc("/artist", f.OtherPages)
	http.HandleFunc("/search", f.SearchPage)
	fmt.Println("http://localhost:6699/")
	http.ListenAndServe(":6699", nil)
}
