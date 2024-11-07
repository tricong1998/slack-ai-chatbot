package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sotatek-dev/hyper-automation-chatbot/internal/dto"
	"github.com/sotatek-dev/hyper-automation-chatbot/internal/services"
)

type UIPathHandler struct {
	uiPathService *services.UIPathService
}

func NewUIPathHandler(uiPathService *services.UIPathService) *UIPathHandler {
	return &UIPathHandler{uiPathService: uiPathService}
}

func (h *UIPathHandler) GreetingNewEmployee(c *gin.Context) {
	var dto dto.UIPathGreetingNewEmployee
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	data, err := h.uiPathService.GreetingNewEmployee(dto)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, data)
}

func (h *UIPathHandler) GetJobDetails(c *gin.Context) {
	jobID, err := strconv.Atoi(c.Param("jobID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	data, err := h.uiPathService.GetJobDetails(jobID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, data)
}
