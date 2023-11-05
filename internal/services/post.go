package services

import (
	"github.com/wlchs/blog/internal/models"
	"github.com/wlchs/blog/internal/transport/types"
	"math"
)

// pageSize defines the number of entries retrieved on one page
const pageSize = 5

func mapPost(p models.Post) types.Post {
	return types.Post{
		URLHandle:    p.URLHandle,
		Title:        p.Title,
		Author:       p.Author.UserName,
		Summary:      p.Summary,
		Body:         p.Body,
		CreationTime: p.CreatedAt,
	}
}
func mapPostMetadata(p models.Post) types.PostMetadata {
	return types.PostMetadata{
		URLHandle:    p.URLHandle,
		Title:        p.Title,
		Author:       p.Author.UserName,
		Summary:      p.Summary,
		CreationTime: p.CreatedAt,
	}
}

func mapPosts(p []models.Post) []types.PostMetadata {
	var posts []types.PostMetadata

	for _, post := range p {
		posts = append(posts, mapPostMetadata(post))
	}

	return posts
}

func mapPostHandles(p []models.Post) []string {
	var handles []string

	for _, post := range p {
		handles = append(handles, post.URLHandle)
	}

	return handles
}

// GetPosts retrieves posts from the database.
// The page parameter is used for pagination and defines the range of posts to retrieve.
func GetPosts(page int) ([]types.PostMetadata, error) {
	startIndex := (page - 1) * pageSize
	p, err := models.GetPosts(startIndex, pageSize)
	return mapPosts(p), err
}

// CountPostPages counts the number of pages available for pagination
func CountPostPages() (int, error) {
	i, err := models.CountPosts()
	if err != nil {
		return 0, err
	}
	return int(math.Ceil(float64(i) / pageSize)), nil
}

func GetPost(id string) (types.Post, error) {
	p, err := models.GetPost(id)
	return mapPost(p), err
}

func AddPost(post types.Post) (types.Post, error) {
	p, err := models.AddPost(post)
	return mapPost(p), err
}
