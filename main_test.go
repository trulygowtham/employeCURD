package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
)

func TestCRUDOperations(t *testing.T) {
	// Initialize empty employee slice
	employees = []Employee{}

	// Create employee
	createEmployeeTest(t)

	// Get all employees
	getAllEmployeesTest(t)

	// Update employee
	updateEmployeeTest(t)

	// Delete employee
	deleteEmployeeTest(t)
}

func createEmployeeTest(t *testing.T) {
	// Create a sample employee
	newEmployee := RequestEmployee{
		Name:     "John",
		Position: "Lead",
		Salary:   3000,
	}

	// Convert employee to JSON
	reqBody, err := json.Marshal(newEmployee)
	if err != nil {
		t.Fatalf("Failed to marshal request body: %v", err)
	}

	// Create HTTP request
	req, err := http.NewRequest("POST", "/employees/create", bytes.NewBuffer(reqBody))
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	// Create HTTP test recorder
	rec := httptest.NewRecorder()

	// Serve HTTP request
	http.HandlerFunc(createEmployee).ServeHTTP(rec, req)

	// Check HTTP status code
	if rec.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, rec.Code)
	}

	// Decode response body
	var createdEmployee Employee
	if err := json.Unmarshal(rec.Body.Bytes(), &createdEmployee); err != nil {
		t.Fatalf("Failed to unmarshal response body: %v", err)
	}

	// Verify if the created employee matches the expected employee
	if createdEmployee.Name != newEmployee.Name ||
		createdEmployee.Position != newEmployee.Position ||
		createdEmployee.Salary != newEmployee.Salary {
		t.Errorf("Created employee does not match expected employee")
	}
}

func getAllEmployeesTest(t *testing.T) {
	// Initialize sample employees
	employees = []Employee{
		{ID: 1, Name: "John", Position: "Engineer", Salary: 30},
		{ID: 2, Name: "Jane", Position: "Testing", Salary: 35},
		{ID: 3, Name: "Michael", Position: "Manager", Salary: 40},
	}

	// Test case: Get first page with default page size
	testGetAllEmployeesWithPagination(t, 1, 3, 3)

	// Test case: Get second page with page size of 2
	testGetAllEmployeesWithPagination(t, 2, 1, 1)

	// Test case: Get third page with page size of 2
	testGetAllEmployeesWithPagination(t, 3, 1, 1)

	// Test case: Get page beyond total number of employees
	testGetAllEmployeesWithPagination(t, 3, 10, 0)

	// Test case: Invalid page number (negative)
	testGetAllEmployeesWithPagination(t, -1, 10, 0)

	// Test case: Invalid page size (zero)
	testGetAllEmployeesWithPagination(t, 1, 0, 0)
}

func testGetAllEmployeesWithPagination(t *testing.T, page, pageSize, expectedCount int) {
	// Create HTTP request with query parameters
	req, err := http.NewRequest("GET", "/employees?page="+strconv.Itoa(page)+"&pageSize="+strconv.Itoa(pageSize), nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	// Create HTTP test recorder
	rec := httptest.NewRecorder()

	// Serve HTTP request
	http.HandlerFunc(getAllEmployees).ServeHTTP(rec, req)

	// Check HTTP status code
	if rec.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, rec.Code)
	}

	// Decode response body
	var response []Employee
	if err := json.Unmarshal(rec.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to unmarshal response body: %v", err)
	}

	// Verify the number of employees returned
	if len(response) != expectedCount {
		t.Errorf("Expected %d employees, got %d", expectedCount, len(response))
	}

}

func updateEmployeeTest(t *testing.T) {
	// Create a sample updated employee
	updatedEmployee := Employee{
		ID:       1,
		Name:     "John",
		Position: "Lead",
		Salary:   3000,
	}

	// Convert updated employee to JSON
	reqBody, err := json.Marshal(updatedEmployee)
	if err != nil {
		t.Fatalf("Failed to marshal request body: %v", err)
	}

	// Create HTTP request
	req, err := http.NewRequest("PUT", "/employees/update/1", bytes.NewBuffer(reqBody))
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	// Create HTTP test recorder
	rec := httptest.NewRecorder()

	// Serve HTTP request
	http.HandlerFunc(updateEmployee).ServeHTTP(rec, req)

	// Check HTTP status code
	if rec.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, rec.Code)
	}

	// Decode response body
	var updatedEmployeeResponse Employee
	if err := json.Unmarshal(rec.Body.Bytes(), &updatedEmployeeResponse); err != nil {
		t.Fatalf("Failed to unmarshal response body: %v", err)
	}

	// Verify if the updated employee matches the expected updated employee
	if updatedEmployeeResponse.ID != updatedEmployee.ID ||
		updatedEmployeeResponse.Name != updatedEmployee.Name ||
		updatedEmployeeResponse.Position != updatedEmployee.Position ||
		updatedEmployeeResponse.Salary != updatedEmployee.Salary {
		t.Errorf("Updated employee does not match expected updated employee")
	}
}

func deleteEmployeeTest(t *testing.T) {
	// Create a sample employee ID to delete
	employeeID := 1

	// Convert employee ID to JSON
	reqBody, err := json.Marshal(employeeID)
	if err != nil {
		t.Fatalf("Failed to marshal request body: %v", err)
	}

	// Create HTTP request
	req, err := http.NewRequest("DELETE", "/employees/delete", bytes.NewBuffer(reqBody))
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	// Create HTTP test recorder
	rec := httptest.NewRecorder()

	// Serve HTTP request
	http.HandlerFunc(deleteEmployee).ServeHTTP(rec, req)

	// Check HTTP status code
	if rec.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, rec.Code)
	}

	// Verify if the employee was deleted by checking the number of employees
	if len(employees) != 2 {
		t.Errorf("Expected 2 employees after deletion, got %d", len(employees))
	}
}
