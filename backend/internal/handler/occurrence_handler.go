// backend/internal/handler/occurence_handler.go
package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/saku-730/specimen-web/backend/internal/model"
	"github.com/saku-730/specimen-web/backend/internal/service"
)

type OccurrenceHandler struct {
	occurrenceService service.OccurrenceService
}

func NewOccurrenceHandler(occurrenceService service.OccurrenceService) *OccurrenceHandler {
	return &OccurrenceHandler{occurrenceService: occurrenceService}
}

// RegisterOccurrenceRoutes はルーターに発生情報関連のエンドポイントを登録するのだ
func (h *OccurrenceHandler) RegisterOccurrenceRoutes(router *gin.RouterGroup) {
	occurrences := router.Group("/occurrences")
	{
		occurrences.GET("", h.GetAllOccurrences)
		occurrences.POST("/search", h.SearchOccurrences) // 条件検索用のエンドポイント
		occurrences.GET("/:id", h.GetOccurrenceByID)
		occurrences.POST("", h.CreateOccurrence)
		occurrences.PUT("/:id", h.UpdateOccurrence)
		occurrences.DELETE("/:id", h.DeleteOccurrence)
	}
}

func (h *OccurrenceHandler) GetOccurrenceByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "無効なIDです"})
		return
	}

	occurrence, err := h.occurrenceService.GetOccurrenceByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "発生情報が見つかりません"})
		return
	}

	c.JSON(http.StatusOK, occurrence)
}

func (h *OccurrenceHandler) GetAllOccurrences(c *gin.Context) {
	occurrences, err := h.occurrenceService.GetAllOccurrences()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "発生情報の取得に失敗しました"})
		return
	}
	c.JSON(http.StatusOK, occurrences)
}

func (h *OccurrenceHandler) SearchOccurrences(c *gin.Context) {
	var conditions model.Occurrence
	if err := c.ShouldBindJSON(&conditions); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "リクエストの形式が正しくありません"})
		return
	}

	occurrences, err := h.occurrenceService.FindOccurrencesByConditions(&conditions)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "検索に失敗しました"})
		return
	}

	c.JSON(http.StatusOK, occurrences)
}

func (h *OccurrenceHandler) CreateOccurrence(c *gin.Context) {
	var req service.CreateOccurrenceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "リクエストの形式が正しくありません"})
		return
	}

	createdOccurrence, err := h.occurrenceService.CreateOccurrence(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "発生情報の作成に失敗しました"})
		return
	}

	c.JSON(http.StatusCreated, createdOccurrence)
}

func (h *OccurrenceHandler) UpdateOccurrence(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "無効なIDです"})
		return
	}

	var req service.UpdateOccurrenceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "リクエストの形式が正しくありません"})
		return
	}

	updatedOccurrence, err := h.occurrenceService.UpdateOccurrence(uint(id), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "発生情報の更新に失敗しました"})
		return
	}

	c.JSON(http.StatusOK, updatedOccurrence)
}

func (h *OccurrenceHandler) DeleteOccurrence(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "無効なIDです"})
		return
	}

	if err := h.occurrenceService.DeleteOccurrence(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "発生情報の削除に失敗しました"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "発生情報を削除しました"})
}
