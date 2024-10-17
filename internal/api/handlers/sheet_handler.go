package handlers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sotatek-dev/hyper-automation-chatbot/internal/dto"
	"github.com/sotatek-dev/hyper-automation-chatbot/internal/services"
)

type SheetHandler struct {
	service *services.GSheetService
}

func NewSheetHandler(service *services.GSheetService) *SheetHandler {
	return &SheetHandler{service: service}
}

func (h *SheetHandler) ReadCandidateOffer(c *gin.Context) {
	var req dto.ReadSheetCandidateOffer
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	candidates, err := h.service.ReadCandidateOffer(req.SheetUrl)
	if err != nil {
		fmt.Println(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, candidates)
}

func (h *SheetHandler) CreateNewSheetInSharedDrive(c *gin.Context) {
	var req dto.CreateNewSheet
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	spreadsheet, err := h.service.CreateNewSheetInSharedDrive(req.SheetName, services.SharedDriveFolderId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, spreadsheet)
}

func (h *SheetHandler) HandleFileCandidateOffer(c *gin.Context) {
	var req dto.ReadSheetCandidateOffer
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	spreadsheet, err := h.service.HandleFileCandidateOffer(req.SheetUrl)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, spreadsheet)
}
