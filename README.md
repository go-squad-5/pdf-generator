# Go Quiz Report PDF Generator & Email Service

A robust backend service built in Go that generates detailed PDF reports from quiz session data and distributes them via email. This project is designed with a clean, layered architecture to ensure maintainability, scalability, and high testability.

---

## Features

- Dynamic PDF Generation: Creates professional, multi-page PDF documents summarizing user quiz performance.
- Concurrent Data Fetching: Utilizes goroutines to fetch session and quiz data from the database simultaneously, improving API response times.
- Paginated Email Service: Intelligently breaks down large reports into smaller, paginated chunks and sends them as a series of emails.
- RESTful API: Provides clear and simple endpoints for generating and emailing reports.
- Dependency Injection: Heavily relies on interfaces for dependency injection, making the entire application highly testable.
- Comprehensive Unit Tests: Includes a full suite of unit tests for each layer of the application (repository, service, and handler).

---

## Project Architecture

The project follows the standard Go project layout, separating concerns into distinct layers for clarity and maintainability.

```
.
├── cmd/api/
│   └── main.go               # Application entry point, server startup, and dependency injection.
├── internal/
│   ├── handler/              # Handles HTTP requests and responses.
│   ├── models/               # Defines the core data structures (structs).
│   ├── repository/           # Data access layer for database interactions.
│   ├── router/               # Defines API routes and maps them to handlers.
│   └── service/              # Contains the core business logic.
├── go.mod
├── go.sum
└── README.md
```

---

## API Endpoints

| Endpoint                     | Method | Description                                             |
|------------------------------|--------|---------------------------------------------------------|
| `/sessions/{id}/report`      | GET    | Generates a PDF report for the session and streams it.  |
| `/sessions/{id}/email-report`| POST   | Triggers background process to email the report.        |

---

## Getting Started

### Prerequisites

- Go (version 1.18 or higher)
- A running MySQL instance
- A local SMTP server for testing emails (e.g., [MailHog](https://github.com/mailhog/MailHog))

### Installation & Setup

1. Clone the repository:
    ```sh
    git clone <your-repository-url>
    cd pdf-generator
    ```

2. Configure Database:

    Update the database connection string (DSN) in `internal/repository/db.go` to match your local MySQL setup.

    ```go
    // internal/repository/db.go
    dsn := "root:my-secret-pw@tcp(127.0.0.1:3306)/quizdb?parseTime=true"
    ```

3. Install Dependencies:
    ```sh
    go mod tidy
    ```

4. Run the Application:
    ```sh
    go run cmd/api/main.go
    ```
    The server will start on port `:8070`.

---

## Testing

This project is built with testability as a first-class citizen. You can run the full suite of unit tests and view coverage reports using standard Go tools.

- Run all tests:
    ```sh
    go test ./...
    ```

- Run tests with coverage details:
    ```sh
    go test -cover ./...
    ```

- Generate an interactive HTML coverage report:
    ```sh
    go test -coverprofile=coverage.out ./... && go tool cover -html=coverage.out
    ```

### Testing Approach

The testing strategy is centered around unit testing in isolation. Each layer (repository, service, handler) is tested independently by mocking its dependencies.

#### Repository Testing

The repository layer is tested without a real database using [`go-sqlmock`](https://github.com/DATA-DOG/go-sqlmock). This allows us to verify that the correct SQL is generated and that data is scanned correctly.

**Example: `quizzes_repository_test.go`**
```go
import (
    "regexp"
    "testing"
    "github.com/DATA-DOG/go-sqlmock"
    "github.com/stretchr/testify/assert"
)

func TestQuizzesRepository_GetQuizzesBySessionID(t *testing.T) {
    db, mock, _ := sqlmock.New()
    repo := NewQuizzesRepository(db)

    rows := sqlmock.NewRows([]string{"id", "question"}).AddRow(1, "What is Go?")
    query := `SELECT ... FROM quizzes q JOIN questions qs ...`

    mock.ExpectQuery(regexp.QuoteMeta(query)).
        WithArgs("session-123").
        WillReturnRows(rows)

    quizzes, err := repo.GetQuizzesBySessionID("session-123")

    assert.NoError(t, err)
    assert.Len(t, quizzes, 1)
}
```

#### Service Testing

Services are tested by providing them with mock repositories. This gives us full control over the data and errors the service receives, allowing us to test all business logic paths.

**Example: `pdf_service_test.go`**
```go
type mockSessionRepo struct { /* ... */ }
func (m *mockSessionRepo) GetSessionByID(id string) (*models.Session, error) { /* ... */ }

func TestPDFService_GenerateQuizReport(t *testing.T) {
    sessionRepo := &mockSessionRepo{
        session: &models.Session{SessionID: "sid-1"},
    }
    // ... create other mock repos ...

    service := NewPDFService(sessionRepo, quizzesRepo)
    pdfBytes, err := service.GenerateQuizReport("sid-1")

    assert.NoError(t, err)
    assert.NotEmpty(t, pdfBytes)
}
```

#### Handler Testing

Handlers are tested using `net/http/httptest`. We provide a mock service and send a simulated HTTP request, then capture the response to verify the status code, headers, and body are correct.

**Example: `pdf_handler_test.go`**
```go
type mockPDFService struct { /* ... */ }
func (m *mockPDFService) GenerateQuizReport(id string) ([]byte, error) { /* ... */ }

func TestPDFHandler_GenerateReportHandler_Success(t *testing.T) {
    mockService := &mockPDFService{pdfBytes: []byte("fake-pdf")}
    handler := NewPDFHandler(mockService)

    req, _ := http.NewRequest("GET", "/sessions/sid-123/report", nil)
    rr := httptest.NewRecorder()

    router := mux.NewRouter()
    router.HandleFunc("/sessions/{id}/report", handler.GenerateReportHandler)
    router.ServeHTTP(rr, req)

    assert.Equal(t, http.StatusOK, rr.Code)
    assert.Equal(t, "application/pdf", rr.Header().Get("Content-Type"))
}
```

---

## Core Dependencies

- [github.com/gorilla/mux](https://github.com/gorilla/mux): For powerful and flexible HTTP routing.
- [github.com/go-sql-driver/mysql](https://github.com/go-sql-driver/mysql): For MySQL database connectivity.
- [github.com/jung-kurt/gofpdf](https://github.com/jung-kurt/gofpdf): For PDF document generation.
- [gopkg.in/gomail.v2](https://gopkg.in/gomail.v2): For sending emails.
- [github.com/stretchr/testify](https://github.com/stretchr/testify): For expressive assertions in tests.
- [github.com/DATA-DOG/go-sqlmock](https://github.com/DATA-DOG/go-sqlmock): For mocking the SQL database in repository tests.

---

## License

MIT License. See [LICENSE](LICENSE) for details.

---

## Contributing

Pull requests are welcome! For major changes, please open an issue first to discuss what you would like to change.

---

## Contact

For questions or support, please open an issue on this repository.
