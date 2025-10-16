package middleware

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// ResponseInterceptor middleware untuk memformat response secara konsisten
func ResponseInterceptor(logger *logrus.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Start timer
		start := time.Now()
		
		// Create custom response writer
		writer := &ResponseWriter{
			ResponseWriter: c.Writer,
			statusCode:     http.StatusOK,
		}
		c.Writer = writer
		
		// Process request
		c.Next()
		
		// Calculate duration
		duration := time.Since(start)
		
		// Log request
		logger.WithFields(logrus.Fields{
			"method":     c.Request.Method,
			"path":       c.Request.URL.Path,
			"status":     writer.statusCode,
			"duration":   duration,
			"client_ip":  c.ClientIP(),
			"user_agent": c.Request.UserAgent(),
		}).Info("HTTP Request")
		
		// Format response based on status code
		if writer.statusCode >= 400 {
			formatErrorResponse(c, writer.statusCode)
		}
	}
}

// ResponseWriter custom response writer untuk tracking status code
type ResponseWriter struct {
	gin.ResponseWriter
	statusCode int
}

// WriteHeader captures the status code
func (rw *ResponseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

// formatErrorResponse memformat error response secara konsisten
func formatErrorResponse(c *gin.Context, statusCode int) {
	// Skip jika response sudah diformat
	if c.Writer.Header().Get("Content-Type") != "" {
		return
	}
	
	var response gin.H
	
	switch statusCode {
	case http.StatusBadRequest:
		response = gin.H{
			"success": false,
			"error": gin.H{
				"code":    "BAD_REQUEST",
				"message": "Invalid request parameters",
				"details": "The request contains invalid or missing parameters",
			},
			"timestamp": time.Now().UTC().Format(time.RFC3339),
		}
	case http.StatusUnauthorized:
		response = gin.H{
			"success": false,
			"error": gin.H{
				"code":    "UNAUTHORIZED",
				"message": "Authentication required",
				"details": "Valid authentication credentials are required",
			},
			"timestamp": time.Now().UTC().Format(time.RFC3339),
		}
	case http.StatusForbidden:
		response = gin.H{
			"success": false,
			"error": gin.H{
				"code":    "FORBIDDEN",
				"message": "Access denied",
				"details": "You don't have permission to access this resource",
			},
			"timestamp": time.Now().UTC().Format(time.RFC3339),
		}
	case http.StatusNotFound:
		response = gin.H{
			"success": false,
			"error": gin.H{
				"code":    "NOT_FOUND",
				"message": "Resource not found",
				"details": "The requested resource was not found",
			},
			"timestamp": time.Now().UTC().Format(time.RFC3339),
		}
	case http.StatusMethodNotAllowed:
		response = gin.H{
			"success": false,
			"error": gin.H{
				"code":    "METHOD_NOT_ALLOWED",
				"message": "Method not allowed",
				"details": "The HTTP method is not allowed for this endpoint",
			},
			"timestamp": time.Now().UTC().Format(time.RFC3339),
		}
	case http.StatusInternalServerError:
		response = gin.H{
			"success": false,
			"error": gin.H{
				"code":    "INTERNAL_SERVER_ERROR",
				"message": "Internal server error",
				"details": "An unexpected error occurred on the server",
			},
			"timestamp": time.Now().UTC().Format(time.RFC3339),
		}
	case http.StatusServiceUnavailable:
		response = gin.H{
			"success": false,
			"error": gin.H{
				"code":    "SERVICE_UNAVAILABLE",
				"message": "Service temporarily unavailable",
				"details": "The service is temporarily unavailable, please try again later",
			},
			"timestamp": time.Now().UTC().Format(time.RFC3339),
		}
	default:
		response = gin.H{
			"success": false,
			"error": gin.H{
				"code":    "UNKNOWN_ERROR",
				"message": "An unknown error occurred",
				"details": "An unexpected error occurred",
			},
			"timestamp": time.Now().UTC().Format(time.RFC3339),
		}
	}
	
	c.JSON(statusCode, response)
}

// SuccessResponse memformat success response secara konsisten
func SuccessResponse(c *gin.Context, data interface{}, message string) {
	response := gin.H{
		"success":   true,
		"message":   message,
		"data":      data,
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	}
	
	c.JSON(http.StatusOK, response)
}

// ErrorResponse memformat error response secara konsisten
func ErrorResponse(c *gin.Context, statusCode int, code, message, details string) {
	response := gin.H{
		"success": false,
		"error": gin.H{
			"code":    code,
			"message": message,
			"details": details,
		},
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	}
	
	c.JSON(statusCode, response)
}

// ValidationErrorResponse memformat validation error response
func ValidationErrorResponse(c *gin.Context, errors map[string]string) {
	response := gin.H{
		"success": false,
		"error": gin.H{
			"code":    "VALIDATION_ERROR",
			"message": "Validation failed",
			"details": "One or more fields failed validation",
			"fields":  errors,
		},
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	}
	
	c.JSON(http.StatusBadRequest, response)
}
