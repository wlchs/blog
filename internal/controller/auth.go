package controller

import (
	"github.com/wlchs/blog/internal/container"
	"github.com/wlchs/blog/internal/errortypes"
	"github.com/wlchs/blog/internal/types"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/wlchs/blog/internal/services"
)

// AuthController interface defining authentication-related methods to handler HTTP requests.
type AuthController interface {
	Login(c *gin.Context)
	Protect(c *gin.Context)
}

// authController is a concrete implementation of the AuthController interface.
type authController struct {
	cont        container.Container
	userService services.UserService
}

// CreateAuthController instantiates the AuthController using the application container.
func CreateAuthController(cont container.Container, userService services.UserService) AuthController {
	return &authController{cont, userService}
}

// Login middleware. Top level handler of /login POST requests.
func (auth authController) Login(c *gin.Context) {
	userService := auth.userService

	var u types.UserLoginInput
	if err := c.BindJSON(&u); err != nil {
		return
	}

	token, err := userService.AuthenticateUser(&u)

	if err != nil {
		_ = c.AbortWithError(http.StatusUnauthorized, err)
		return
	}

	c.Header("X-Auth-Token", token)
	c.Status(http.StatusOK)
}

// Protect middleware. Can be used before any middleware to make sure only authenticated users are able to use an endpoint.
func (auth authController) Protect(c *gin.Context) {
	jwtUtils := auth.cont.GetJWTUtils()
	token := c.Request.Header.Get("X-Auth-Token")

	if token == "" {
		_ = c.AbortWithError(http.StatusUnauthorized, errortypes.MissingAuthTokenError{})
	} else if u, err := jwtUtils.ParseJWT(token); err == nil {
		c.Set("user", u)
		c.Next()
	} else {
		_ = c.AbortWithError(http.StatusUnauthorized, errortypes.InvalidAuthTokenError{})
	}
}
