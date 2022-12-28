//go:build integration
// +build integration

package expense

import (
	"fmt"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetExpenseByIDApi(t *testing.T) {

	// setup echo server
	eh := setupServer(t)
	pingServer()

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
	err = res.Decode(&e)

	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, res.StatusCode)
	assert.Equal(t, createExpense.ID, e.ID)
	assert.Equal(t, createExpense.Title, e.Title)
	assert.Equal(t, createExpense.Amount, e.Amount)
	assert.Equal(t, createExpense.Note, e.Note)
	assert.Equal(t, createExpense.Tags, e.Tags)

	// clean up data
	res = request(t, http.MethodDelete, uri(fmt.Sprint(serverPort), fmt.Sprintf("expenses/%d", createExpense.ID)), nil)
	// Assert
	assert.Equal(t, http.StatusOK, res.StatusCode)

	// teardown echo server
	teardownServer(t, eh)
}

func TestGetExpensesAllNotLen0(t *testing.T) {

	// setup echo server
	eh := setupServer(t)
	pingServer()

	// Act
	var expenses []Expense
	res := request(t, http.MethodGet, uri(fmt.Sprint(serverPort), "expenses"), nil)
	err := res.Decode(&expenses)
	// Assert
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, res.StatusCode)
	assert.NotEqual(t, 0, len(expenses))

	teardownServer(t, eh)
}
