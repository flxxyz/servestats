package server

import (
	"ServerStatus/cmd"
	"fmt"
	"github.com/kataras/iris/v12"
)

var encodeMap = map[string]bool{
	"json": true,
	"xml":  true,
}

func LoggerMiddleware(ctx iris.Context) {
	//ctx.Application().Logger().Infof("Runs before %s", ctx.Path())
	ctx.Next()
}

func CorsMiddleware(ctx iris.Context) {
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
		switch ctx.URLParamDefault("encode", "json") {
		case "xml":
			ctx.XML(r)
		default:
			ctx.JSON(r)
		}
	})
	app.Any("*", func(ctx iris.Context) {
		ctx.WriteString("hello world!")
	})

	app.Listen(fmt.Sprintf("%s:%d", p.Host, p.HTTPPort))
}
