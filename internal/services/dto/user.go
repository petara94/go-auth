package dto

type User struct {
	ID            uint64 `json:"id"`
	Login         string `json:"login"`
	Password      string `json:"password"`
	CheckPassword bool   `json:"check_password"`
	IsAdmin       bool   `json:"is_admin"`
	IsBlocked     bool   `json:"is_blocked"`
}
