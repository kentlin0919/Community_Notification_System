package middlewares

import (
	"bytes"
	"io/ioutil"
	"strings"
	"time"

	"Community_Notification_System/app/repositories/action_log"
	"Community_Notification_System/database/ActionLog_DB"
	"Community_Notification_System/utils"

	"github.com/gin-gonic/gin"
)

func ActionLogMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		startTime := time.Now()

		// Read the body
		var bodyBytes []byte
		if c.Request.Body != nil {
			bodyBytes, _ = ioutil.ReadAll(c.Request.Body)
		}
		// Restore the io.ReadCloser to its original state
		c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))

		// Process request
		c.Next()

		// After request
		// endTime := time.Now()
		// latency := endTime.Sub(startTime)

		// Get user ID from JWT token
		userID := ""
		claims, exists := c.Get("claims")
		if exists {
			if claimsMap, ok := claims.(*utils.Claims); ok {
				userID = claimsMap.UserID
			}
		}

		// Get module from path
		path := c.Request.URL.Path
		pathSegments := strings.Split(strings.Trim(path, "/"), "/")
		module := ""
		if len(pathSegments) > 2 {
			module = pathSegments[2] // e.g., /api/v1/user -> user
		}

		logEntry := &ActionLog_DB.ActionLog{
			Timestamp: startTime,
			Module:    module,
			APIPath:   c.Request.URL.Path,
			UserID:    userID,
		}

		repo := action_log.ActionLogRepository{}
		err := repo.CreateLog(logEntry)
		if err != nil {
			// Handle error, maybe log to a file
		}
	}
}
