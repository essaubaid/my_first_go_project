package middleware

import (
	"net/http"

	"github.com/essaubaid/my_first_go_project/helpers"
	"github.com/gin-gonic/gin"
)

func Authentication() gin.HandlerFunc {
	return func(c *gin.Context) {
		clientToken := c.Request.Header.Get("token")
		if clientToken == "" {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "No Authorization header provided",
			})
			c.Abort()
			return
		}

		claim, err := helpers.ValidateToken(clientToken)

		if err != "" {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err,
			})
			c.Abort()
			return
		}

		c.Set("email", claim.Email)
		c.Set("first_name", claim.First_name)
		c.Set("last_name", claim.Last_name)
		c.Set("uid", claim.Uid)

		c.Next()
	}
}
