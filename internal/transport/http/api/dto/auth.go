package dto

type LoginReq struct {
	Login    string  `json:"login"`
	Password string  `json:"password"`
	TTL      *string `json:"ttl,omitempty"`
}
