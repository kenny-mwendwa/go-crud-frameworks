package routers

import (
	"github.com/julienschmidt/httprouter"
	"github.com/kenny-mwendwa/go-crud-frameworks/internals/handlers"
)

func HttpRouter() *httprouter.Router {
	r := httprouter.New()

	r.GET("/users", handlers.HttpGetUsers)
	r.POST("/users", handlers.HttpCreateUser)
	r.GET("/users/:id", handlers.HttpGetUser)
	r.PUT("/users/:id", handlers.HttpUpdateUser)
	r.DELETE("/users/:id", handlers.HttpDeleteUser)

	return r
}
