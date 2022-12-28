package expense

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

type Response struct {
	*http.Response
	err error
}

const serverPort = 2565

func setupServer(t *testing.T) *echo.Echo {
	eh := echo.New()
	go func(e *echo.Echo) {
		db, err := sql.Open("postgres", "postgresql://root:root@db/go-example-db?sslmode=disable")
		if err != nil {
			t.Errorf("error truncating test database tables: %v", err)
		}

		expense := NewApplication(db)

		e.GET("/expenses/:id", expense.GetExpenseByID)
		e.GET("/expenses", expense.GetExpensesAll)
		e.POST("/expenses", expense.CreateExpense)
		e.PUT("/expenses/:id", expense.UpdateExpenseByID)
		e.DELETE("/expenses/:id", expense.DeleteExpenseByID)

		e.Start(fmt.Sprintf(":%d", serverPort))
	}(eh)

	return eh
}

func pingServer() {
	for {
		conn, err := net.DialTimeout("tcp", fmt.Sprintf("localhost:%d", serverPort), 30*time.Second)
		if err != nil {
			log.Println(err)
		}
		if conn != nil {
			conn.Close()
			break
		}
	}
}

func teardownServer(t *testing.T, eh *echo.Echo) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err := eh.Shutdown(ctx)
	assert.NoError(t, err)
}

// Decode decodes the response body into the given interface.
func (res *Response) Decode(v interface{}) error {
	if res.err != nil {
		return res.err
	}
	result, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}
	return json.Unmarshal(result, v)
}

// request creates a new request and returns the response.
func request(t *testing.T, method, url string, body io.Reader) *Response {

	if body == nil {
		body = bytes.NewBufferString("")
	}

	req, err := http.NewRequest(method, url, body)
	assert.NoError(t, err)

	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	client := http.Client{}
	resp, err := client.Do(req)

	assert.NoError(t, err)

	return &Response{resp, err}
}

func uri(port string, paths ...string) string {
	host := "http://localhost:" + port
	if paths == nil {
		return host
	}

	url := append([]string{host}, paths...)
	return strings.Join(url, "/")
}
