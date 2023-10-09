package main

type Blog struct {
	Id        string
	Image     string
	Title     string
	ShortText string
	LongText  string
	Date      string
	Comments  []Comment
	Session   Session
}

type Comment struct {
	Id     string
	BlogId string
	Name   string
	Email  string
	Text   string
	Date   string
}

type User struct {
	Id   int
	Name string
}

type Session struct {
	Id              int
	Authenticated   bool
	Unauthenticated bool
	User            User
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
