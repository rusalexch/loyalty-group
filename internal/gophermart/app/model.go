package app

type User struct {
	ID       int    `json:"id" db:"id"`
	Login    string `json:"login" db:"login"`
	Password string `json:"-" db:"password"`
}

type CreateUser struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}
