package auth

type Config struct {
	UserService userService
	JwtSecret   string
}

type authModule struct {
	userService userService
	jwtSecret   string
}

type signup struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type login struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}
