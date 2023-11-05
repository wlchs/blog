package rest

import (
	"github.com/gin-gonic/gin"
	"github.com/wlchs/blog/internal/errortypes"
	"github.com/wlchs/blog/internal/services"
	"github.com/wlchs/blog/internal/transport/types"
	"net/http"
	"strconv"
)

// getPostMiddleware is the middleware responsible for querying posts. If the page parameter is provided, it is used for pagination.
func getPostsMiddleware(c *gin.Context) {
	page := 1
	p, found := c.GetQuery("page")

	if found {
		pageInt, err := strconv.Atoi(p)

		if err != nil {
			_ = c.AbortWithError(http.StatusBadRequest, err)
			return
		}

		page = pageInt
	}

	posts, err := services.GetPosts(page)
	if err != nil {
		_ = c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	c.IndentedJSON(http.StatusOK, posts)
}

func getPostMiddleware(c *gin.Context) {
	id, found := c.Params.Get("id")

	if !found {
		c.AbortWithStatusJSON(http.StatusBadRequest, "No id provided!")
		return
	}

	post, err := services.GetPost(id)
	if err != nil {
		_ = c.AbortWithError(http.StatusNotFound, err)
		return
	}

	c.IndentedJSON(http.StatusOK, post)
}

func addPostMiddleware(c *gin.Context) {
	var post types.Post

	if err := c.BindJSON(&post); err != nil {
		_ = c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	// Set author from context
	post.Author = c.GetString("user")
	post, err := services.AddPost(post)

	switch err.(type) {
	case nil:
		c.IndentedJSON(http.StatusCreated, post)

	case errortypes.DuplicateElementError:
		_ = c.AbortWithError(http.StatusConflict, err)

	default:
		_ = c.AbortWithError(http.StatusBadRequest, err)
	}
}
