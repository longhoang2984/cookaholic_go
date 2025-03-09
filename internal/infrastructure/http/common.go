package http

import (
	"github.com/gin-gonic/gin"
)

type ErrorResponse struct {
	Message string `json:"message"`
}

func NewErrorResponse(message string) *ErrorResponse {
	return &ErrorResponse{
		Message: message,
	}
}

func AuthorizedPermission(c *gin.Context) (*uint, *ErrorResponse) {
	userID, exists := c.Get("user_id")
	if !exists {
		return nil, NewErrorResponse("unauthorized")
	}

	uid, ok := userID.(uint)
	if !ok {
		return nil, NewErrorResponse("invalid user ID format")
	}

	return &uid, nil
}
