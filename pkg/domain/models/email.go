package models

type Email struct {
	Title string
	Body  string
}

func NewEmail(title, body string) *Email {
	return &Email{
		Title: title,
		Body:  body,
	}
}
