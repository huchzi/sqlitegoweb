package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"net/url"
	"strconv"

	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB
var albums []Album
var templates = template.Must(template.ParseFiles("entryForm.html", "query.html", "result.html"))

type Album struct {
	Id     int64
	Title  string
	Artist string
	Price  float32
}

// albumsByArtist querys the db database for records by a specific artist
func albumsByArtist(name string) ([]Album, error) {
	var albums []Album

	rows, err := db.Query("SELECT * FROM album WHERE artist = ?", name)
	if err != nil {
		return nil, fmt.Errorf("albumsByArtist: %q: %v", name, err)
	}
	defer rows.Close()

	for rows.Next() {
		var alb Album
		if err := rows.Scan(&alb.Id, &alb.Title, &alb.Artist, &alb.Price); err != nil {
			return nil, fmt.Errorf("albumsByArtist: %q: %v", name, err)
		}
		albums = append(albums, alb)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("albumsByArtist: %q: %v", name, err)
	}

	return albums, nil
}

func resultHandler(w http.ResponseWriter, r *http.Request) {
	var err error
	var artistName string

	fmt.Println(r.FormValue("artist"))
	if r.FormValue("artist") != "" {
		artistName = r.FormValue("artist")
		writeToDB(r)
	} else {
		artistName = r.FormValue("name")
	}

	albums, err = albumsByArtist(artistName)
	if err != nil {
		log.Fatal(err)
	}
	templates.ExecuteTemplate(w, "result.html", albums)
}

func queryHandler(w http.ResponseWriter, r *http.Request) {
	templates.ExecuteTemplate(w, "query.html", nil)
}

func newEntryHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	if r.FormValue("artist") == "" {
		templates.ExecuteTemplate(w, "entryForm.html", nil)
		fmt.Println("First entry")
	} else {
		artistValid := r.FormValue("artist") != ""
		titleValid := r.FormValue("title") != ""
		_, err := strconv.ParseFloat(r.FormValue("price"), 32)
		priceValid := err == nil

		if artistValid && titleValid && priceValid {
			fmt.Println("Entry accepted")
			r.ParseForm()
			fmt.Println(r.Form)
			v := url.Values{}
			v.Set("artist", "John Coltrane")
			r.Form = v
			http.Redirect(w, r, "/result", http.StatusFound)
		} else {
			fmt.Println("Entry error")
			templates.ExecuteTemplate(w, "query.html", nil)
			templates.ExecuteTemplate(w, "entryForm.html", "Entry error")
		}
	}
}

func writeToDB(r *http.Request) {
	price, _ := strconv.ParseFloat(r.FormValue("price"), 32)

	newAlbum := Album{Id: 1, Artist: r.FormValue("artist"), Title: r.FormValue("title"), Price: float32(price)}
	_, err := db.Exec("INSERT INTO album (title, artist, price) VALUES (?, ?, ?)", newAlbum.Title, newAlbum.Artist, newAlbum.Price)
	if err != nil {
		fmt.Println(err.Error())
	}
}

func main() {
	var err error

	db, err = sql.Open("sqlite3", "assets/recordings.sqlite")
	if err != nil {
		log.Fatal("Couldn't open database.")
	}
	defer db.Close()

	http.HandleFunc("/", queryHandler)
	http.HandleFunc("/result", resultHandler)
	http.HandleFunc("/entryForm", newEntryHandler)
	http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("assets"))))
	log.Fatal(http.ListenAndServe(":8080", nil))
}
