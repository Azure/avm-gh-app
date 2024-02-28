package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"avm-gh-app/pkg"
	"github.com/gregjones/httpcache"
	"github.com/palantir/go-githubapp/githubapp"
	"github.com/rcrowley/go-metrics"
	"github.com/rs/zerolog"
)

func main() {
	_, exist := os.LookupEnv("GH_APP_CONFIG")
	if exist {
		err := pkg.LoadConfigFromEnv()
		if err != nil {
			panic(err)
		}
	} else {
		err := pkg.LoadConfig("githubapp_config.yml")
		if err != nil {
			panic(err)
		}
	}
	logger := zerolog.New(os.Stdout).With().Timestamp().Logger()
	zerolog.DefaultContextLogger = &logger

	metricsRegistry := metrics.DefaultRegistry

	cc, err := githubapp.NewDefaultCachingClientCreator(
		pkg.Cfg.Github,
		githubapp.WithClientUserAgent("avmbot/1.0.0"),
		githubapp.WithClientTimeout(3*time.Second),
		githubapp.WithClientCaching(false, func() httpcache.Cache { return httpcache.NewMemoryCache() }),
		githubapp.WithClientMiddleware(
			githubapp.ClientMetrics(metricsRegistry),
		),
	)
	if err != nil {
		panic(err)
	}

	pushHandler := &pkg.PushHandler{
		ClientCreator: cc,
	}
	dispatcher := githubapp.NewEventDispatcher([]githubapp.EventHandler{pushHandler}, pkg.Cfg.Github.App.WebhookSecret, githubapp.WithScheduler(
		githubapp.AsyncScheduler(),
	))

	//webhookHandler := githubapp.NewDefaultEventDispatcher(config.Github, pushHandler)

	http.Handle(githubapp.DefaultWebhookRoute, dispatcher)

	addr := fmt.Sprintf("%s:%d", pkg.Cfg.Server.Address, pkg.Cfg.Server.Port)
	logger.Info().Msgf("Starting server on %s...", addr)
	err = http.ListenAndServe(addr, nil)
	if err != nil {
		panic(err)
	}
}
