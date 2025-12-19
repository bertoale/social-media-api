package comment

import (
	"go-sosmed/pkg/response"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Controller struct {
	service Service
}

// ==========================================
// Helper to get user ID from context
// ==========================================
func GetUserIDFromContext(c *gin.Context) (uint, bool) {
	uid, exists := c.Get("userID")
	if !exists {
		return 0, false
	}
	return uid.(uint), true
}
func ParsePostID(c *gin.Context) (uint, error) {
	postIDParam := c.Param("post_id")
	postID, err := strconv.ParseUint(postIDParam, 10, 32)
	if err != nil {
		return 0, err
	}
	return uint(postID), nil
}
func ParseCommentID(c *gin.Context) (uint, error) {
	idParam := c.Param("comment_id")
	commentID, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		return 0, err
	}
	return uint(commentID), nil
}

// ==========================================
// Controller Methods
// ==========================================
// CreateComment godoc
// @Summary Create a comment
// @Description Add a new comment to a post
// @Tags Comment
// @Accept json
// @Produce json
// @Param post_id path int true "Post ID"
// @Param data body CommentRequest true "Comment data"
// @Security BearerAuth
// @Success 201 {object} response.SuccessResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 401 {object} response.ErrorResponse
// @Router /api/posts/{post_id}/comments [post]
func (ctrl *Controller) CreateComment(c *gin.Context) {
	var req CommentRequest
	if err := c.ShouldBind(&req); err != nil {
		response.Error(c, 400, "invalid request: "+err.Error())
		return
	}

	userID, ok := GetUserIDFromContext(c)
	if !ok {
		response.Error(c, 401, "user not authenticated")
		return
	}

	postID, err := ParsePostID(c)
	if err != nil {
		response.Error(c, 400, "invalid post id")
		return
	}

	comment := &Comment{
		PostID:  uint(postID),
		UserID:  userID,
		Content: req.Content,
	}

	saved, err := ctrl.service.CreateComment(comment)
	if err != nil {
		response.Error(c, 400, err.Error())
		return
	}

	response.Success(c, 201, "comment created successfully", saved)
}

// ReplyToComment godoc
// @Summary Reply to a comment
// @Description Reply to an existing comment (supports nested replies up to 2 levels)
// @Tags Comment
// @Accept json
// @Produce json
// @Param post_id path int true "Post ID"
// @Param comment_id path int true "Parent Comment ID"
// @Param data body ReplyCommentRequest true "Reply data"
// @Security BearerAuth
// @Success 201 {object} response.SuccessResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 401 {object} response.ErrorResponse
// @Failure 404 {object} response.ErrorResponse
// @Router /api/posts/{post_id}/comments/{comment_id}/reply [post]
func (ctrl *Controller) ReplyToComment(c *gin.Context) {
	// Parent comment ID
	parentID, err := ParseCommentID(c)
	if err != nil {
		response.Error(c, 400, "invalid parent comment ID")
		return
	}

	// Bind request
	var req ReplyCommentRequest
	if err := c.ShouldBind(&req); err != nil {
		response.Error(c, 400, "invalid request: "+err.Error())
		return
	}

	// Ambil user ID dari token
	userID, ok := GetUserIDFromContext(c)
	if !ok {
		response.Error(c, 401, "user not authenticated")
		return
	}

	// Ambil parent comment lengkap
	parent, err := ctrl.service.GetByID(uint(parentID))
	if err != nil {
		response.Error(c, 404, "parent comment not found")
		return
	}

	// **HARUS** isi post_id → untuk mencegah FK error
	reply := &Comment{
		PostID:  parent.PostID, // FIX
		UserID:  userID,
		Content: req.Content,
	}

	// ---- LEVEL RULE ----
	// Level 1 → reply ke root
	// Level 2 → reply ke reply tetapi parent tetap root
	if parent.ParentID == nil {
		reply.ParentID = &parent.ID
		reply.ReplyToUserID = &parent.UserID
	} else {
		reply.ParentID = parent.ParentID
		reply.ReplyToUserID = &parent.UserID
	}

	// Save reply
	saved, err := ctrl.service.CreateComment(reply)
	if err != nil {
		response.Error(c, 500, err.Error())
		return
	}

	// Gunakan serializer
	response.Success(c, 201, "reply created successfully", ToCommentResponse(saved))
}

// UpdateComment godoc
// @Summary Update a comment
// @Description Update a comment (only author can update)
// @Tags Comment
// @Accept json
// @Produce json
// @Param comment_id path int true "Comment ID"
// @Param data body UpdateCommentRequest true "Update data"
// @Security BearerAuth
// @Success 200 {object} response.SuccessResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 401 {object} response.ErrorResponse
// @Failure 403 {object} response.ErrorResponse
// @Router /api/comments/{comment_id} [put]
func (ctrl *Controller) UpdateComment(c *gin.Context) {
	commentID, err := ParseCommentID(c)
	if err != nil {
		response.Error(c, 400, "invalid comment ID")
		return
	}

	var req UpdateCommentRequest
	if err := c.ShouldBind(&req); err != nil {
		response.Error(c, 400, "invalid request: "+err.Error())
		return
	}

	userID, ok := GetUserIDFromContext(c)
	if !ok {
		response.Error(c, 401, "user not authenticated")
		return
	}

	updated, err := ctrl.service.UpdateComment(userID, uint(commentID), req)
	if err != nil {
		response.Error(c, 400, err.Error())
		return
	}

	response.Success(c, 200, "comment updated successfully", updated)
}

// DeleteComment godoc
// @Summary Delete a comment
// @Description Delete a comment (only author can delete)
// @Tags Comment
// @Accept json
// @Produce json
// @Param comment_id path int true "Comment ID"
// @Security BearerAuth
// @Success 200 {object} response.SuccessResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 401 {object} response.ErrorResponse
// @Failure 403 {object} response.ErrorResponse
// @Router /api/comments/{comment_id} [delete]
func (ctrl *Controller) DeleteComment(c *gin.Context) {
	commentID, err := ParseCommentID(c)
	if err != nil {
		response.Error(c, 400, "invalid comment ID")
		return
	}

	userID, ok := GetUserIDFromContext(c)
	if !ok {
		response.Error(c, 401, "user not authenticated")
		return
	}

	err = ctrl.service.DeleteComment(userID, uint(commentID))
	if err != nil {
		response.Error(c, 400, err.Error())
		return
	}

	response.Success(c, 200, "comment deleted successfully", nil)
}

// GetCommentTree godoc
// @Summary Get comment tree for a post
// @Description Retrieve all comments for a post in hierarchical tree structure with nested replies
// @Tags Comment
// @Accept json
// @Produce json
// @Param post_id path int true "Post ID"
// @Security BearerAuth
// @Success 200 {object} response.SuccessResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 401 {object} response.ErrorResponse
// @Router /api/posts/{post_id}/comments [get]
func (ctrl *Controller) GetCommentTree(c *gin.Context) {
	postID, err := ParsePostID(c)
	if err != nil {
		response.Error(c, 400, "invalid post ID")
		return
	}

	tree, err := ctrl.service.GetCommentTree(uint(postID))
	if err != nil {
		response.Error(c, 400, err.Error())
		return
	}

	// MAP model -> DTO
	var resp []CommentResponse
	for i := range tree {
		resp = append(resp, ToCommentResponse(&tree[i]))
	}

	response.Success(c, 200, "comments retrieved successfully", resp)
}

// GetReplies godoc
// @Summary Get replies to a comment
// @Description Retrieve all direct replies to a specific comment
// @Tags Comment
// @Accept json
// @Produce json
// @Param post_id path int true "Post ID"
// @Param comment_id path int true "Comment ID"
// @Security BearerAuth
// @Success 200 {object} response.SuccessResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 401 {object} response.ErrorResponse
// @Router /api/posts/{post_id}/comments/{comment_id}/replies [get]
func (ctrl *Controller) GetReplies(c *gin.Context) {
	commentID, err := ParseCommentID(c)
	if err != nil {
		response.Error(c, 400, "invalid comment ID")
		return
	}

	replies, err := ctrl.service.GetReplies(uint(commentID))
	if err != nil {
		response.Error(c, 400, err.Error())
		return
	}

	// FIX → convert model → DTO
	var resp []CommentResponse
	for _, reply := range replies {
		resp = append(resp, ToCommentResponse(&reply))
	}

	response.Success(c, 200, "replies retrieved successfully", resp)
}

// GetCommentCount godoc
// @Summary Get total comments of a post
// @Description Get number of comments (including replies) for a post
// @Tags Comment
// @Accept json
// @Produce json
// @Param post_id path int true "Post ID"
// @Success 200 {object} response.SuccessResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /api/posts/{post_id}/comments/count [get]
func (ctrl *Controller) GetCommentCount(c *gin.Context) {
	postID, err := ParsePostID(c)
	if err != nil {
		response.Error(c, 400, "invalid post ID")
		return
	}

	total, err := ctrl.service.GetCommentCount(postID)
	if err != nil {
		response.Error(c, 500, err.Error())
		return
	}

	response.Success(c, 200, "comment count retrieved", gin.H{
		"post_id":       postID,
		"total_comment": total,
	})
}



func NewController(service Service) *Controller {
	return &Controller{service: service}
}
