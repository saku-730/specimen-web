
// backend/internal/handler/identification_handler.go
package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/saku-730/specimen-web/backend/internal/service"
)

type IdentificationHandler struct {
	identificationService service.IdentificationService
}

func NewIdentificationHandler(identificationService service.IdentificationService) *IdentificationHandler {
	return &IdentificationHandler{identificationService: identificationService}
}

// RegisterIdentificationRoutes はルーターに同定情報関連のエンドポイントを登録するのだ
func (h *IdentificationHandler) RegisterIdentificationRoutes(router *gin.RouterGroup) {
	identifications := router.Group("/identifications")
	{
		identifications.GET("", h.GetAllIdentifications)
		identifications.GET("/:id", h.GetIdentificationByID)
		identifications.POST("", h.CreateIdentification)
		identifications.PUT("/:id", h.UpdateIdentification)
		identifications.DELETE("/:id", h.DeleteIdentification)
	}
}

func (h *IdentificationHandler) GetIdentificationByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "無効なIDです"})
		return
	}

	identification, err := h.identificationService.GetIdentificationByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "同定情報が見つかりません"})
		return
	}

	c.JSON(http.StatusOK, identification)
}

func (h *IdentificationHandler) GetAllIdentifications(c *gin.Context) {
	identifications, err := h.identificationService.GetAllIdentifications()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "同定情報の取得に失敗しました"})
		return
	}
	c.JSON(http.StatusOK, identifications)
}

func (h *IdentificationHandler) CreateIdentification(c *gin.Context) {
	var req service.CreateIdentificationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "リクエストの形式が正しくありません"})
		return
	}

	created, err := h.identificationService.CreateIdentification(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "同定情報の作成に失敗しました"})
		return
	}

	c.JSON(http.StatusCreated, created)
}

func (h *IdentificationHandler) UpdateIdentification(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "無効なIDです"})
		return
	}

	var req service.UpdateIdentificationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "リクエストの形式が正しくありません"})
		return
	}

	updated, err := h.identificationService.UpdateIdentification(uint(id), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "同定情報の更新に失敗しました"})
		return
	}

	c.JSON(http.StatusOK, updated)
}

func (h *IdentificationHandler) DeleteIdentification(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "無効なIDです"})
		return
	}

	if err := h.identificationService.DeleteIdentification(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "同定情報の削除に失敗しました"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "同定情報を削除しました"})
}
