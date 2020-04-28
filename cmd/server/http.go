package server

import (
	"ServerStatus/cmd"
	"fmt"
	"github.com/kataras/iris/v12"
)

func LoggerMiddleware(ctx iris.Context) {
	//ctx.Application().Logger().Infof("Runs before %s", ctx.Path())
	ctx.Next()
}

func CorsMiddleware(ctx iris.Context) {
	ctx.Header("Content-Type", "application/json")
	ctx.Header("Access-Control-Allow-Methods", "GET")
	ctx.Header("Access-Control-Allow-Credentials", "true")
	ctx.Header("Access-Control-Allow-origin", "*")
	ctx.Next()
}

func httpRun(p *cmd.Cmd) {
	app := iris.Default()
	app.Use(LoggerMiddleware)
	app.Use(CorsMiddleware)
	app.Use(iris.NoCache)
	app.Use(iris.Gzip)
	app.Get("/", func(ctx iris.Context) {
		ctx.Write(response("update"))
	})
	app.Any("*", func(ctx iris.Context) {
		ctx.WriteString("hello world!")
	})

	app.Listen(fmt.Sprintf("%s:%d", p.Host, p.HTTPPort))
}
