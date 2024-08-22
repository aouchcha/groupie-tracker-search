package functions

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	"strings"
)

type Err struct {
	Message string
	Title   string
	Code    int
}

var Error Err

func ServeStyle(w http.ResponseWriter, r *http.Request) {
	v := http.StripPrefix("/styles/", http.FileServer(http.Dir("./styles")))
	tmpl1, err2 := template.ParseFiles("templates/errors.html")
	if err2 != nil {
		http.Error(w, "Error 500", http.StatusInternalServerError)
		return
	}
	if r.URL.Path == "/styles/" {
		ChooseError(w, 403)
		tmpl1.Execute(w, Error)
		return
	}
	v.ServeHTTP(w, r)
}

func FirstPage(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/welcome.html")
	tmpl1, err2 := template.ParseFiles("templates/errors.html")

	if err != nil || err2 != nil {
		if err2 != nil {
			http.Error(w, "Error 500", http.StatusInternalServerError)
			return
		} else {
			ChooseError(w, 500)
			tmpl1.Execute(w, Error)
			return
		}
	}
	fmt.Println(r.URL.Path)
	if r.URL.Path != "/" {
		ChooseError(w, 404)
		tmpl1.Execute(w, Error)
		return
	}
	if r.Method != http.MethodGet {
		ChooseError(w, 405)
		tmpl1.Execute(w, Error)
		return
	}
	tmpl.Execute(w, artists)
}

// var IND int

func OtherPages(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/details.html")
	tmpl1, err2 := template.ParseFiles("templates/errors.html")

	if err != nil || err2 != nil {
		if err2 != nil {
			http.Error(w, "Error 500", http.StatusInternalServerError)
			return
		} else {
			ChooseError(w, 500)
			tmpl1.Execute(w, Error)
			return
		}
	}

	fmt.Println(r.URL.Path)
	if r.URL.Path != "/artist" {
		// fmt.Println("vvvv")
		ChooseError(w, 404)
		tmpl1.Execute(w, Error)
		return
	}
	max := artists[len(artists)-1].ID
	fmt.Println(max)
	url := r.URL.Query().Get("ID")
	// fmt.Println(url)
	index, err := strconv.Atoi(string(url))
	if err != nil || index < 1 || index > max {
		// fmt.Println("hna")
		ChooseError(w, 404)
		tmpl1.Execute(w, Error)
		return
	}
	index -= 1
	// IND = index
	if r.Method != http.MethodGet {
		ChooseError(w, 405)
		tmpl1.Execute(w, Error)
		return
	}

	artistinfo := struct {
		ID            int
		Name          string
		Image         string
		Members       []string
		CreationDate  int
		FirstAlbum    string
		Localisations []string
		Relations     map[string][]string
	}{
		ID:            artists[index].ID,
		Name:          artists[index].Name,
		Image:         artists[index].Image,
		Members:       artists[index].Members,
		CreationDate:  artists[index].CreationDate,
		FirstAlbum:    artists[index].FirstAlbum,
		Localisations: locals.Index[index].Location,
		Relations:     rel.Index[index].DateLocations,
	}
	tmpl.Execute(w, artistinfo)
}

func SearchPage(w http.ResponseWriter, r *http.Request) {
	tmpl, _ := template.ParseFiles("./templates/search.html")
	tmpl1, _ := template.ParseFiles("templates/errors.html")

	text := r.FormValue("search")
	var ss []Artist
	for i := range artists {
		if strings.Contains(artists[i].Name, text) {
			ss = append(ss, artists[i])
		}
	}
	if len(ss) == 0 {
		ChooseError(w, 1000)
		tmpl1.Execute(w, Error)
		return
	}
	tmpl.Execute(w, ss)
}

func ChooseError(w http.ResponseWriter, code int) {
	if code == 404 {
		Error.Title = "Error 404"
		Error.Message = "The page web doesn't exist\nError 404"
		Error.Code = code
		w.WriteHeader(code)
	} else if code == 405 {
		Error.Title = "Error 405"
		Error.Message = "The method is not alloweded\nError 405"
		Error.Code = code
		w.WriteHeader(code)
	} else if code == 400 {
		Error.Title = "Error 400"
		Error.Message = "Bad Request\nError 400"
		Error.Code = code
		w.WriteHeader(code)
	} else if code == 500 {
		Error.Title = "Error 500"
		Error.Message = "Internal Server Error\nError 500"
		Error.Code = code
		w.WriteHeader(code)
	} else if code == 403 {
		Error.Title = "Error 403"
		Error.Message = "This page web is forbidden\nError 403"
		Error.Code = code
		w.WriteHeader(code)
	}
}
