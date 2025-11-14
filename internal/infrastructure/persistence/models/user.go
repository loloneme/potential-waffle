package models

type User struct {
	ID       string `db:"user_id"`
	Username string `db:"username"`
	IsActive bool   `db:"is_active"`
	TeamName string `db:"team_name"`
}
