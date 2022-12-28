//go:build integration
// +build integration

package expense

import (
	"bytes"
	"database/sql"
	"fmt"
	"log"
	"net"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestUpdateExpenseByIdSuccess(t *testing.T) {

	eh := echo.New()
	go func(e *echo.Echo) {
		db, err := sql.Open("postgres", "postgresql://root:root@db/go-example-db?sslmode=disable")
		if err != nil {
			log.Fatal(err)
		}

		h := NewApplication(db)

		e.POST("/expenses", h.CreateExpense)
		e.PUT("/expenses/:id", h.UpdateExpenseByID)
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

	// Create for update
	bodyCreate := `{
		"title":"market",
		"amount": 100.00,
		"note": "debt",
		"tags": ["markets"]
	}`

	var createExpense Expense
	res := request(t, http.MethodPost, uri(fmt.Sprint(serverPort), "expenses"), strings.NewReader(bodyCreate))
	assert.Nil(t, res.err)

	err := res.Decode(&createExpense)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusCreated, res.StatusCode)

	// Update
	body := bytes.NewBufferString(`{
		"title":"pay market",
		"amount": 9999.00,
		"note": "clear debt",
		"tags": ["markets", "debt"]
	}`)

	var update Expense
	res = request(t, http.MethodPut, uri(fmt.Sprint(serverPort), fmt.Sprintf("expenses/%d", createExpense.ID)), body)
	err = res.Decode(&update) // ใช้ decode ข้อมูลที่ได้จาก response body มาเก็บไว้ในตัวแปร u

	// Assert
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, res.StatusCode)
	assert.Equal(t, update.ID, createExpense.ID)
	assert.Equal(t, update.Title, "pay market")
	assert.Equal(t, update.Amount, 9999.00)
	assert.Equal(t, update.Note, "clear debt")
	assert.Equal(t, update.Tags, []string{"markets", "debt"})

}
