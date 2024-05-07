# Employee CRUD REST API

This is a RESTful API for managing employee records. It allows users to perform CRUD operations on employee data using HTTP requests.

## Features

- **Create**: Add new employee records to the database.
- **Read**: Retrieve employee records from the database.
- **Update**: Modify existing employee records in the database.
- **Delete**: Remove employee records from the database.

## Technologies Used

- Go (Golang) programming language

## API Endpoints

The following endpoints are available:

- **POST /employees**: Create a new employee record.
- **GET /employees**: Retrieve a list of all employee records.
- **PUT /employees/{id}**: Update an existing employee record by ID.
- **DELETE /employees**: Delete an employee record by ID.

## Setup

1. **Clone the repository**: 
    ```
    git clone https://github.com/trulygowtham/employeCURD.git
    ```
    
2. **Run the application**:
    ```
    go run main.go
    ```

## Usage

- **Create Employee**: Send a POST request to `/employees` with JSON data containing employee details.
- **Retrieve Employees**: Send a GET request to `/employees` to retrieve a list of all employee records.
- **Update Employee**: Send a PUT request to `/employees/{id}` with JSON data containing updated employee details.
- **Delete Employee**: Send a DELETE request to `/employees` with JSON data containing to delete an employee record by ID.

## Contributing

Contributions are welcome! Please open an issue or submit a pull request with any improvements or bug fixes.

## License

This project is licensed under the [MIT License](LICENSE).
