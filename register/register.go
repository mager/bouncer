package register

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/bwmarrin/discordgo"
	"github.com/gorilla/mux"

	"github.com/kelseyhightower/envconfig"
	"github.com/mager/bouncer/config"
	"github.com/mager/bouncer/premint"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

func Register(
	lc fx.Lifecycle,
	logger *zap.SugaredLogger,
	router *mux.Router,
	premint premint.PremintClient,
) (
	*zap.SugaredLogger,
	*mux.Router,
	premint.PremintClient,
	config.Config,
	*discordgo.Session,
) {
	var cfg config.Config
	err := envconfig.Process("bouncer", &cfg)
	if err != nil {
		log.Fatal(err.Error())
	}

	var authToken = fmt.Sprintf("Bot %s", cfg.DiscordAuthToken)
	discord, err := discordgo.New(authToken)
	if err != nil {
		log.Fatal(err.Error())
	}

	lc.Append(
		fx.Hook{
			OnStart: func(context.Context) error {
				addr := ":8080"
				logger.Info("Listening on ", addr)

				go http.ListenAndServe(addr, router)
				return nil
			},
			OnStop: func(context.Context) error {
				defer logger.Sync()
				return nil
			},
		},
	)

	return logger, router, premint, cfg, discord
}
