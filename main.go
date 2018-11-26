package main

import (
	"fmt"
	"html/template"
	"net/http"

	"awesomeProject/models"
)

var posts map[string]*models.Post

func IndexHendler(writer http.ResponseWriter, request *http.Request) {
	t, err := template.ParseFiles("templates/index.html", "templates/header.html", "templates/footer.html")
	if err != nil {
		fmt.Fprintf(writer, err.Error())
	}

	fmt.Println(posts)

	t.ExecuteTemplate(writer, "index", posts)

}

func WriteHendler(writer http.ResponseWriter, request *http.Request) {
	t, err := template.ParseFiles("templates/write.html", "templates/header.html", "templates/footer.html")
	if err != nil {
		fmt.Fprintf(writer, err.Error())
	}

	t.ExecuteTemplate(writer, "write", nil)

}

func EditHendler(writer http.ResponseWriter, request *http.Request) {
	t, err := template.ParseFiles("templates/write.html", "templates/header.html", "templates/footer.html")
	if err != nil {
		fmt.Fprintf(writer, err.Error())
	}

	id := request.FormValue("id")
	post, found := posts[id]
	if !found {
		http.NotFound(writer, request)
	}

	t.ExecuteTemplate(writer, "write", post)

}

func savePostHendler(writer http.ResponseWriter, request *http.Request) {
	id := request.FormValue("id")
	title := request.FormValue("title")
	content := request.FormValue("content")

	var post *models.Post
	if id != "" {
		post = posts[id]
		post.Title = title
		post.Content = content
	} else {
		id := GenerateId()
		post := models.NewPost(id, title, content)
		posts[post.Id] = post

	}

	http.Redirect(writer, request, "/", 302)

}

func DeleteHendler(writer http.ResponseWriter, request *http.Request) {
	id := request.FormValue("id")
	if id == "" {
		http.NotFound(writer, request)
	}

	delete(posts, id)

	http.Redirect(writer, request, "/", 302)

}

func main() {
	fmt.Println("qwertydfsf sdfdsffsuj")

	posts = make(map[string]*models.Post, 0)

	http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("./assets/"))))
	http.HandleFunc("/", IndexHendler)
	http.HandleFunc("/write", WriteHendler)
	http.HandleFunc("/edit", EditHendler)
	http.HandleFunc("/delete", DeleteHendler)
	http.HandleFunc("/SavePost", savePostHendler)

	http.ListenAndServe(":3000", nil)
}
