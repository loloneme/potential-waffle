package models

type Reviewer struct {
	PullRequestID string `db:"pr_id"`
	ReviewerID    string `db:"reviewer_id"`
}
