package main

import (
	"net/http"

	"go.uber.org/fx"
	"go.uber.org/zap"
)

func main() {
	fx.New(
		fx.Provide(NewHTTPServer,
			fx.Annotate(
				NewServeMux,
				fx.ParamTags(`name:"echo"`, `name:"hello"`),
			),
			fx.Annotate(
				NewEchoHandler,
				fx.As(new(Route)),
				fx.ResultTags(`name:"echo"`),
			),
			fx.Annotate(NewHelloHandler,
				fx.As(new(Route)),
				fx.ResultTags(`name:"hello"`),
			),
			zap.NewExample),
		fx.Invoke(func(*http.Server) {}),
	).Run()

}
