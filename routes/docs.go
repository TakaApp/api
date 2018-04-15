package routes

import (
	"net/http"

	"github.com/labstack/echo"
)

// GetDocs returns the docs on / and a 200 Code
// we cant remove it:
// it is important otherwise services (eq cloudfare) will health check on /
// and will think the server died
func GetDocs(c echo.Context) error {
	return c.String(http.StatusOK, "Version 0.5")
}
