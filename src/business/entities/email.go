package entities

type Attachments struct {
	Base64Content string
	Filename      string
	Size          int64
	ContentType   string
}

type Email struct {
	Subject     string
	From        string
	To          string
	PlainText   string
	HTML        string
	Attachments []Attachments
}
