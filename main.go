package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"
)

// Counter for generating unique IDs
var counter uint64
var mu sync.Mutex

// Employee struct represents an employee
type Employee struct {
	ID       uint64  `json:"id"`
	Name     string  `json:"name"`
	Position string  `json:"position"`
	Salary   float64 `json:"salary"`
}

// Employee struct represents an employee payload
type RequestEmployee struct {
	Name     string  `json:"name"`
	Position string  `json:"position"`
	Salary   float64 `json:"salary"`
}

var employees []Employee

func main() {
	// Define API endpoints
	http.HandleFunc("/employees", getAllEmployees)
	http.HandleFunc("/employee/create", createEmployee)
	http.HandleFunc("/employee/update/", updateEmployee)
	http.HandleFunc("/employee/delete", deleteEmployee)

	// Start the server
	fmt.Println("Server listening on port 8080...")
	log.Fatal(http.ListenAndServe(":9080", nil))
}

// getAllEmployees retrieves all employees
func getAllEmployees(w http.ResponseWriter, r *http.Request) {
	// Parse query parameters
	pageStr := r.URL.Query().Get("page")
	pageSizeStr := r.URL.Query().Get("pageSize")

	// Set default values for page and page size
	page := 1
	pageSize := 10

	// Convert query parameters to integers
	if pageStr != "" {
		p, err := strconv.Atoi(pageStr)
		if err == nil && p > 0 {
			page = p
		}
	}

	if pageSizeStr != "" {
		ps, err := strconv.Atoi(pageSizeStr)
		if err == nil && ps > 0 {
			pageSize = ps
		}
	}

	// Calculate the start and end indices for pagination
	startIndex := (page - 1) * pageSize
	endIndex := startIndex + pageSize
	if startIndex >= len(employees) {
		// Serialize and send the empty JSON response
		json.NewEncoder(w).Encode([]Employee{})
		return
	}
	if endIndex > len(employees) {
		endIndex = len(employees)
	}
	// Extract the subset of employees for the requested page
	pagedEmployees := employees[startIndex:endIndex]

	// Serialize and send the pagedEmployees as JSON response
	json.NewEncoder(w).Encode(pagedEmployees)
}

// createEmployee adds a new employee
func createEmployee(w http.ResponseWriter, r *http.Request) {
	var requestEmployee RequestEmployee
	err := json.NewDecoder(r.Body).Decode(&requestEmployee)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	var employee Employee
	employee.Name = requestEmployee.Name
	employee.Position = requestEmployee.Position
	employee.Salary = requestEmployee.Salary
	employee.ID = GenerateUniqueID()
	employees = append(employees, employee)
	json.NewEncoder(w).Encode(employee)
}

// updateEmployee updates an existing employee
func updateEmployee(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 4 {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	// Convert string to uint64
	empID, _ := strconv.ParseUint(parts[3], 10, 64)
	var updatedEmployee RequestEmployee
	err := json.NewDecoder(r.Body).Decode(&updatedEmployee)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	for i, emp := range employees {
		if emp.ID == empID {
			employees[i].Name = updatedEmployee.Name
			employees[i].Position = updatedEmployee.Position
			employees[i].Salary = updatedEmployee.Salary
			json.NewEncoder(w).Encode(employees[i])
			return
		}
	}

	http.Error(w, "Employee not found", http.StatusNotFound)
}

// deleteEmployee deletes an existing employee
func deleteEmployee(w http.ResponseWriter, r *http.Request) {
	var employeeID uint64
	err := json.NewDecoder(r.Body).Decode(&employeeID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	// Convert string to uint64
	//empID, _ := strconv.ParseUint(employeeID, 10, 64)
	for i, emp := range employees {
		if emp.ID == employeeID {
			employees = append(employees[:i], employees[i+1:]...)
			fmt.Fprintf(w, "Employee with ID %d deleted", employeeID)
			return
		}
	}

	http.Error(w, "Employee not found", http.StatusNotFound)
}

// GenerateUniqueID generates a unique ID based on timestamp and counter
func GenerateUniqueID() uint64 {
	mu.Lock()
	defer mu.Unlock()

	// Get current timestamp in milliseconds
	timestamp := uint64(time.Now().UnixNano() / int64(time.Millisecond))

	// Increment counter
	counter++

	// Combine timestamp and counter to generate unique ID
	uniqueID := (timestamp << 32) | (counter & 0xFFFFFFFF)

	return uniqueID
}
