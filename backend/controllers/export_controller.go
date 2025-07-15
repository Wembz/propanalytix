package controllers

import (
	"net/http"
	"strconv"

	"github.com/Wembz/propanalytix/backend/services"
	"github.com/gin-gonic/gin"
	
)

type ExportController struct {
	ExportService *services.ExportService
}

func NewExportController() *ExportController {
	return &ExportController{
		ExportService: services.NewExportService(),
	}
}

// GET /export/pdf/:id
func (ec *ExportController) ExportPDF(c *gin.Context) {
	userID := c.MustGet("userID").(uint)

	idStr := c.Param("id")
	calcID64, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid calculation ID"})
		return
	}

	calculationID := uint(calcID64)

	file, err := ec.ExportService.ExportToPDF(userID, calculationID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Export failed"})
		return
	}
	c.Header("Content-Disposition", "attachment; filename=report.pdf")
	c.File( file)
}

// GET /export/excel/:id
func (ec *ExportController) ExportExcel(c *gin.Context) {
	userID := c.MustGet("userID").(uint)
	idStr := c.Param("id")
	calcID64, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid calculation ID"})
		return
	}
	calculationID := uint(calcID64)

	file, err := ec.ExportService.ExportToExcel(userID, calculationID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Export failed"})
		return
	}
	c.Header("Content-Disposition", "attachment; filename=report.xlsx")
	c.File( file)
}