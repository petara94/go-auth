package dto

type UserGroup struct {
	ID      uint64 `json:"id"`
	Name    string `json:"name"`
	IsAdmin bool   `json:"is_admin"`
}
