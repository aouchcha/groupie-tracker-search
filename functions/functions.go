package functions

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type Artist struct {
	ID           int      `json:"id"`
	Image        string   `json:"image"`
	Name         string   `json:"name"`
	Members      []string `json:"members"`
	CreationDate int      `json:"creationDate"`
	FirstAlbum   string   `json:"firstAlbum"`
}

type Dates struct {
	Index []SubDates `json:"index"`
}

type SubDates struct {
	Id            int                 `json:"id"`
	DateLocations map[string][]string `json:"datesLocations"`
}

type Locations struct {
	Index []SubLocal `json:"index"`
 }

 type SubLocal struct {
	Id int `json:"id"`
	Location []string `json:"locations"`
 }

var (
	artists []Artist
	rel     Dates
	locals Locations
)

func GetArtistData() {
	res, err := http.Get("https://groupietrackers.herokuapp.com/api/artists")
	if err != nil {
		fmt.Println("Error in getting the data from the artist link:", err)
		return
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println("Error in reading the body:", err)
		return
	}

	err = json.Unmarshal(body, &artists)
	if err != nil {
		fmt.Println("Error with JSON unmarshal:", err)
		return
	}
}

func GetRelationData() {
	res, err := http.Get("https://groupietrackers.herokuapp.com/api/relation")
	if err != nil {
		fmt.Println("Error in getting the data from the artist link:", err)
		return
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println("Error in reading the body:", err)
		return
	}

	err = json.Unmarshal(body, &rel)
	if err != nil {
		fmt.Println("Error with JSON unmarshal:", err)
		return
	}
}

func GetLocationData() {
	res, err := http.Get("https://groupietrackers.herokuapp.com/api/locations")
	if err != nil {
		fmt.Println("Error in getting the data from the artist link:", err)
		return
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println("Error in reading the body:", err)
		return
	}

	err = json.Unmarshal(body, &locals)
	if err != nil {
		fmt.Println("Error with JSON unmarshal:", err)
		return
	}
}
