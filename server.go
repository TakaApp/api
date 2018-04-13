package main

import (
	"github.com/TakaApp/api/routes"
	"github.com/labstack/echo"
)

func main() {
	e := echo.New()

	e.POST("/trip", routes.GetTrip)

	e.Logger.Fatal(e.Start(":1323"))
}
