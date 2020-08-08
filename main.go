package main

import (
	"log"
	"net/http"
	"os"
	"path"
	"strings"

	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()
	e.GET("/*", serveWML)
	e.GET("/dl/*", serveDL)

	e.Start(":80")
}

func serveWML(c echo.Context) error {
	f, err := os.Open("./static/wap.wml")

	if err == os.ErrNotExist {
		c.String(http.StatusNotFound, "")
	} else if err != nil {
		log.Panicln(err)
		return c.String(http.StatusInternalServerError, "")
	}

	return c.Stream(http.StatusOK, "text/vnd.wap.wml", f)
}

func serveDL(c echo.Context) error {
	c.Path()
	req := c.Request()
	if req == nil {
		return c.String(http.StatusInternalServerError, "")
	}
	_, file := path.Split(req.URL.Path)
	f, err := os.Open("./static/dl/" + file)

	if err == os.ErrNotExist {
		return c.String(http.StatusNotFound, "")
	} else if err != nil {
		log.Panicln(err)
		return c.String(http.StatusInternalServerError, "")
	}

	// Go is incorrect when judging Java files
	// this overrides that
	contentType := "binary/octet-stream"
	if strings.HasSuffix(file, ".jar") {
		contentType = "application/java-archive"
	}

	return c.Stream(http.StatusOK, contentType, f)
}
