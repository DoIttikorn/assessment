package main

import (
	"context"
	"crypto/subtle"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	_ "github.com/lib/pq"

	"github.com/Doittikorn/assessment/pkg/config"
	"github.com/Doittikorn/assessment/pkg/db"
	"github.com/Doittikorn/assessment/pkg/expense"
)

func main() {
	c := config.NewConfig()
	// connect to db
	db.ConnectDB()
	// setup expense database
	expense := expense.NewApplication(db.GetDB())

	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	g := e.Group("expenses")
	g.POST("", expense.CreateExpense)
	g.PUT("/:id", expense.UpdateExpenseByID)
	g.GET("", expense.GetExpensesAll, middleware.BasicAuth(basicAuth))
	g.GET("/:id", expense.GetExpenseByID)

	// graceful shutdown
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt, syscall.SIGTERM)

	go func() {
		e.Logger.Fatal(e.Start(":" + c.Port()))
	}()
	fmt.Println("App started.")

	killSignal := <-signals
	switch killSignal {
	case os.Interrupt:
		fmt.Println("Got SIGINT...")
	case syscall.SIGTERM:
		fmt.Println("got SIGTERM...")
	}
	fmt.Println("App is shutting down...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		log.Fatal(err)
	}

}

// basic auth middleware
func basicAuth(username, password string, c echo.Context) (bool, error) {
	// Be careful to use constant time comparison to prevent timing attacks
	if subtle.ConstantTimeCompare([]byte(username), []byte("username")) == 1 &&
		subtle.ConstantTimeCompare([]byte(password), []byte("password")) == 1 {
		return true, nil
	}
	return false, nil
}
