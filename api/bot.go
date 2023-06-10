package api

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/glog"
	telebot "github.com/jxo-me/gfbot"
	"github.com/jxo-me/gfbot/middleware"
	"github.com/jxo-me/netx/core/api"
	"log"
)

func NewBot(ctx context.Context, domain, token, pathPrefix string) (*api.TGBot, error) {
	host := fmt.Sprintf("https://%s", domain)
	if pathPrefix != "" {
		host = fmt.Sprintf("%s%s", host, pathPrefix)
	}
	hook := &telebot.HttpHook{
		Endpoint: &telebot.HttpHookEndpoint{
			PublicURL: fmt.Sprintf("%s/v1/bot/hook", host),
		},
	}

	// "5548720536:AAFY-wb4ir22eF5vRMQXft_sj-RDhaB54EQ"
	pref := telebot.Settings{
		Token:        token,
		Poller:       hook,
		Hook:         hook,
		StateStorage: telebot.NewInMemoryStorage(telebot.KeyStrategySenderAndChat),
	}

	b, err := telebot.NewBot(pref)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	logger := glog.New()
	err = logger.SetConfigWithMap(g.Map{
		"path":   "./logs/log",
		"file":   "telegram-bot-{Y-m-d}.log",
		"level":  "all",
		"stdout": true,
	})
	if err != nil {
		return nil, err
	}
	// Global-scoped middleware:
	b.Use(middleware.Logger(ctx, logger))
	b.Use(middleware.AutoRespond())

	return &api.TGBot{
		Domain: domain,
		Token:  token, // dev.us.jxo.me
		Bot:    b,
	}, nil
}
