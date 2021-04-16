package core

import (
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
)

func NewServer()  {
	initRouter()
	s := g.Server()
	s.SetIndexFolder(true)
	s.Run()
}

func initRouter()  {
	s := g.Server()

	s.Group("/", func(group *ghttp.RouterGroup) {
		group.GET("/sw.js", ServiceWorkerHandler)
		group.ALL("/:service/*", RequestHandler)
	})
}
