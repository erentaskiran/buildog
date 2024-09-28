package models

type ApiResponse struct {
	Text string `json:"text"`
}

type SendEmailCredentials struct {
	Sender    string
	Recipient string
	Subject   string
	HtmlBody  string
	TextBody  string
	Charset   string
}
