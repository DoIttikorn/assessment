//go:build integration
// +build integration

package expense

import (
	"bytes"
	"fmt"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUpdateExpenseByIdSuccess(t *testing.T) {

	// setup echo server
	eh := setupServer(t)
	pingServer()

	// setup data for update
	bodyCreate := `{
		"title":"market",
		"amount": 100.00,
		"note": "debt",
		"tags": ["markets"]
	}`

	var createExpense Expense
	res := request(t, http.MethodPost, uri(fmt.Sprint(serverPort), "expenses"), strings.NewReader(bodyCreate))
	assert.Nil(t, res.err)

	// check create data
	err := res.Decode(&createExpense)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusCreated, res.StatusCode)

	// start update case
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

	// cleanup data
	res = request(t, http.MethodDelete, uri(fmt.Sprint(serverPort), fmt.Sprintf("expenses/%d", createExpense.ID)), strings.NewReader(bodyCreate))
	// Assert
	assert.Equal(t, http.StatusOK, res.StatusCode)

	// teardown echo server
	teardownServer(t, eh)

}
