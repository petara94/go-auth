package dto

type Session struct {
	Token  string `json:"token"`
	UserID uint64 `json:"user_id"`
}

type Auth struct {
	Login    string `json:"user_id"`
	Password string `json:"password"`
}
