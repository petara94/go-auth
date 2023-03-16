package dto

type Message struct {
	Message string `json:"message"`
}

// SuccessMessage func get success message
func SuccessMessage() *Message {
	return &Message{Message: "success"}
}
