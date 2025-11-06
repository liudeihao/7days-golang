package main

import (
	"geeweb/gee"
	"net/http"
)

func main() {
	r := gee.Default()
	r.GET("/", func(c *gee.Context) {
		c.String(http.StatusOK, "Hello Geektutu\n")
	})
	// index out of range for testing Recovery()
	r.GET("/hello/:name/space", func(c *gee.Context) {
		c.String(http.StatusOK, "Hello %s, this is your space!\n ", c.Param("name"))
	})
	r.Run(":9999")
}
