package main

import (
	"errors"
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

const (
	CodeCreateEmployeeReportFailure = 1001
)

type CustomError struct {
	Code    int
	Message string
	Err     error
}

func (e *CustomError) Error() string {
	return fmt.Sprintf("Error %d: %s, %v", e.Code, e.Message, e.Err)
}

func (e *CustomError) Unwrap() error {
	return e.Err
}

func NewCustomError(code int, message string, err error) *CustomError {
	return &CustomError{
		Code:    code,
		Message: message,
		Err:     err,
	}
}

func ErrorHandlerMiddleware(logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Process the request
		c.Next()

		// Check if errors occurred during the request
		if len(c.Errors) > 0 {
			// Retrieve the last error
			err := c.Errors.Last().Err

			// Check if it's a CustomError
			var customErr *CustomError
			if errors.As(err, &customErr) {
				c.JSON(http.StatusInternalServerError, gin.H{
					"error":   customErr.Message,
					"code":    customErr.Code,
					"details": customErr.Error(),
				})
			} else {
				// For generic errors, respond with a generic message
				logger.Error("An unexpected error occurred", zap.Error(err))
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": "An unexpected error occurred",
				})
			}
		}
	}
}

type Employee struct {
	Name string `json:"name"`
}

func loadEmployeeDataFromFile(filename string) ([]Employee, error) {
	_, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("%w: could not read file %s", err, filename)
	}

	return nil, nil
}

func fetchEmployeeData() ([]Employee, error) {
	employees, err := loadEmployeeDataFromFile("employees.txt")
	if err != nil {
		return nil, fmt.Errorf("%w: failed to load employee data", err)
	}

	return employees, nil
}

func createEmployeeReport() ([]Employee, error) {
	employees, err := fetchEmployeeData()
	if err != nil {
		return nil, fmt.Errorf("%w: failed to fetch employee data", err)
	}

	return employees, nil
}

func handleReportGeneration(logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		employees, err := createEmployeeReport()
		if err != nil {
			errMsg := "failed to create employee report"
			err := fmt.Errorf("%w", NewCustomError(CodeCreateEmployeeReportFailure, errMsg, err))
			logger.Error(errMsg, zap.Error(err))
			c.Error(err)
			return
		}

		logger.Info("Report generated successfully", zap.Int("employee_count", len(employees)))
		c.JSON(http.StatusOK, employees)
	}
}

func main() {
	// Initialize zap logger
	logger, err := zap.NewProduction()
	if err != nil {
		fmt.Printf("Failed to initialize zap logger: %v\n", err)
		return
	}
	defer logger.Sync()

	router := gin.Default()

	// Use the ErrorHandlerMiddleware with zap logger
	router.Use(ErrorHandlerMiddleware(logger))

	// Register the route with zap logger
	router.GET("/report", handleReportGeneration(logger))

	// Start the server
	if err := router.Run("127.0.0.1:8080"); err != nil {
		logger.Fatal("Failed to start server", zap.Error(err))
	}
}
