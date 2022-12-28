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

func TestCreateExpenseApi(t *testing.T) {
	// setup echo server
	eh := setupServer(t)
	pingServer()

	// Arrange
	body := `{
		"title":"expense",
		"amount": 1000.00,
		"note": "note test",
		"tags": ["dodo", "learn"]
	}`
	var e Expense

	// Act
	res := request(t, http.MethodPost, uri(fmt.Sprint(serverPort), "expenses"), strings.NewReader(body))
	err := res.Decode(&e)

	// Assert
	assert.Nil(t, err)
	assert.Equal(t, http.StatusCreated, res.StatusCode)
	assert.NotEqual(t, 0, e.ID)
	assert.Equal(t, "expense", e.Title)
	assert.Equal(t, 1000.00, e.Amount)
	assert.Equal(t, "note test", e.Note)
	assert.Equal(t, []string{"dodo", "learn"}, e.Tags)

	// cleanup data
	res = request(t, http.MethodDelete, uri(fmt.Sprint(serverPort), fmt.Sprintf("expenses/%d", e.ID)), nil)
	// Assert
	assert.Equal(t, http.StatusOK, res.StatusCode)

	// teardown echo server
	teardownServer(t, eh)
}
