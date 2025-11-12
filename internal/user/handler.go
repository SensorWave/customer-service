package user

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
    r.POST("/companies/:companyID/users", h.CreateUser)
}

func (h *Handler) CreateUser(c *gin.Context) {
    companyID := c.Param("companyID")
    var user User
    if err := c.ShouldBindJSON(&user); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    user.ID = uuid.NewString()
    user.CompanyID = companyID

    if err := h.Service.CreateUser(&user); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusCreated, user)
}
