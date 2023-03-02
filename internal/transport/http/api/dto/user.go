package dto

type User struct {
	ID          uint64 `json:"id"`
	Login       string `json:"login"`
	Password    string `json:"password"`
	UserGroupID uint64 `json:"user_group_id"`
}
