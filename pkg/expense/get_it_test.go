//go:build integration
// +build integration

package expense

import (
	"context"
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

func TestGetExpenseByIDApi(t *testing.T) {

	eh := echo.New()
	go func(e *echo.Echo) {
		db, err := sql.Open("postgres", "postgresql://root:root@db/go-example-db?sslmode=disable")
		if err != nil {
			log.Fatal(err)
		}

		h := NewApplication(db)

		e.GET("/expenses/:id", h.GetExpenseByID)
		e.POST("/expenses", h.CreateExpense)
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
		"title":"pay market",
		"amount": 9999.00,
		"note": "clear debt",
		"tags": ["markets", "debt"]
	}`

	var createExpense Expense
	res := request(t, http.MethodPost, uri(fmt.Sprint(serverPort), "expenses"), strings.NewReader(body))
	err := res.Decode(&createExpense) // ใช้ decode ข้อมูลที่ได้จาก response body มาเก็บไว้ในตัวแปร u
	assert.Nil(t, err)
	assert.Equal(t, http.StatusCreated, res.StatusCode)

	var e Expense
	res = request(t, http.MethodGet, uri(fmt.Sprint(serverPort), fmt.Sprintf("expenses/%d", createExpense.ID)), nil)
	err = res.Decode(&e) // ใช้ decode ข้อมูลที่ได้จาก response body มาเก็บไว้ในตัวแปร u

	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, res.StatusCode)
	assert.Equal(t, createExpense.ID, e.ID)
	assert.Equal(t, createExpense.Title, e.Title)
	assert.Equal(t, createExpense.Amount, e.Amount)
	assert.Equal(t, createExpense.Note, e.Note)
	assert.Equal(t, createExpense.Tags, e.Tags)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err = eh.Shutdown(ctx)
	assert.NoError(t, err)
}

func TestGetExpensesAll(t *testing.T) {

	// eh := echo.New()
	// go func(e *echo.Echo) {
	// 	db, err := sql.Open("postgres", "postgresql://root:root@db/go-example-db?sslmode=disable")
	// 	if err != nil {
	// 		log.Fatal(err)
	// 	}

	// 	h := NewApplication(db)

	// 	e.GET("", h.GetExpensesAll)
	// 	e.Start(fmt.Sprintf(":%d", serverPort))
	// }(eh)

	// for {
	// 	conn, err := net.DialTimeout("tcp", fmt.Sprintf("localhost:%d", serverPort), 30*time.Second)
	// 	if err != nil {
	// 		log.Println(err)
	// 	}
	// 	if conn != nil {
	// 		conn.Close()
	// 		break
	// 	}
	// }

	// var expenses []Expense
	// res := request(t, http.MethodGet, uri(fmt.Sprint(serverPort), "expenses"), nil)
	// err := res.Decode(&expenses) // ใช้ decode ข้อมูลที่ได้จาก response body มาเก็บไว้ในตัวแปร u
	// assert.Nil(t, err)
	// assert.Equal(t, http.StatusCreated, res.StatusCode)

	// assert.Nil(t, err)
	// assert.Equal(t, http.StatusOK, res.StatusCode)
	// assert.NotEqual(t, 0, len(expenses))

}
