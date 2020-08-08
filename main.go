package main

import (
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()
	e.GET("/*", serveWML)
	e.Static("/dl/*", "./static/dl")

	e.Start(":80")
}

func serveWML(c echo.Context) error {
	f, err := os.Open("./static/wap.wml")

	if err == os.ErrNotExist {
		c.String(http.StatusNotFound, "")
	} else if err != nil {
		c.String(http.StatusInternalServerError, "")
	}

	return c.Stream(http.StatusOK, "text/vnd.wap.wml", f)
}
