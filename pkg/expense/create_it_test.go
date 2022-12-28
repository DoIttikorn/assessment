//go:build integration
// +build integration

package expense

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

const serverPort = 2565

func TestCreateExpenseApi(t *testing.T) {
	eh := echo.New()
	go func(e *echo.Echo) {
		db, err := sql.Open("postgres", "postgresql://root:root@db/go-example-db?sslmode=disable")
		if err != nil {
			t.Errorf("error truncating test database tables: %v", err)
		}

		expense := NewApplication(db)

		e.POST("/expenses", expense.CreateExpense)
		e.DELETE("/expenses/:id", expense.DeleteExpenseByID)

		e.Start(fmt.Sprintf(":%d", serverPort))
	}(eh)

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

	body := `{
		"title":"expense",
		"amount": 1000.00,
		"note": "note test",
		"tags": ["dodo", "learn"]
	}`
	var e Expense
	res := request(t, http.MethodPost, uri(fmt.Sprint(serverPort), "expenses"), strings.NewReader(body))
	err := res.Decode(&e)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusCreated, res.StatusCode)
	assert.NotEqual(t, 0, e.ID)
	assert.Equal(t, "expense", e.Title)
	assert.Equal(t, 1000.00, e.Amount)
	assert.Equal(t, "note test", e.Note)
	assert.Equal(t, []string{"dodo", "learn"}, e.Tags)

	res = request(t, http.MethodDelete, uri(fmt.Sprint(serverPort), fmt.Sprintf("expenses/%d", e.ID)), strings.NewReader(body))
	assert.Equal(t, http.StatusOK, res.StatusCode)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err = eh.Shutdown(ctx)
	assert.NoError(t, err)
}

type Response struct {
	*http.Response
	err error
}

// decode ของที่เราอยากจะได้จาก response
func (res *Response) Decode(v interface{}) error {
	if res.err != nil {
		return res.err
	}
	result, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}
	return json.Unmarshal(result, v)
}

// ใช้ในการยิง request ไปยัง server เพื่อทดสอบโดยใช้ http.NewRequest ในการสร้าง request
func request(t *testing.T, method, url string, body io.Reader) *Response {

	if body == nil {
		body = bytes.NewBufferString("")
	}

	req, err := http.NewRequest(method, url, body)
	assert.NoError(t, err)

	// req.Header.Add("Authorization", os.Getenv("AUTH_TOKEN"))
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
