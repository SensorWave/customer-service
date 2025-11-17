package company

import (
    "github.com/gin-gonic/gin"
    "github.com/google/uuid"
    "net/http"
)

type Handler struct {
    Service *Service
}

func NewHandler(s *Service) *Handler {
    return &Handler{Service: s}
}

func (h *Handler) RegisterRoutes(r *gin.Engine) {
    companies := r.Group("/companies")
    {
        companies.POST("/", h.CreateCompany)
        companies.GET("/", h.GetAllCompanies)
        companies.GET("/:id", h.GetCompanyByID)
        companies.PUT("/:id", h.UpdateCompany)
        companies.DELETE("/:id", h.DeleteCompany)
    }
}

func (h *Handler) CreateCompany(c *gin.Context) {
    var company Company
    if err := c.ShouldBindJSON(&company); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    // Génère un ID UUID v4
    company.ID = uuid.NewString()

    // Passe le context HTTP -> service
    if err := h.Service.CreateCompany(c.Request.Context(), &company); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusCreated, company)
}

func (h *Handler) GetCompanyByID(c *gin.Context) {
    id := c.Param("id")

    company, err := h.Service.GetCompanyByID(c.Request.Context(), id)
    if err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "company not found"})
        return
    }

    c.JSON(http.StatusOK, company)
}

func (h *Handler) GetAllCompanies(c *gin.Context) {
    companies, err := h.Service.GetAllCompanies(c.Request.Context())
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, companies)
}

func (h *Handler) UpdateCompany(c *gin.Context) {
    id := c.Param("id")

    var input Company
    if err := c.ShouldBindJSON(&input); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    input.ID = id

    if err := h.Service.UpdateCompany(c.Request.Context(), &input); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, input)
}

func (h *Handler) DeleteCompany(c *gin.Context) {
    id := c.Param("id")

    if err := h.Service.DeleteCompany(c.Request.Context(), id); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "company deleted"})
}
