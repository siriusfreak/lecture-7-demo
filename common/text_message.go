package common

type TextMessage struct {
	ID 				int64 	`json:"id"`
	Text 			string	`json:"text"`
	CallbackUrl 	string	`json:"callback_url"`
}
