package auth

type signup struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type login struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}
