package dto

import "time"

type Session struct {
	Token  string     `json:"token"`
	UserID uint64     `json:"user_id"`
	Expr   *time.Time `json:"expr"`
}

type Auth struct {
	Login    string     `json:"login"`
	Password string     `json:"password"`
	TTL      *time.Time `json:"ttl,omitempty"`
}
