package dto

// ChangePassword change password
type ChangePassword struct {
	OldPassword string `json:"old_password"`
	NewPassword string `json:"new_password"`
}

// PublicUser public user
type PublicUser struct {
	ID    uint64 `json:"id"`
	Login string `json:"login"`
}

// LoginPassword login and password
type LoginPassword struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}
