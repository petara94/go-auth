package dto

import "time"

type Session struct {
	Token  string     `json:"token"`
	UserID uint64     `json:"user_id"`
	Expr   *time.Time `json:"expr,omitempty"`
}

type Auth struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}
