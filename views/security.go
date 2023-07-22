package views

import (
	"github.com/gin-gonic/gin"
)

// Adds several headers to the response to improve security.
// For example, headers prevent embedding a site and passing information about the referrer.
func SecurityHeadersMiddleware(c *gin.Context) {
	c.Writer.Header().Set("X-Frame-Options", "SAMEORIGIN")
	c.Writer.Header().Set("X-Content-Type-Options", "nosniff")
	c.Writer.Header().Set("Referrer-Policy", "same-origin")
	c.Writer.Header().Set("Content-Security-Policy", "script-src 'self'; ")
}
