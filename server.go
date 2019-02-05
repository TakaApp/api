package main

import (
	"github.com/TakaApp/api/routes"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

func main() {
	e := echo.New()

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{echo.GET, echo.PUT, echo.POST, echo.DELETE, echo.OPTIONS},
	}))

	e.GET("/", routes.GetDocs)
	e.POST("/trip", routes.GetTrip)
	e.GET("/search-location/:text", routes.GetSearchLocation)
	e.POST("/reverse-geocoding", routes.ReverseGeocoding)
	e.GET("/stop-details/:routeID/:from/:to", routes.GetStopDetails)

	e.Logger.Fatal(e.Start(":1323"))
}
