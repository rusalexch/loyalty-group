package db

type reward struct {
	ID     string  `db:"id"`
	Type   int     `db:"type"`
	Reward float64 `db:"reward"`
}
