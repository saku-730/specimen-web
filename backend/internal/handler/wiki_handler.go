
// backend/internal/handler/wiki_handler.go
package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/saku-730/specimen-web/backend/internal/service"
)

type WikiHandler struct {
	wikiService service.WikiService
}

func NewWikiHandler(wikiService service.WikiService) *WikiHandler {
	return &WikiHandler{wikiService: wikiService}
}

// RegisterWikiRoutes はルーターにWiki関連のエンドポイントを登録するのだ
func (h *WikiHandler) RegisterWikiRoutes(router *gin.RouterGroup) {
	wikis := router.Group("/wikis")
	{
		wikis.GET("", h.GetAllWikiPages)
		wikis.GET("/:id", h.GetWikiPageByID)
		wikis.POST("", h.CreateWikiPage)
		wikis.PUT("/:id", h.UpdateWikiPage)
		wikis.DELETE("/:id", h.DeleteWikiPage)
	}
}

func (h *WikiHandler) GetWikiPageByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "無効なIDです"})
		return
	}

	page, err := h.wikiService.GetWikiPageByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Wikiページが見つかりません"})
		return
	}

	c.JSON(http.StatusOK, page)
}

func (h *WikiHandler) GetAllWikiPages(c *gin.Context) {
	pages, err := h.wikiService.GetAllWikiPages()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Wikiページの取得に失敗しました"})
		return
	}
	c.JSON(http.StatusOK, pages)
}

func (h *WikiHandler) CreateWikiPage(c *gin.Context) {
	var req service.CreateWikiPageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "リクエストの形式が正しくありません"})
		return
	}

	created, err := h.wikiService.CreateWikiPage(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Wikiページの作成に失敗しました"})
		return
	}

	c.JSON(http.StatusCreated, created)
}

func (h *WikiHandler) UpdateWikiPage(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "無効なIDです"})
		return
	}

	var req service.UpdateWikiPageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "リクエストの形式が正しくありません"})
		return
	}

	updated, err := h.wikiService.UpdateWikiPage(uint(id), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Wikiページの更新に失敗しました"})
		return
	}

	c.JSON(http.StatusOK, updated)
}

func (h *WikiHandler) DeleteWikiPage(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "無効なIDです"})
		return
	}

	if err := h.wikiService.DeleteWikiPage(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Wikiページの削除に失敗しました"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Wikiページを削除しました"})
}
