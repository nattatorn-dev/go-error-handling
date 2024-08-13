package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
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

func ErrorHandlerMiddleware() gin.HandlerFunc {
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
				// If it's a CustomError, respond with detailed information
				c.JSON(http.StatusInternalServerError, gin.H{
					"error":   customErr.Message,
					"code":    customErr.Code,
					"details": customErr.Error(),
				})
			} else {
				// For generic errors, respond with a generic message
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
		return nil, fmt.Errorf("%w: failed to create employee report", err)
	}
	return employees, nil
}

func handleReportGeneration(c *gin.Context) {
	employees, err := createEmployeeReport()
	if err != nil {
		errMsg := "failed to create employee report"
		err := fmt.Errorf("%w", NewCustomError(CodeCreateEmployeeReportFailure, errMsg, err))
		log.Println(errMsg, err)
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, employees)
}

func main() {
	router := gin.Default()

	router.Use(ErrorHandlerMiddleware())

	router.GET("/report", handleReportGeneration)

	if err := router.Run("127.0.0.1:8080"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
