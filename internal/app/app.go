package app

import (
	"fmt"
	"github.com/gin-gonic/gin"
	repos "goapi/internal/repositories"
	httpserver "goapi/pkg/server"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
)

func oneUserHandler(c *gin.Context) {
	userName := c.Param("name")
	repository := repos.CreateUserRepository()
	user := repository.GetUser(userName)
	c.JSON(http.StatusOK, user)
}

func usersHandler(c *gin.Context) {
	repository := repos.CreateUserRepository()
	size, _ := strconv.ParseInt(c.Query("size"), 10, 64)
	page, _ := strconv.ParseInt(c.Query("page"), 10, 64)
	result := repository.GetUsers(page, size)
	c.JSON(http.StatusOK, result)
}

func Run(port int) {
	repos.FillDb()
	app := gin.New()

	app.GET("/users", usersHandler)
	app.GET("/users/:name", oneUserHandler)

	httpServer := httpserver.New(app)

	// Waiting signal
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	select {
	case s := <-interrupt:
		fmt.Println("app - Run - signal: " + s.String())
	case err := <-httpServer.Notify():
		fmt.Println(fmt.Errorf("app - Run - httpServer.Notify: %w", err))

		err = httpServer.Shutdown()
		if err != nil {
			fmt.Println(fmt.Errorf("app - Run - httpServer.Shutdown: %w", err))
		}
	}
}
