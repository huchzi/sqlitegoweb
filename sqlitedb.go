package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"
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

	albums, err = albumsByArtist(r.FormValue("name"))
	if err != nil {
		log.Fatal(err)
	}

	templates.ExecuteTemplate(w, "result.html", albums)
}

func queryHandler(w http.ResponseWriter, r *http.Request) {
	templates.ExecuteTemplate(w, "query.html", nil)
}

func newEntryHandler(w http.ResponseWriter, r *http.Request) {
	templates.ExecuteTemplate(w, "entryForm.html", nil)
}

func writeToDB(w http.ResponseWriter, r *http.Request) {
	price, _ := strconv.ParseFloat(r.FormValue("price"), 32)

	newAlbum := Album{Id: 1, Artist: r.FormValue("artist"), Title: r.FormValue("title"), Price: float32(price)}
	_, err := db.Exec("INSERT INTO album (title, artist, price) VALUES (?, ?, ?)", newAlbum.Title, newAlbum.Artist, newAlbum.Price)
	if err != nil {
		fmt.Println(err.Error())
	}
	http.Redirect(w, r, "/", http.StatusFound)
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
	http.HandleFunc("/writeToDB", writeToDB)
	http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("assets"))))
	log.Fatal(http.ListenAndServe(":8080", nil))
}
