
// backend/internal/handler/project_handler.go
package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/saku-730/specimen-web/backend/internal/service"
)

type ProjectHandler struct {
	projectService service.ProjectService
}

func NewProjectHandler(projectService service.ProjectService) *ProjectHandler {
	return &ProjectHandler{projectService: projectService}
}

// RegisterProjectRoutes はルーターにプロジェクト関連のエンドポイントを登録するのだ
func (h *ProjectHandler) RegisterProjectRoutes(router *gin.RouterGroup) {
	projects := router.Group("/projects")
	{
		projects.GET("", h.GetAllProjects)
		projects.GET("/:id", h.GetProjectByID)
		projects.POST("", h.CreateProject)
		projects.PUT("/:id", h.UpdateProject)
		projects.DELETE("/:id", h.DeleteProject)
	}
}

func (h *ProjectHandler) GetProjectByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "無効なIDです"})
		return
	}

	project, err := h.projectService.GetProjectByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "プロジェクトが見つかりません"})
		return
	}

	c.JSON(http.StatusOK, project)
}

func (h *ProjectHandler) GetAllProjects(c *gin.Context) {
	projects, err := h.projectService.GetAllProjects()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "プロジェクトの取得に失敗しました"})
		return
	}
	c.JSON(http.StatusOK, projects)
}

func (h *ProjectHandler) CreateProject(c *gin.Context) {
	var req service.CreateProjectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "リクエストの形式が正しくありません"})
		return
	}

	createdProject, err := h.projectService.CreateProject(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "プロジェクトの作成に失敗しました"})
		return
	}

	c.JSON(http.StatusCreated, createdProject)
}

func (h *ProjectHandler) UpdateProject(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "無効なIDです"})
		return
	}

	var req service.UpdateProjectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "リクエストの形式が正しくありません"})
		return
	}

	updatedProject, err := h.projectService.UpdateProject(uint(id), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "プロジェクトの更新に失敗しました"})
		return
	}

	c.JSON(http.StatusOK, updatedProject)
}

func (h *ProjectHandler) DeleteProject(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "無効なIDです"})
		return
	}

	if err := h.projectService.DeleteProject(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "プロジェクトの削除に失敗しました"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "プロジェクトを削除しました"})
}
