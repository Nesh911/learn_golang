package main

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/go-martini/martini"
	"github.com/martini-contrib/render"
	"github.com/russross/blackfriday"
	"learn_golang/models"
)

var posts map[string]*models.Post

func IndexHendler(rnd render.Render) {
	rnd.HTML(200, "index", posts)

}

func WriteHendler(rnd render.Render) {
	rnd.HTML(200, "write", nil)
}

func EditHendler(rnd render.Render, request *http.Request, params martini.Params) {
	id := params["id"]
	post, found := posts[id]
	if !found {
		rnd.Redirect("/")
	}

	rnd.HTML(200, "write", post)
}

func savePostHendler(rnd render.Render, request *http.Request) {
	id := request.FormValue("id")
	title := request.FormValue("title")
	contentMarkdown := request.FormValue("content")
	contentHtml := string(blackfriday.MarkdownBasic([]byte(contentMarkdown)))

	var post *models.Post
	if id != "" {
		post = posts[id]
		post.Title = title
		post.ContentHtml = contentHtml
		post.ContentMarkdown = contentMarkdown
	} else {
		id := GenerateId()
		post := models.NewPost(id, title, contentHtml, contentMarkdown)
		posts[post.Id] = post

	}

	rnd.Redirect("/")

}

func DeleteHendler(rnd render.Render, request *http.Request, params martini.Params) {
	id := params["id"]
	if id == "" {
		rnd.Redirect("/")
	}

	delete(posts, id)

	rnd.Redirect("/")

}

func getHtmlHendler(rnd render.Render, request *http.Request) {
	md := request.FormValue("md")
	htmlBytes := blackfriday.MarkdownBasic([]byte(md))

	rnd.JSON(200, map[string]interface{}{"html": string(htmlBytes)})
}

func unescape(x string) interface{} {
	return template.HTML(x)
}

func main() {
	fmt.Println("Hello")

	m := martini.Classic()

	unescapeFuncMap := template.FuncMap{"unescape": unescape}

	m.Use(render.Renderer(render.Options{
		Directory:  "templates",                         // Specify what path to load the templates from.
		Layout:     "layout",                            // Specify a layout template. Layouts can call {{ yield }} to render the current template.
		Extensions: []string{".tmpl", ".html"},          // Specify extensions to load for templates.
		Funcs:      []template.FuncMap{unescapeFuncMap}, // Specify helper function maps for templates to access.
		Charset:    "UTF-8",                             // Sets encoding for json and html content-types. Default is "UTF-8".
		IndentJSON: true,                                // Output human readable JSON
		//IndentXML:       true,                    // Output human readable XML
		//HTMLContentType: "application/xhtml+xml", // Output XHTML content type instead of default "text/html"
	}))

	posts = make(map[string]*models.Post, 0)

	staticOptions := martini.StaticOptions{Prefix: "assets"}
	m.Use(martini.Static("assets", staticOptions))
	m.Get("/", IndexHendler)
	m.Get("/write", WriteHendler)
	m.Get("/edit/:id", EditHendler)
	m.Get("/delete/:id", DeleteHendler)
	m.Post("/SavePost", savePostHendler)
	m.Post("/gethtml", getHtmlHendler)

	m.Run()
}
