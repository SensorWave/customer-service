package user

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	Service *Service
}

func NewHandler(s *Service) *Handler {
	return &Handler{Service: s}
}

// --------- CREATE ----------
func (h *Handler) CreateUser(c *gin.Context) {
	companyIdStr := c.Param("id")
	companyId, err := strconv.Atoi(companyIdStr)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid company id"})
		return
	}

	var user User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user.CompanyID = companyId

	if err := h.Service.CreateUser(c.Request.Context(), &user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, user)
}

// --------- GET BY COMPANY ----------
func (h *Handler) GetUsersByCompany(c *gin.Context) {
	companyIdStr := c.Param("id")
	companyId, err := strconv.Atoi(companyIdStr)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid companyId"})
		return
	}

	users, err := h.Service.GetUsersByCompany(c.Request.Context(), companyId)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, users)
}

// --------- GET BY ID ----------
func (h *Handler) GetUserByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		return
	}

	user, err := h.Service.GetUserByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	c.JSON(http.StatusOK, user)
}

// --------- GET BY KEYCLOAK ID ----------
func (h *Handler) GetUserByKeycloakID(c *gin.Context) {
	keycloakID := c.Param("keycloak_id")
	if keycloakID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid keycloak id"})
		return
	}

	user, err := h.Service.GetUserByKeycloakID(c.Request.Context(), keycloakID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	c.JSON(http.StatusOK, user)
}

// --------- GET ROLE BY KEYCLOAK ID ----------
func (h *Handler) GetUserRoleByKeycloakID(c *gin.Context) {
	keycloakID := c.Param("keycloak_id")
	if keycloakID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid keycloak id"})
		return
	}

	role, err := h.Service.GetUserRoleByKeycloakID(c.Request.Context(), keycloakID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"role": role})
}

// --------- GET COMPANY ID BY KEYCLOAK ID ----------
func (h *Handler) GetUserCompanyIDByKeycloakID(c *gin.Context) {
	keycloakID := c.Param("keycloak_id")
	if keycloakID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid keycloak id"})
		return
	}

	companyID, err := h.Service.GetUserCompanyIDByKeycloakID(c.Request.Context(), keycloakID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"company_id": companyID})
}

// --------- UPDATE ----------
func (h *Handler) UpdateUser(c *gin.Context) {
	idStr := c.Param("id")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		return
	}

	var input User
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	input.ID = id

	if err := h.Service.UpdateUser(c.Request.Context(), &input); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, input)
}

// --------- DELETE ----------
func (h *Handler) DeleteUser(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		return
	}

	if err := h.Service.DeleteUser(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "user deleted"})
}
