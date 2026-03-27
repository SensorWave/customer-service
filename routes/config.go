package routes

import (
	appconfig "customer-service/config"

	"github.com/gin-gonic/gin"
)

const configContextKey = "app_config"

func WithConfig(cfg *appconfig.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set(configContextKey, cfg)
		c.Next()
	}
}

func ConfigFromContext(c *gin.Context) (*appconfig.Config, bool) {
	value, exists := c.Get(configContextKey)
	if !exists {
		return nil, false
	}

	cfg, ok := value.(*appconfig.Config)
	if !ok {
		return nil, false
	}

	return cfg, true
}

func MustConfig(c *gin.Context) *appconfig.Config {
	return c.MustGet(configContextKey).(*appconfig.Config)
}
