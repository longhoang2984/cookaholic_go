package http

import (
	"cookaholic/internal/common"
	"cookaholic/internal/interfaces"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type UserFollowerHandler struct {
	userFollowerService interfaces.UserFollowerService
}

func NewUserFollowerHandler(userFollowerService interfaces.UserFollowerService) *UserFollowerHandler {
	return &UserFollowerHandler{
		userFollowerService: userFollowerService,
	}
}

// FollowUser handles the request to follow a user
func (h *UserFollowerHandler) FollowUser(c *gin.Context) {
	// Get current user ID from context
	currentUserID, errResp := AuthorizedPermission(c)
	if errResp != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": errResp.Message})
		return
	}

	// Get target user ID from URL
	targetUserIDStr := c.Param("id")
	targetUserID, err := uuid.Parse(targetUserIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	// Follow the user
	err = h.userFollowerService.FollowUser(c.Request.Context(), *currentUserID, targetUserID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Successfully followed user"})
}

// UnfollowUser handles the request to unfollow a user
func (h *UserFollowerHandler) UnfollowUser(c *gin.Context) {
	// Get current user ID from context
	currentUserID, errResp := AuthorizedPermission(c)
	if errResp != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": errResp.Message})
		return
	}

	// Get target user ID from URL
	targetUserIDStr := c.Param("id")
	targetUserID, err := uuid.Parse(targetUserIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	// Unfollow the user
	err = h.userFollowerService.UnfollowUser(c.Request.Context(), *currentUserID, targetUserID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Successfully unfollowed user"})
}

// GetFollowers handles the request to get a user's followers
func (h *UserFollowerHandler) GetFollowers(c *gin.Context) {
	// Get user ID from URL
	userIDStr := c.Param("id")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	// Get pagination parameters
	cursorStr := c.Query("cursor")
	var cursor *time.Time
	if cursorStr != "" {
		cursorTime, err := common.CursorToTime(cursorStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid cursor"})
			return
		}
		cursor = cursorTime
	}

	limit, err := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if err != nil || limit < 1 || limit > 100 {
		limit = 10
	}

	// Get followers with user info
	followers, nextCursorTime, err := h.userFollowerService.GetFollowers(c.Request.Context(), userID, cursor, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve followers"})
		return
	}

	// Convert next cursor time to UUID string
	nextCursor := ""
	if nextCursorTime != nil {
		nextCursor = common.TimeToCursor(nextCursorTime)
	}

	c.JSON(http.StatusOK, gin.H{
		"data": followers,
		"meta": gin.H{
			"next_cursor": nextCursor,
			"limit":       limit,
		},
	})
}

// GetFollowing handles the request to get users whom a user follows
func (h *UserFollowerHandler) GetFollowing(c *gin.Context) {
	// Get user ID from URL
	userIDStr := c.Param("id")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	// Get pagination parameters
	cursorStr := c.Query("cursor")
	var cursor *time.Time
	if cursorStr != "" {
		cursorTime, err := common.CursorToTime(cursorStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid cursor"})
			return
		}
		cursor = cursorTime
	}

	limit, err := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if err != nil || limit < 1 || limit > 100 {
		limit = 10
	}

	// Get following with user info
	following, nextCursorTime, err := h.userFollowerService.GetFollowing(c.Request.Context(), userID, cursor, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve following users"})
		return
	}

	// Convert next cursor time to UUID string
	nextCursor := ""
	if nextCursorTime != nil {
		nextCursor = common.TimeToCursor(nextCursorTime)
	}

	c.JSON(http.StatusOK, gin.H{
		"data": following,
		"meta": gin.H{
			"next_cursor": nextCursor,
			"limit":       limit,
		},
	})
}

// IsFollowing handles the request to check if a user is following another user
func (h *UserFollowerHandler) IsFollowing(c *gin.Context) {
	// Get current user ID from context
	currentUserID, errResp := AuthorizedPermission(c)
	if errResp != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": errResp.Message})
		return
	}

	// Get target user ID from URL
	targetUserIDStr := c.Param("id")
	targetUserID, err := uuid.Parse(targetUserIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	// Check if following
	isFollowing, err := h.userFollowerService.IsFollowing(c.Request.Context(), *currentUserID, targetUserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check following status"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"is_following": isFollowing})
}

// GetFollowersCount handles the request to get a user's follower count
func (h *UserFollowerHandler) GetFollowersCount(c *gin.Context) {
	// Get user ID from URL
	userIDStr := c.Param("id")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	// Get follower count
	count, err := h.userFollowerService.GetFollowersCount(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve follower count"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"count": count})
}

// GetFollowingCount handles the request to get the count of users a user is following
func (h *UserFollowerHandler) GetFollowingCount(c *gin.Context) {
	// Get user ID from URL
	userIDStr := c.Param("id")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	// Get following count
	count, err := h.userFollowerService.GetFollowingCount(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve following count"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"count": count})
}
