package controller_test

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/wlchs/blog/internal/container"
	"github.com/wlchs/blog/internal/controller"
	"github.com/wlchs/blog/internal/errortypes"
	"github.com/wlchs/blog/internal/logger"
	"github.com/wlchs/blog/internal/mocks"
	"github.com/wlchs/blog/internal/test"
	"github.com/wlchs/blog/internal/types"
	"net/http/httptest"
	"testing"
)

// userTestContext contains commonly used services, controllers and other objects relevant for testing the UserController.
type userTestContext struct {
	mockUserService *mocks.MockUserService
	sut             controller.UserController
	ctx             *gin.Context
	rec             *httptest.ResponseRecorder
}

// createUserControllerContext creates the context for testing the UserController and reduces code duplication.
func createUserControllerContext(t *testing.T) *userTestContext {
	t.Helper()

	mockCtrl := gomock.NewController(t)
	mockUserService := mocks.NewMockUserService(mockCtrl)
	cont := container.CreateContainer(logger.CreateLogger(), nil, nil, nil)
	sut := controller.CreateUserController(cont, mockUserService)
	ctx, rec := test.CreateControllerContext()

	return &userTestContext{mockUserService, sut, ctx, rec}
}

// TestUserController_GetUser tests retrieving a user from the blog.
func TestUserController_GetUser(t *testing.T) {
	t.Parallel()
	c := createUserControllerContext(t)

	expectedOutput := types.User{
		UserName: "testAuthor",
		Posts:    []string{"urlHandle1", "urlHandle2"},
	}

	c.ctx.AddParam("userName", expectedOutput.UserName)
	c.mockUserService.EXPECT().GetUser(expectedOutput.UserName).Return(expectedOutput, nil)

	c.sut.GetUser(c.ctx)

	var output types.User
	_ = json.Unmarshal(c.rec.Body.Bytes(), &output)

	assert.Nil(t, c.ctx.Errors, "should complete without error")
	assert.Equal(t, expectedOutput, output, "response body should match")
	assert.Equal(t, 200, c.rec.Code, "incorrect response status")
}

// TestUserController_GetUser_Missing_Username tests retrieving a user from the blog without username.
func TestUserController_GetUser_Missing_Username(t *testing.T) {
	t.Parallel()
	c := createUserControllerContext(t)

	expectedError := errortypes.MissingUsernameError{}

	c.sut.GetUser(c.ctx)

	errors := c.ctx.Errors.Errors()
	assert.Equal(t, 1, len(errors), "expected exactly 1 error")
	assert.Equal(t, expectedError.Error(), errors[0], "incorrect error type")
	assert.Equal(t, 400, c.rec.Code, "incorrect response status")
}

// TestUserController_GetUser_Incorrect_Username tests retrieving a non-existing user from the blog.
func TestUserController_GetUser_Incorrect_Username(t *testing.T) {
	t.Parallel()
	c := createUserControllerContext(t)

	userName := "testAuthor"
	expectedError := errortypes.UserNotFoundError{User: types.User{UserName: userName}}
	c.ctx.AddParam("userName", userName)
	c.mockUserService.EXPECT().GetUser(userName).Return(types.User{}, expectedError)

	c.sut.GetUser(c.ctx)

	errors := c.ctx.Errors.Errors()
	assert.Equal(t, 1, len(errors), "expected exactly 1 error")
	assert.Equal(t, expectedError.Error(), errors[0], "incorrect error type")
	assert.Equal(t, 404, c.rec.Code, "incorrect response status")
}

// TestUserController_GetUser_Unexpected_Error tests handling an unexpected error while retrieving a user from the blog.
func TestUserController_GetUser_Unexpected_Error(t *testing.T) {
	t.Parallel()
	c := createUserControllerContext(t)

	userName := "testAuthor"
	expectedError := errortypes.UnexpectedUserError{User: types.User{UserName: userName}}
	c.ctx.AddParam("userName", userName)
	c.mockUserService.EXPECT().GetUser(userName).Return(types.User{}, fmt.Errorf("unexpected error"))

	c.sut.GetUser(c.ctx)

	errors := c.ctx.Errors.Errors()
	assert.Equal(t, 1, len(errors), "expected exactly 1 error")
	assert.Equal(t, expectedError.Error(), errors[0], "incorrect error type")
	assert.Equal(t, 500, c.rec.Code, "incorrect response status")
}

// TestUserController_GetUsers tests retrieving every user from the blog.
func TestUserController_GetUsers(t *testing.T) {
	t.Parallel()
	c := createUserControllerContext(t)

	expectedOutput := []types.User{
		{
			UserName: "testAuthor",
			Posts:    []string{"urlHandle1", "urlHandle2"},
		},
	}

	c.mockUserService.EXPECT().GetUsers().Return(expectedOutput, nil)

	c.sut.GetUsers(c.ctx)

	var output []types.User
	_ = json.Unmarshal(c.rec.Body.Bytes(), &output)

	assert.Nil(t, c.ctx.Errors, "should complete without error")
	assert.Equal(t, expectedOutput, output, "response body should match")
	assert.Equal(t, 200, c.rec.Code, "incorrect response status")
}

// TestUserController_GetUsers_Unexpected_Error tests handling an unexpected error while retrieving users from the blog.
func TestUserController_GetUsers_Unexpected_Error(t *testing.T) {
	t.Parallel()
	c := createUserControllerContext(t)

	expectedError := errortypes.UnexpectedUserError{}
	c.mockUserService.EXPECT().GetUsers().Return(nil, fmt.Errorf("unexpected error"))

	c.sut.GetUsers(c.ctx)

	errors := c.ctx.Errors.Errors()
	assert.Equal(t, 1, len(errors), "expected exactly 1 error")
	assert.Equal(t, expectedError.Error(), errors[0], "incorrect error type")
	assert.Equal(t, 500, c.rec.Code, "incorrect response status")
}

// TestUserController_UpdateUser tests updating a user's password.
func TestUserController_UpdateUser(t *testing.T) {
	t.Parallel()
	c := createUserControllerContext(t)

	input := types.UserUpdateInput{
		OldPassword: "oldPW",
		NewPassword: "newPW",
	}
	expectedOutput := types.User{
		UserName: "testAuthor",
		Posts:    []string{"urlHandle1", "urlHandle2"},
	}
	mockOld := types.UserLoginInput{
		UserName: expectedOutput.UserName,
		Password: input.OldPassword,
	}
	mockNew := types.UserLoginInput{
		UserName: expectedOutput.UserName,
		Password: input.NewPassword,
	}

	test.MockJsonPost(c.ctx, input)

	c.ctx.AddParam("userName", expectedOutput.UserName)
	c.mockUserService.EXPECT().UpdateUser(&mockOld, &mockNew).Return(expectedOutput, nil)
	c.sut.UpdateUser(c.ctx)

	var output types.User
	_ = json.Unmarshal(c.rec.Body.Bytes(), &output)

	assert.Nil(t, c.ctx.Errors, "should complete without error")
	assert.Equal(t, expectedOutput, output, "response body should match")
	assert.Equal(t, 200, c.rec.Code, "incorrect response status")
}

// TestUserController_UpdateUser_Invalid_Input tests updating a user's password with invalid input.
func TestUserController_UpdateUser_Invalid_Input(t *testing.T) {
	t.Parallel()
	c := createUserControllerContext(t)

	userName := "testAuthor"
	c.ctx.AddParam("userName", userName)

	c.sut.UpdateUser(c.ctx)

	errors := c.ctx.Errors.Errors()
	assert.Equal(t, 1, len(errors), "expected exactly 1 error")
	assert.Equal(t, 400, c.rec.Code, "incorrect response status")
}

// TestUserController_UpdateUser_Incorrect_Username tests updating a user's password with incorrect old password.
func TestUserController_UpdateUser_Incorrect_Username(t *testing.T) {
	t.Parallel()
	c := createUserControllerContext(t)

	userName := "testAuthor"
	input := types.UserUpdateInput{
		OldPassword: "oldPW",
		NewPassword: "newPW",
	}
	mockOld := types.UserLoginInput{
		UserName: userName,
		Password: input.OldPassword,
	}
	mockNew := types.UserLoginInput{
		UserName: userName,
		Password: input.NewPassword,
	}
	expectedError := errortypes.IncorrectUsernameOrPasswordError{}

	test.MockJsonPost(c.ctx, input)

	c.ctx.AddParam("userName", userName)
	c.mockUserService.EXPECT().UpdateUser(&mockOld, &mockNew).Return(types.User{}, expectedError)
	c.sut.UpdateUser(c.ctx)

	errors := c.ctx.Errors.Errors()
	assert.Equal(t, 1, len(errors), "expected exactly 1 error")
	assert.Equal(t, expectedError.Error(), errors[0], "incorrect error type")
	assert.Equal(t, 401, c.rec.Code, "incorrect response status")
}

// TestUserController_UpdateUser_Unexpected_Error tests handling an unexpected error while updating a user's password.
func TestUserController_UpdateUser_Unexpected_Error(t *testing.T) {
	t.Parallel()
	c := createUserControllerContext(t)

	userName := "testAuthor"
	input := types.UserUpdateInput{
		OldPassword: "oldPW",
		NewPassword: "newPW",
	}
	mockOld := types.UserLoginInput{
		UserName: userName,
		Password: input.OldPassword,
	}
	mockNew := types.UserLoginInput{
		UserName: userName,
		Password: input.NewPassword,
	}
	expectedError := errortypes.UnexpectedUserError{User: types.User{UserName: userName}}

	test.MockJsonPost(c.ctx, input)

	c.ctx.AddParam("userName", userName)
	c.mockUserService.EXPECT().UpdateUser(&mockOld, &mockNew).Return(types.User{}, fmt.Errorf("unexpected error"))
	c.sut.UpdateUser(c.ctx)

	errors := c.ctx.Errors.Errors()
	assert.Equal(t, 1, len(errors), "expected exactly 1 error")
	assert.Equal(t, expectedError.Error(), errors[0], "incorrect error type")
	assert.Equal(t, 500, c.rec.Code, "incorrect response status")
}
