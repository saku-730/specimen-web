
// backend/internal/handler/specimen_handler.go
package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/saku-730/specimen-web/backend/internal/service"
)

type SpecimenHandler struct {
	specimenService service.SpecimenService
}

func NewSpecimenHandler(specimenService service.SpecimenService) *SpecimenHandler {
	return &SpecimenHandler{specimenService: specimenService}
}

// RegisterSpecimenRoutes はルーターに標本関連のエンドポイントを登録するのだ
func (h *SpecimenHandler) RegisterSpecimenRoutes(router *gin.RouterGroup) {
	specimens := router.Group("/specimens")
	{
		specimens.GET("", h.GetAllSpecimens)
		specimens.GET("/:id", h.GetSpecimenByID)
		specimens.POST("", h.CreateSpecimen)
		specimens.PUT("/:id", h.UpdateSpecimen)
		specimens.DELETE("/:id", h.DeleteSpecimen)
	}
}

func (h *SpecimenHandler) GetSpecimenByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "無効なIDです"})
		return
	}

	specimen, err := h.specimenService.GetSpecimenByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "標本が見つかりません"})
		return
	}

	c.JSON(http.StatusOK, specimen)
}

func (h *SpecimenHandler) GetAllSpecimens(c *gin.Context) {
	specimens, err := h.specimenService.GetAllSpecimens()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "標本の取得に失敗しました"})
		return
	}
	c.JSON(http.StatusOK, specimens)
}

func (h *SpecimenHandler) CreateSpecimen(c *gin.Context) {
	var req service.CreateSpecimenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "リクエストの形式が正しくありません"})
		return
	}

	createdSpecimen, err := h.specimenService.CreateSpecimen(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "標本の作成に失敗しました"})
		return
	}

	c.JSON(http.StatusCreated, createdSpecimen)
}

func (h *SpecimenHandler) UpdateSpecimen(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "無効なIDです"})
		return
	}

	var req service.UpdateSpecimenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "リクエストの形式が正しくありません"})
		return
	}

	updatedSpecimen, err := h.specimenService.UpdateSpecimen(uint(id), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "標本の更新に失敗しました"})
		return
	}

	c.JSON(http.StatusOK, updatedSpecimen)
}

func (h *SpecimenHandler) DeleteSpecimen(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "無効なIDです"})
		return
	}

	if err := h.specimenService.DeleteSpecimen(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "標本の削除に失敗しました"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "標本を削除しました"})
}
