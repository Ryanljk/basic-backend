package middleware

import (
    "bytes"
    "io"
    "log"
    "github.com/gin-gonic/gin"
)

//logging function to track what form data is being sent by the client
func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		bodyBytes, err := io.ReadAll(c.Request.Body)
		if err != nil {
			log.Println("Failed to read request body:", err)
			c.AbortWithStatus(500)
			return
		}

		c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
		if (string(bodyBytes) != "") {
			log.Printf("%s %s \nBody: %s\n", c.Request.Method, c.Request.URL.Path, string(bodyBytes))
		}

		c.Next()
	}
}
