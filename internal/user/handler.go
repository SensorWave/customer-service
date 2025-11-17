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
    users := r.Group("/companies/:companyID/users")
    {
        users.POST("/", h.CreateUser)
        users.GET("/", h.GetUsersByCompany)
        users.GET("/:id", h.GetUserByID)
        users.PUT("/:id", h.UpdateUser)
        users.DELETE("/:id", h.DeleteUser)
    }
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

    if err := h.Service.CreateUser(c.Request.Context(), &user); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusCreated, user)
}

func (h *Handler) GetUsersByCompany(c *gin.Context) {
    companyID := c.Param("companyID")
    users, err := h.Service.GetUsersByCompany(c.Request.Context(), companyID)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, users)
}

func (h *Handler) GetUserByID(c *gin.Context) {
    id := c.Param("id")

    user, err := h.Service.GetUserByID(c.Request.Context(), id)
    if err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
        return
    }

    c.JSON(http.StatusOK, user)
}

func (h *Handler) UpdateUser(c *gin.Context) {
    id := c.Param("id")
    companyID := c.Param("companyID")

    var input User
    if err := c.ShouldBindJSON(&input); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    input.ID = id
    input.CompanyID = companyID

    if err := h.Service.UpdateUser(c.Request.Context(), &input); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, input)
}

func (h *Handler) DeleteUser(c *gin.Context) {
    id := c.Param("id")

    if err := h.Service.DeleteUser(c.Request.Context(), id); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "user deleted"})
}
