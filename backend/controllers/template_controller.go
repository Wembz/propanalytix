package controllers

import (
	
	"net/http"
	"strconv"

	"github.com/Wembz/propanalytix/backend/services"
	"github.com/gin-gonic/gin"
	
)

type TemplateController struct {
	TemplateService *services.TemplateService
}

func NewTemplateController() *TemplateController {
	return &TemplateController{
		TemplateService: services.NewTemplateService(),
	}
}

// POST /template
func (tc *TemplateController) Save(c *gin.Context) {
	userID := c.MustGet("userID").(uint)
	var req struct {
		Name string `json:"name"`
		Data string `json:"data"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	err := tc.TemplateService.Save(userID, req.Name, req.Data)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save template"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Template saved"})
}

// GET /template
func (tc *TemplateController) List(c *gin.Context) {
	userID := c.MustGet("userID").(uint)
	id, _ := strconv.Atoi(c.Param("id"))
	template, err := tc.TemplateService.GetByID(userID, uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Template not found"})
		return
	}
	c.JSON(http.StatusOK, template)
}

// GET /template/:id
func (tc *TemplateController) Get(c *gin.Context) {
	userID := c.MustGet("userID").(uint)
	id, _ := strconv.Atoi(c.Param("id"))
	template, err := tc.TemplateService.GetByID(userID, uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Template not found"})
	}
	c.JSON(http.StatusOK, template)
}

// PUT /template/:id
func (tc *TemplateController) Update(c *gin.Context) {
	userID := c.MustGet("userID").(uint)
	id, _ := strconv.Atoi(c.Param("id"))
	var req struct {
		Name string `json:"name"`
		Data string `json:"data"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}
	err := tc.TemplateService.Update(userID, uint(id), req.Name, req.Data)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Update failed"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Template updated"})
}

// DELETE /template/:id
func (tc *TemplateController) Delete(c *gin.Context) {
	userID := c.MustGet("userID").(uint)
	id, _ := strconv.Atoi(c.Param("id"))
	err := tc.TemplateService.Delete(userID, uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Delete failed"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Template deleted"})
}