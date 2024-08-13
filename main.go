package main

import (
	"errors"
	"fmt"
	"log"
	"os"
)

var (
	ErrFailedToReadFile = errors.New("failed to read file")
)

type Employee struct {
	Name string
}

func readEmployeeDataFromFile(filename string) ([]Employee, error) {
	_, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", ErrFailedToReadFile, filename)
	}

	return nil, nil
}

func getEmployeeData() ([]Employee, error) {
	employees, err := readEmployeeDataFromFile("employees.txt")
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve employee data: %w", err)
	}
	return employees, nil
}

func generateEmployeeReport() ([]Employee, error) {
	employees, err := getEmployeeData()
	if err != nil {
		return nil, fmt.Errorf("failed to generate employee report: %w", err)
	}
	return employees, nil
}

func main() {
	if _, err := generateEmployeeReport(); err != nil {
		log.Println("Error generating report:", err)
	}
}
