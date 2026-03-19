package handlers

import (
	"contract-manage/models"
	"contract-manage/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type WorkflowHandler struct {
	workflowService *services.WorkflowService
}

func NewWorkflowHandler(db *gorm.DB) *WorkflowHandler {
	return &WorkflowHandler{
		workflowService: services.NewWorkflowService(db),
	}
}

func (h *WorkflowHandler) GetWorkflow(c *gin.Context) {
	contractID, err := strconv.ParseUint(c.Param("contract_id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid contract ID"})
		return
	}

	workflow, err := h.workflowService.GetWorkflowByContractID(contractID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Workflow not found"})
		return
	}

	c.JSON(http.StatusOK, workflow)
}

func (h *WorkflowHandler) CreateWorkflow(c *gin.Context) {
	var input struct {
		ContractID  uint64 `json:"contract_id" binding:"required"`
		CreatorRole string `json:"creator_role" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	workflow, err := h.workflowService.CreateWorkflow(input.ContractID, input.CreatorRole)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create workflow"})
		return
	}

	c.JSON(http.StatusCreated, workflow)
}

func (h *WorkflowHandler) Approve(c *gin.Context) {
	var input struct {
		WorkflowID uint64 `json:"workflow_id" binding:"required"`
		Level      int    `json:"level" binding:"required"`
		Comment    string `json:"comment"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID := c.GetUint("user_id")

	err := h.workflowService.Approve(input.WorkflowID, 0, input.Level, uint64(userID), input.Comment)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to approve"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Approved successfully"})
}

func (h *WorkflowHandler) Reject(c *gin.Context) {
	var input struct {
		WorkflowID uint64 `json:"workflow_id" binding:"required"`
		Level      int    `json:"level" binding:"required"`
		Comment    string `json:"comment" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID := c.GetUint("user_id")

	err := h.workflowService.Reject(input.WorkflowID, input.Level, uint64(userID), input.Comment)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to reject"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Rejected successfully"})
}

func (h *WorkflowHandler) GetMyPendingApproval(c *gin.Context) {
	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
		return
	}

	approvals, err := h.workflowService.GetPendingApprovals(string(user.(*models.User).Role))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get pending approvals"})
		return
	}

	c.JSON(http.StatusOK, approvals)
}
