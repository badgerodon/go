package main

import (
	"code.google.com/p/go.net/websocket"
	"github.com/badgerodon/web"
)

func main() {
	web.Route("/", func(ctx web.Context) {
		ctx.Render(nil)
	})
	web.Listen(":80")
}