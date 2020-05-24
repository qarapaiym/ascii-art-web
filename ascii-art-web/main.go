package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"

	student ".."
)

type Post struct {
	Id      string
	Title   string
	Content string
}

func NewPost(id, title, content string) *Post {
	return &Post{id, title, content}
}

var posts map[string]*Post

func indexHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		errorHandler(w, r, http.StatusNotFound)
		return
	}
	t, err := template.ParseFiles("templates/index.html", "templates/header.html", "templates/footer.html", "templates/write.html")
	if err != nil {
		fmt.Fprintf(w, err.Error())
		return
	}
	fmt.Println(posts)

	t.ExecuteTemplate(w, "index", posts)
}
func writeHandler(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("write.html")
	if err != nil {
		fmt.Fprintf(w, err.Error())
		return
	}

	t.ExecuteTemplate(w, "index", nil)
}

func savePostHandler(w http.ResponseWriter, r *http.Request) {
	id := r.FormValue("id")
	title := r.FormValue("title")
	content := r.FormValue("content")
	change := ""
	if ValidAscii(content) {
		fmt.Println("Content: ", content)
		change = student.AsciiWeb(string(content), title)
	} else {
		w.WriteHeader(http.StatusBadRequest)
		t, err := template.ParseFiles("templates/400.html")
		if err != nil {
			internalServerError(w, r)
		}
		t.Execute(w, nil)
	}

	post := NewPost(id, change, content)
	posts[post.Id] = post
	http.Redirect(w, r, "/", 302)
}

func internalServerError(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusInternalServerError)
	t, _ := template.ParseFiles("templates/500.html")
	err := t.Execute(w, nil)
	if err != nil {
		log.Fatal(err)
	}
}
func notFound(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	t, _ := template.ParseFiles("templates/lol.html")
	err := t.Execute(w, nil)
	if err != nil {
		log.Fatal(err)
	}
}
func ValidAscii(s string) bool {
	for _, i := range []byte(s) {
		if i > 127 {
			return false
		}
	}
	return true
}
func badRequest(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusBadRequest)
	t, _ := template.ParseFiles("templates/400.html")
	err := t.Execute(w, nil)
	if err != nil {
		log.Fatal(err)
	}
}
func errorHandler(w http.ResponseWriter, r *http.Request, status int) {
	w.WriteHeader(status)
	if status == http.StatusNotFound {
		notFound(w, r)
	} else if status == 500 {
		internalServerError(w, r)
	} else if status == 400 {
		badRequest(w, r)
	}
}
func main() {
	posts = make(map[string]*Post, 0)
	fmt.Println("Listening on port :3000")
	http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("./assets/"))))
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("write", writeHandler)
	http.HandleFunc("/SavePost", savePostHandler)
	http.ListenAndServe(":3000", nil)
}
