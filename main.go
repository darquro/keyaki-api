package main

import (
	"encoding/json"
	"net/http"
	"strconv"

	"google.golang.org/appengine"
	"google.golang.org/appengine/urlfetch"
)

// func main() {
// 	http.HandleFunc("/blog", blogHandler)
// 	// http.HandleFunc("/news", newsHandler)
// 	http.ListenAndServe(":8080", nil)
// }

func init() {
	http.HandleFunc("/blog", blogHandler)
	http.HandleFunc("/news", newsHandler)
}

func blogHandler(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	member, err := strconv.Atoi(q.Get("m"))
	if err != nil {
		member = -1
	}
	page, err := strconv.Atoi(q.Get("p"))
	if err != nil {
		page = -1
	}

	blogURL, err := getBlogURL(member, page)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	client := urlfetch.Client(appengine.NewContext(r))
	resp, err := client.Get(blogURL.String())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	results, err := parseBlogResponse(resp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	b, err := json.Marshal(results)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Write(b)
}

func newsHandler(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	year, err := strconv.Atoi(q.Get("y"))
	if err != nil {
		year = 0
	}
	month, err := strconv.Atoi(q.Get("m"))
	if err != nil {
		month = 0
	}

	newsURL, err := getNewsURL(year, month)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	client := urlfetch.Client(appengine.NewContext(r))
	resp, err := client.Get(newsURL.String())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	results, err := parseNewsResponse(resp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	b, err := json.Marshal(results)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Write(b)
}
