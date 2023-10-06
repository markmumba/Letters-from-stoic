package main

type Blog struct {
	Id        string `json:id`
	Image     string `json:image`
	Title     string `json:title`
	ShortText string `json:shorttext`
	LongText  string `json:longtext`
	Date      string `json:date`
	Comments  []Comment
}

type Comment struct {
	Id     string `json:id`
	BlogId string `json:blog_id`
	Name   string `json:name`
	Email  string `json:email`
	Text   string `json:text`
	Date   string `json:date`
}

func (b Blog) ShortenText() string {
	char := 0
	for i := range b.ShortText {
		char++
		if char > 200 {
			return b.ShortText[:i] + `......`
		}
	}
	return b.ShortText
}
