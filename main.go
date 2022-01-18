package main

import (
	"net/http"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"

	"log"
)

func doNothing(c echo.Context) error {
	return nil
}

func main() {

	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.BodyDump(func(c echo.Context, reqBody, resBody []byte) {
		log.Println(string(reqBody))
		log.Println(string(resBody))
	}))
	e.Use(middleware.RemoveTrailingSlash())

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})

	//PLATFORMS
	e.GET("/login/test",
		doNothing,
		RProxy("/login/test",
			[]string{"mspid"}, SimpleForwarder))

	e.Any("/*", doNothing,
		RProxy("/login/test",
			[]string{"mspid"}, SimpleForwarder))

	e.Logger.Fatal(e.Start(":3200"))
}
