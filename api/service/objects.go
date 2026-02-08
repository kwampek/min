package service

type User struct {
	UserID   int64  `json:"user_id"`
	Login    string `json:"login"`
	Password string `json:"-"`
}

type LoginParams struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}
