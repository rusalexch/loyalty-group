package app

const (
	ContentType = "Content-Type"
	AppJSON     = "application/json"
	Text        = "text/plain"
	AuthHeader  = "Authorization"
)

type ContextKey int

const (
	UserKey ContextKey = iota
)
