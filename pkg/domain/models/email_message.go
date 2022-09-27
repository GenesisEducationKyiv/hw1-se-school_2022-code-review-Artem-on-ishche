package models

type EmailMessage struct {
	Title string
	Body  string
}

func NewEmail(title, body string) *EmailMessage {
	return &EmailMessage{
		Title: title,
		Body:  body,
	}
}
