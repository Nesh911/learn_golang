package models

type Post struct {
	Id              string
	Title           string
	ContentHtml     string
	ContentMarkdown string
}

func NewPost(id, title, contentHtml, ContentMarkdown string) *Post {
	return &Post{id, title, contentHtml, ContentMarkdown}
}
