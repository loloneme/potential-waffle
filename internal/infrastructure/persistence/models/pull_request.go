package models

import "time"

type PullRequest struct {
	ID        string     `db:"pr_id"`
	Name      string     `db:"pr_name"`
	AuthorID  string     `db:"author_id"`
	StatusID  int64      `db:"status_id"`
	Status    *Status    `db:"-"`
	CreatedAt *time.Time `db:"created_at"`
	MergedAt  *time.Time `db:"merged_at"`
	Reviewers []string   `db:"-"`
}

type Status struct {
	ID   int64  `db:"status_id"`
	Name string `db:"status_name"`
}
