package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Wembz/propanalytix/backend/services"
	"github.com/gin-gonic/gin"
)

type CalculationController struct {
	CalculationService *services.CalculationService
}

func NewCalculationController() *CalculationController {
	return &CalculationController{
		CalculationService: services.NewCalculationService(),
	}
}

// POST /calculation
func (cc *CalculationController) Submit(c *gin.Context) {
	userID := c.MustGet("userID").(uint)
	input := c.PostForm("input")
	result := c.PostForm("result")
	templateIDStr := c.PostForm("template_id")
	exportType := c.PostForm("export")

	var templateID *uint
	if templateIDStr != "" {
		tid, err := parseUint(templateIDStr)
		if err == nil {
			templateID = &tid
		}
		
	}

	err  := cc.CalculationService.SaveCalculation(userID, input, result, templateID, exportType)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save calculation"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Calculation saved"})
}

// GET /calculation/recent
func (cc *CalculationController) Recent(c *gin.Context) {
	userID := c.MustGet("userID").(uint)
	items, err := cc.CalculationService.GetRecentCalculations(userID, 5)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve calculations"})
		return
	}
	c.JSON(http.StatusOK, items)
}

// DELETE /calculation/:id
 func (cc *CalculationController) Delete(c *gin.Context) {
	userID := c.MustGet("userID").(uint)
	id, err := parseUint(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}
	err = cc.CalculationService.DeleteCalculation(userID, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Delete failed"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Calculation deleted"})
 }
 func parseUint(s string) (uint, error) {
	num, err := json.Number(s).Int64()
	if err != nil {
		return 0, err
	} 
	if num < 0 {
		return 0, fmt.Errorf("negative value cannot be parsed as uint")
	}
	return uint(num), nil
 }
