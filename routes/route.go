package routes

import (
	company "customer-service/internal/company"
	user "customer-service/internal/user"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine, companyHandler *company.Handler, userHandler *user.Handler) {
	// Healthcheck
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// ----------------- COMPANY ROUTES -----------------
	companyGroup := r.Group("/companies")
	{
		companyGroup.POST("", companyHandler.CreateCompany)
		companyGroup.GET("", companyHandler.GetAllCompanies)
		companyGroup.GET("/:id", companyHandler.GetCompanyByID)
		companyGroup.PUT("/:id", companyHandler.UpdateCompany)
		companyGroup.DELETE("/:id", companyHandler.DeleteCompany)

		// ----------------- NESTED USER ROUTES -----------------
		companyGroup.POST("/:id/users-company", userHandler.CreateUser)
		companyGroup.GET("/:id/users-company", userHandler.GetUsersByCompany)
	}

	// ----------------- FLAT USER ROUTES -----------------
	r.GET("/users-company/:id", userHandler.GetUserByID)
	r.PUT("/users-company/:id", userHandler.UpdateUser)
	r.DELETE("/users-company/:id", userHandler.DeleteUser)
}
