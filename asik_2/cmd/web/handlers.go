package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"path/filepath"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

var templates *template.Template

type Article struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

var db *sql.DB

func init() {
	templatesPath := "ui/html"
	templatesPattern := filepath.Join(templatesPath, "*.html")

	templateFiles, err := filepath.Glob(templatesPattern)
	if err != nil {
		panic(err)
	}

	if len(templateFiles) == 0 {
		panic("No template files found")
	}

	templateNames := make([]string, len(templateFiles))
	for i, file := range templateFiles {
		_, fileName := filepath.Split(file)
		templateNames[i] = strings.TrimSuffix(fileName, filepath.Ext(fileName))
	}

	templates = template.Must(template.New("").Funcs(template.FuncMap{"safeHTML": func(b string) template.HTML { return template.HTML(b) }}).ParseFiles(templateFiles...))
}

func Home(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query("SELECT title, content FROM articles")
	if err != nil {
		http.Error(w, "Failed to retrieve articles from MySQL", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var retrievedArticles []Article
	for rows.Next() {
		var article Article
		err := rows.Scan(&article.Title, &article.Content)
		if err != nil {
			http.Error(w, "Failed to decode articles from MySQL", http.StatusInternalServerError)
			return
		}
		retrievedArticles = append(retrievedArticles, article)
	}

	renderTemplate(w, "index", retrievedArticles)
}

func renderTemplate(w http.ResponseWriter, tmpl string, data interface{}) {
	tmpl = filepath.Join("ui/html", tmpl) // указывает путь к вашим шаблонам
	t, err := template.ParseFiles(tmpl + ".html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = t.Execute(w, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func Blog(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "blog", nil)
}

func Contact(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "contact", nil)
}

func FullWidth(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "fullwidth", nil)
}

func addArticleHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse JSON data from the request body
	var newArticle Article
	err := json.NewDecoder(r.Body).Decode(&newArticle)
	if err != nil {
		http.Error(w, "Failed to parse request body", http.StatusBadRequest)
		return
	}

	// Validate the data
	if err := validateArticle(newArticle); err != nil {
		http.Error(w, fmt.Sprintf("Validation failed: %s", err), http.StatusBadRequest)
		return
	}

	// Insert the article into MySQL
	_, err = db.Exec("INSERT INTO articles (title, content) VALUES (?, ?)", newArticle.Title, newArticle.Content)
	if err != nil {
		http.Error(w, "Failed to insert article into MySQL", http.StatusInternalServerError)
		return
	}

	response := map[string]string{"message": "Article created successfully"}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func validateArticle(article Article) error {
	if article.Title == "" {
		return fmt.Errorf("Title is required")
	}
	if article.Content == "" {
		return fmt.Errorf("Content is required")
	}
	return nil
}
