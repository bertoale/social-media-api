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
func ParseBlogID(c *gin.Context) (uint, error) {
	blogIDParam := c.Param("blog_id")
	blogID, err := strconv.ParseUint(blogIDParam, 10, 32)
	if err != nil {
		return 0, err
	}
	return uint(blogID), nil
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

	blogID, err := ParseBlogID(c)
	if err != nil {
		response.Error(c, 400, "invalid blog id")
		return
	}

	comment := &Comment{
		BlogID:  uint(blogID),
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

	// **HARUS** isi blog_id → untuk mencegah FK error
	reply := &Comment{
		BlogID:  parent.BlogID, // FIX
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

func (ctrl *Controller) GetCommentTree(c *gin.Context) {
	blogID, err := ParseBlogID(c)
	if err != nil {
		response.Error(c, 400, "invalid blog ID")
		return
	}

	tree, err := ctrl.service.GetCommentTree(uint(blogID))
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

func NewController(service Service) *Controller {
	return &Controller{service: service}
}
