package routes

import (
	"net/http"

	gtfs "github.com/TakaApp/api/gtfs/stops"
	"github.com/labstack/echo"
)

func GetSearchStop(c echo.Context) error {
	text := c.Param("text")

	results := gtfs.SearchStop(text)

	return c.JSON(http.StatusOK, results)
}
