package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB
var albums []Album

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
	tmpl, err := template.ParseFiles("result.html")
	if err != nil {
		fmt.Fprintln(w, err.Error())
	}
	err = tmpl.Execute(w, albums)
	if err != nil {
		fmt.Fprintln(w, err.Error())
	}
}

func queryHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("query.html"))
	tmpl.Execute(w, nil)
}

func main() {
	var err error

	db, err = sql.Open("sqlite3", "assets/recordings.sqlite")
	if err != nil {
		log.Fatal("Couldn't open database.")
	}
	defer db.Close()

	name := make([]rune, 0, 20)
	for i, s := range os.Args {
		switch {
		case i == 0:
			continue
		case i > 1:
			name = append(name, ' ')
		}
		name = append(name, []rune(s)...)
	}
	namestring := string(name)

	fmt.Printf("\nLooking for albums by '%s'\n\n", namestring)

	albums, err = albumsByArtist(namestring)
	if err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/result", resultHandler)
	http.HandleFunc("/", queryHandler)
	http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("assets"))))
	log.Fatal(http.ListenAndServe(":8080", nil))
}
