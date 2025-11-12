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
    r.POST("/companies", h.CreateCompany)
}

func (h *Handler) CreateCompany(c *gin.Context) {
    var company Company
    if err := c.ShouldBindJSON(&company); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    company.ID = uuid.NewString()
    if err := h.Service.CreateCompany(&company); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusCreated, company)
}
