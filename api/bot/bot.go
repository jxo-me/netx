package bot

import (
	telebot "github.com/jxo-me/gfbot"
)

var (
	insBotRouter = Routers{
		List: map[string]telebot.IHandler{
			"/Test2":             NewConversation("/Test2", "/cancel"),
			telebot.OnText:       telebot.HandlerFunc(Event.OnText),
			telebot.OnCallback:   telebot.HandlerFunc(Event.OnCallback),
			telebot.OnUserJoined: telebot.HandlerFunc(Event.OnUserJoined),
			// Click Callback
			OnClickAdmissions:  telebot.HandlerFunc(Event.OnClickService),
			OnClickAuthers:     telebot.HandlerFunc(Event.OnClickService),
			OnClickBypass:      telebot.HandlerFunc(Event.OnClickService),
			OnClickHops:        telebot.HandlerFunc(Event.OnClickService),
			OnClickIngress:     telebot.HandlerFunc(Event.OnClickService),
			OnClickServices:    telebot.HandlerFunc(Event.OnClickService),
			OnClickChains:      telebot.HandlerFunc(Event.OnClickService),
			OnClickHosts:       telebot.HandlerFunc(Event.OnClickService),
			OnClickResolver:    telebot.HandlerFunc(Event.OnClickService),
			OnClickLimiter:     telebot.HandlerFunc(Event.OnClickService),
			OnClickConnLimiter: telebot.HandlerFunc(Event.OnClickService),
			OnClickRateLimiter: telebot.HandlerFunc(Event.OnClickService),
			OnClickConfig:      telebot.HandlerFunc(Event.OnClickService),
			// back services
			OnBackServices: telebot.HandlerFunc(Event.OnBackServices),
			// TextCommand
			HostTextCommand: telebot.HandlerFunc(Event.OnHostTextCommand),
		},
		Btns: map[*telebot.Btn]telebot.IHandler{
			//&bot.BtnBetting:       bot.Event.OnBtnBetting,
		},
	}
)

type Routers struct {
	List map[string]telebot.IHandler
	Btns map[*telebot.Btn]telebot.IHandler
	Test map[string]telebot.IHandler
}

func Router() *Routers {
	return &insBotRouter
}
