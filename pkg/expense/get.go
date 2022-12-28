package expense

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/lib/pq"
)

func (h *handler) GetExpenseByID(echo echo.Context) error {
	var e Expense
	id := echo.Param("id")
	if id == "" {
		return echo.JSON(http.StatusBadRequest, Error{Message: "Invalid request param id"})
	}
	err := h.DB.QueryRow("SELECT id, title, amount, note, tags FROM expenses WHERE id = $1", id).Scan(&e.ID, &e.Title, &e.Amount, &e.Note, pq.Array(&e.Tags))
	if err != nil {
		return echo.JSON(http.StatusInternalServerError, Error{Message: "Error getting expense by id"})
	}
	return echo.JSON(http.StatusOK, e)
}

func (h *handler) GetExpensesAll(echo echo.Context) error {
	var expenses []Expense
	rows, err := h.DB.Query("SELECT id, title, amount, note, tags FROM expenses")
	if err != nil {
		return echo.JSON(http.StatusInternalServerError, Error{Message: "Error getting all expenses"})
	}
	defer rows.Close()
	for rows.Next() {
		var e Expense
		err := rows.Scan(&e.ID, &e.Title, &e.Amount, &e.Note, pq.Array(&e.Tags))
		if err != nil {
			return echo.JSON(http.StatusInternalServerError, Error{Message: "Error getting all expenses"})
		}
		expenses = append(expenses, e)
	}
	return echo.JSON(http.StatusOK, expenses)
}
