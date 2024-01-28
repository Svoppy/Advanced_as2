package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	var err error
	db, err = sql.Open("mysql", "root:708204@tcp(127.0.0.1:3306)/asik2")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to MySQL.")
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("ui/static"))))

	http.HandleFunc("/", Home)

	http.HandleFunc("/contact", Contact)
	http.HandleFunc("/blog", Blog)
	http.HandleFunc("/fullwidth", FullWidth)

	http.HandleFunc("/form", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "ui/html/form.html")
	})

	http.HandleFunc("/addArticle", addArticleHandler)

	http.HandleFunc("/success", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, `<p>Article created successfully!</p>
                  <a href="/index" id="homenav">Go to Home</a>`)
	})
	http.ListenAndServe(":8080", nil)
}
