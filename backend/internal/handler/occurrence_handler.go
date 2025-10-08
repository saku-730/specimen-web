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

	// フォーム用の選択肢などを取得するエンドポイント
	router.GET("/languages", h.GetAllLanguages)

	// フォーム全体を一度に登録するエンドポイント
	router.POST("/full-occurrence", h.CreateFullOccurrence)

	// search occurrence data with some others
	router.GET("/search", h.Search)
}

// CreateFullOccurrence はフォームからの全入力をまとめて登録するハンドラなのだ
func (h *OccurrenceHandler) CreateFullOccurrence(c *gin.Context) {
	var req service.FullOccurrenceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "リクエストの形式が正しくありません: " + err.Error()})
		return
	}

	if err := h.occurrenceService.CreateFullOccurrence(req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "データの登録に失敗しました: " + err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "データが正常に登録されました"})
}

func (h *OccurrenceHandler) GetAllLanguages(c *gin.Context) {
	languages, err := h.occurrenceService.GetAllLanguages()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Falied to get language list"})
		return
	}
	c.JSON(http.StatusOK, languages)
}

func (h *OccurrenceHandler) Search(c *gin.Context) {
	var req service.SearchRequest //for c.SholdBindQuery

	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error":"invalid search parameters"})
		return
	}

	results, err := h.occurrenceService.Search(req)
	if err != nil{
		c.JSON(http.StatusInternalServarError, gin.H{"error":"Failed search"})

	c.JSON(http.StatusOK,results)
}
