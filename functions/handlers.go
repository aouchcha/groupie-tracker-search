package functions

import (
	"encoding/json"
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

func SuggestHandler(w http.ResponseWriter, r *http.Request) {
	input := r.URL.Query().Get("search")

	suggestions := getSuggestions(input)

	w.Header().Set("Content-Type", "text/plain")
	json.NewEncoder(w).Encode(suggestions)
}

func getSuggestions(input string) []string {
	var suggestions []string
	input = strings.ToLower(input)
	for i := range artists {
		if strings.HasPrefix(strings.ToLower(artists[i].Name), input) {
			suggestions = append(suggestions, artists[i].Name+"-> Band")
		}
		if strings.HasPrefix(strings.ToLower(artists[i].FirstAlbum), input) {
			suggestions = append(suggestions, artists[i].FirstAlbum+"-> First Album Date")
		}
		if strings.HasPrefix(strings.ToLower(strconv.Itoa(artists[i].CreationDate)), input) {
			suggestions = append(suggestions, strconv.Itoa(artists[i].CreationDate)+"-> First Album Date")
		}
		for j := range artists[i].Members {
			if strings.HasPrefix(strings.ToLower(artists[i].Members[j]), input) {
				suggestions = append(suggestions, artists[i].Members[j]+"->Member")
			}
		}
	}
	for i := range locals.Index {
		for j := range locals.Index[i].Location {
			if strings.Contains(strings.ToLower(locals.Index[i].Location[j]), input) {
				suggestions = append(suggestions, locals.Index[i].Location[j]+"->Location")
			}
		}
	}
	return suggestions
}

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

	if r.URL.Path != "/artist" {
		ChooseError(w, 404)
		tmpl1.Execute(w, Error)
		return
	}
	max := artists[len(artists)-1].ID
	url := r.URL.Query().Get("ID")
	index, err := strconv.Atoi(string(url))
	if err != nil || index < 1 || index > max {
		ChooseError(w, 404)
		tmpl1.Execute(w, Error)
		return
	}
	index -= 1
	if r.Method != http.MethodGet {
		ChooseError(w, 405)
		tmpl1.Execute(w, Error)
		return
	}

	artistinfos := struct {
		ID            int
		Name          string
		Image         string
		Members       []string
		CreationDate  int
		FirstAlbum    string
		Localisations []string
		Relations     map[string][]string
		Dates         []string
	}{
		ID:            artists[index].ID,
		Name:          artists[index].Name,
		Image:         artists[index].Image,
		Members:       artists[index].Members,
		CreationDate:  artists[index].CreationDate,
		FirstAlbum:    artists[index].FirstAlbum,
		Localisations: locals.Index[index].Location,
		Relations:     rel.Index[index].DateLocations,
		Dates:         dat.Index[index].Date,
	}
	tmpl.Execute(w, artistinfos)
}

func SearchPage(w http.ResponseWriter, r *http.Request) {
	tmpl, _ := template.ParseFiles("./templates/search.html")
	tmpl1, _ := template.ParseFiles("templates/errors.html")
	types := r.FormValue("typessearch")
	// fmt.Println(types)
	text := strings.ToLower(strings.TrimSuffix(r.FormValue("search"), " "))
	temp := strings.Split(text, "->")
	text = temp[0]
	if text == "" {
		ChooseError(w, 400)
		tmpl1.Execute(w, Error)
		return
	}

	var ss []Artist
	if types == "Band" || types == "fistalbum" || types == "creation" {
		for i := range artists {
			if strings.HasPrefix(strings.ToLower(artists[i].Name), text) || strings.HasPrefix(strings.ToLower(artists[i].FirstAlbum), text) || strings.HasPrefix(strings.ToLower(strconv.Itoa(artists[i].CreationDate)), text) {
				ss = append(ss, artists[i])
			}
		}
	} else if types == "Members" {
		for i := range artists {
			for j := range artists[i].Members {
				if strings.HasPrefix(strings.ToLower(artists[i].Members[j]), text) {
					ss = append(ss, artists[i])
				}
			}
		}
	} else if types == "location" {
		for i := range locals.Index {
			for j := range locals.Index[i].Location {
				if strings.Contains(strings.ToLower(locals.Index[i].Location[j]), text) {
					if len(ss) == 0 {
						ss = append(ss, artists[locals.Index[i].Id-1])
					} else {
						var check bool
						for k := range ss {
							if ss[k].ID == locals.Index[i].Id {
								check = true
							} else {
								check = false
							}
						}
						if !check {
							ss = append(ss, artists[locals.Index[i].Id-1])
						}
					}

				}
			}
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
