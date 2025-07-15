package routes

import (
	"github.com/Wembz/propanalytix/backend/controllers"
	"github.com/Wembz/propanalytix/backend/middleware"
	webhook "github.com/Wembz/propanalytix/backend/webhooks"
	"github.com/gin-gonic/gin"
	"github.com/mehanizm/airtable"
	"gopkg.in/gomail.v2"
)

func SetupRouter(airtableClient *airtable.Client, emailClient *gomail.Dialer) *gin.Engine {
	r := gin.Default()

	// Public routes
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "OK"})
	})

	// Auth routes go here...
	r.POST("/login", controllers.Login)
    r.POST("/register", controllers.Register)

	// Stripe Webhooks (public)
	webhookController := controllers.NewWebhookController()
    r.POST("/webhooks/stripe", webhookController.HandleStripeWebhook)

	// Protected routes
	auth := r.Group("/")
	auth.Use(middleware.AuthMiddleware())

	// Admin-only routes
	admin := auth.Group("/admin")
	admin.Use(middleware.AdminOnly())
	admin.GET("/dashboard", controllers.AdminDashboard)
	admin.GET("/users", controllers.GetAllUsers)
	admin.POST("/plans/upgrade", controllers.AdminUpgradeUserPlan)

	// User profile routes
	user := controllers.NewUserController()
	auth.GET("/user/profile", user.GetProfile)
	auth.PUT("/user/profile", user.UpdateProfile)
	auth.DELETE("/user/deactivate", user.DeactivateAccount)
	auth.PUT("/user/plan", user.SetPlan)

	// Calculation routes
	calc := controllers.NewCalculationController()
	auth.POST("/calculation", calc.Submit)
	auth.GET("/calculation/recent", calc.Recent)
	auth.DELETE("/calculation/:id", calc.Delete)
		
	// Template routes
	template := controllers.NewTemplateController()
	auth.POST("/template", template.Save)
	auth.GET("/template", template.List)
	auth.GET("/template/:id", template.Get)
	auth.PUT("/template/:id", template.Update)
	auth.DELETE("/template/:id", template.Delete)

	// Subscription routes
	sub := controllers.NewSubscriptionController()
	auth.PUT("/subscription/upgrade", sub.UpgradeMembership)
	auth.GET("/subscription/status", sub.GetPlanStatus)

	// Export routes
	export := controllers.NewExportController()
	auth.GET("/export/pdf/:id", export.ExportPDF)
	auth.GET("/export/excel/:id", export.ExportExcel)

	// Notification routes
	notif := controllers.NewNotificationController()
	auth.POST("/notifications/summary", notif.SendSummaryEmail)

	// Internal cron route (can be protected via secret header/token later)
	r.POST("/notifications/expiry-alerts", notif.NotifyExpiringSubscriptions)

	return r

}