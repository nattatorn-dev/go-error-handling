package main

import (
	"errors"
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
		log.Println(ErrFailedToReadFile, filename, err)
		return nil, err
	}

	return nil, nil
}

func getEmployeeData() ([]Employee, error) {
	employees, err := readEmployeeDataFromFile("employees.txt")
	if err != nil {
		log.Println("failed to retrieve employee data", err)
		return nil, err
	}
	return employees, nil
}

func generateEmployeeReport() ([]Employee, error) {
	employees, err := getEmployeeData()
	if err != nil {
		log.Println("failed to generate employee report", err)
		return nil, err
	}
	return employees, nil
}

func main() {
	if _, err := generateEmployeeReport(); err != nil {
		log.Println("Error generating report:", err)
	}
}
