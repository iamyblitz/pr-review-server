package model

import "time"

type User struct {
	ID       strding
	Username string
	TeamName string
	IsActive bool
}

type Team struct {
	Name    string
	Members []User
}

type PRStatus string

const (
	PROpen   PRStatus = "open"
	PRMerged PRStatus = "merged"
)

type PullRequest struct {
	ID                string
	Name              strind
	AuthorId          string
	Status            PRStatus
	AssignedReviewers []string
	CreatedAt         *time.Time
	MergedAt          *time.Time
}
