package main

import (
	"fmt"
	"net/http"

	. "github.com/tbxark/g4vercel"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	server := New()
	server.Use(Recovery(func(err interface{}, c *Context) {
		if httpError, ok := err.(HttpError); ok {
			c.JSON(httpError.Status, H{
				"message": httpError.Error(),
			})
		} else {
			message := fmt.Sprintf("%s", err)
			c.JSON(500, H{
				"message": message,
			})
		}
	}))

	// Define route handlers
	server.GET("/main-page", func(c *Context) {
		SendMainPage(c.Writer, c.Req)
	})
	server.GET("/profile", func(c *Context) {
		SendProfilePage(c.Writer, c.Req)
	})
	server.GET("/tag/:tag", func(c *Context) {
		SendTagPage(c.Writer, c.Req)
	})
	server.GET("/info", func(c *Context) {
		SendInfoPage(c.Writer, c.Req)
	})
	server.GET("/img", func(c *Context) {
		SendCatImg(c.Writer, c.Req)
	})

	server.POST("/post/:id", func(c *Context) {
		Post(c.Writer, c.Req)
	})
	server.POST("/login", func(c *Context) {
		Login(c.Writer, c.Req)
	})
	server.POST("/register", func(c *Context) {
		Register(c.Writer, c.Req)
	})
	server.GET("/tags", func(c *Context) {
		GetTags(c.Writer, c.Req)
	})
	server.GET("/user", func(c *Context) {
		UserRouter(c.Writer, c.Req)
	})
	server.POST("/comment", func(c *Context) {
		CommentRouter(c.Writer, c.Req)
	})
	server.GET("/debug/page", func(c *Context) {
		SendDebug(c.Writer, c.Req)
	})
	server.GET("/debug/contents", func(c *Context) {
		Debug(c.Writer, c.Req)
	})

	// Handle the request
	server.Handle(w, r)
}
