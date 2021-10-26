package main

import (
	"github.com/gorilla/mux"

	"github.com/mager/bouncer/handler"
	"github.com/mager/bouncer/logger"
	"github.com/mager/bouncer/premint"
	"github.com/mager/bouncer/register"
	"github.com/mager/bouncer/router"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

func main() {
	fx.New(
		fx.Provide(
			logger.Options,
			router.Options,
			premint.Options,
		),
		fx.Invoke(Register),
	).Run()
}

func Register(
	lc fx.Lifecycle,
	logger *zap.SugaredLogger,
	router *mux.Router,
	premint premint.PremintClient,
) {
	logger, router, premint, cfg, discord := register.Register(lc, logger, router, premint)

	handler.New(logger, router, premint, cfg, discord)
}
