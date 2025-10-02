// backend/internal/handler/observation_handler.go
package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/saku-730/specimen-web/backend/internal/service"
)

type ObservationHandler struct {
	observationService service.ObservationService
}

func NewObservationHandler(observationService service.ObservationService) *ObservationHandler {
	return &ObservationHandler{observationService: observationService}
}

// RegisterObservationRoutes はルーターに観察情報関連のエンドポイントを登録するのだ
func (h *ObservationHandler) RegisterObservationRoutes(router *gin.RouterGroup) {
	observations := router.Group("/observations")
	{
		observations.GET("", h.GetAllObservations)
		observations.GET("/:id", h.GetObservationByID)
		observations.POST("", h.CreateObservation)
		observations.PUT("/:id", h.UpdateObservation)
		observations.DELETE("/:id", h.DeleteObservation)
	}
	router.GET("/observation-methods", h.GetAllObservationMethods)
}

func (h *ObservationHandler) GetObservationByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "無効なIDです"})
		return
	}

	observation, err := h.observationService.GetObservationByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "観察情報が見つかりません"})
		return
	}

	c.JSON(http.StatusOK, observation)
}

func (h *ObservationHandler) GetAllObservations(c *gin.Context) {
	observations, err := h.observationService.GetAllObservations()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "観察情報の取得に失敗しました"})
		return
	}
	c.JSON(http.StatusOK, observations)
}

func (h *ObservationHandler) CreateObservation(c *gin.Context) {
	var req service.CreateObservationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "リクエストの形式が正しくありません"})
		return
	}

	created, err := h.observationService.CreateObservation(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "観察情報の作成に失敗しました"})
		return
	}

	c.JSON(http.StatusCreated, created)
}

func (h *ObservationHandler) UpdateObservation(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "無効なIDです"})
		return
	}

	var req service.UpdateObservationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "リクエストの形式が正しくありません"})
		return
	}

	updated, err := h.observationService.UpdateObservation(uint(id), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "観察情報の更新に失敗しました"})
		return
	}

	c.JSON(http.StatusOK, updated)
}

func (h *ObservationHandler) DeleteObservation(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "無効なIDです"})
		return
	}

	if err := h.observationService.DeleteObservation(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "観察情報の削除に失敗しました"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "観察情報を削除しました"})
}

// GetAllObservationMethods は全ての観察方法を取得するのだ
func (h *ObservationHandler) GetAllObservationMethods(c *gin.Context) {
	methods, err := h.observationService.GetAllObservationMethods()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "観察方法リストの取得に失敗しました"})
		return
	}
	c.JSON(http.StatusOK, methods)
}