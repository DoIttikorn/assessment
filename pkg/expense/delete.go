package expense

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func (h *handler) DeleteExpenseByID(c echo.Context) error {
	expenseId := c.Param("id")
	if expenseId == "" {
		return c.JSON(http.StatusBadRequest, Error{Message: "Invalid request param id"})
	}

	_, err := h.DB.Exec("DELETE FROM expenses WHERE id = $1", expenseId)
	if err != nil {
		return c.JSON(http.StatusBadRequest, Error{Message: err.Error()})
	}

	return c.JSON(http.StatusOK, nil)
}
