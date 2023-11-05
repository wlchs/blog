package types

import "time"

type PostMetadata struct {
	URLHandle    string    `json:"urlHandle"`
	Title        string    `json:"title"`
	Author       string    `json:"author"`
	Summary      string    `json:"summary"`
	CreationTime time.Time `json:"creationTime"`
}

type PostsQueryData struct {
	Posts []PostMetadata `json:"posts"`
	Pages int            `json:"pages"`
}

type Post struct {
	URLHandle    string    `json:"urlHandle"`
	Title        string    `json:"title"`
	Author       string    `json:"author"`
	Summary      string    `json:"summary"`
	Body         string    `json:"body"`
	CreationTime time.Time `json:"creationTime"`
}
