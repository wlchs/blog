package types

import "time"

type PostMetadata struct {
	URLHandle    string    `json:"urlHandle"`
	Title        string    `json:"title"`
	Author       string    `json:"author"`
	Summary      string    `json:"summary"`
	CreationTime time.Time `json:"creationTime"`
}

type Post struct {
	URLHandle    string    `json:"urlHandle"`
	Title        string    `json:"title"`
	Author       string    `json:"author"`
	Summary      string    `json:"summary"`
	Body         string    `json:"body"`
	CreationTime time.Time `json:"creationTime"`
}
