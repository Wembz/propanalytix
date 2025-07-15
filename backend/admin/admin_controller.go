package admin

import (
	"net/http"

	
	"github.com/Wembz/propanalytix/backend/services"
	"github.com/gin-gonic/gin"
	
)

func GetAllUsers(c *gin.Context) {
	users, err := services.FetchAllUsers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch users"})
		return
	}
	c.JSON(http.StatusOK, users)
}

func GetRevenueReport(c *gin.Context) {
	report, err := services.GenerateRevenueReport()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate report"})
		return
	}
	c.JSON(http.StatusOK, report)
}

func GetAuditLogs(c *gin.Context) {
	logs, err  := services.FetchAuditLogs()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch audit logs"})
		return
	}
	c.JSON(http.StatusOK, logs)
}

func GetCalculationAnalytics(c *gin.Context) {
	analytics, err := services.GenerateCalculationAnalytics()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load calculation analytics"})
		return
	}
	c.JSON(http.StatusOK, analytics)
}

func BanUser(c *gin.Context) {
	userID := c.Param("id")
	err := services.BanUserByID(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to ban user"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "User banned successfully"})	
}

func PromoteUser(c *gin.Context) {
	userID := c.Param("id")
	role := c.Query("role")
	if role == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing role query param"})
		return
	}
	
	err := services.UpdateUserRole(userID, role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user role"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "User role  updated"})

}
func DeleteUser(c *gin.Context) {
	userID := c.Param("id")
	err  := services.DeleteUserByID(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete user"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}

func GetPlanStats(c *gin.Context) {
	stats, err := services.GeneratePlanStatistics()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get plan stats"})
		return
	}
	c.JSON(http.StatusOK, stats)
}

func ExportUsersCSV(c *gin.Context) {
	csvData, err := services.ExportUserAsCSV()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to export user"})
		return
	}
	c.Header("Content-Disposition", "attachment; filename-user.csv")
	c.Data(http.StatusOK, "text/csv", csvData)
}
