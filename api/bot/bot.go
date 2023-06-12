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
			OnClickServices:    telebot.HandlerFunc(Event.OnClickServices),
			OnClickChains:      telebot.HandlerFunc(Event.OnClickService),
			OnClickHosts:       telebot.HandlerFunc(Event.OnClickService),
			OnClickResolver:    telebot.HandlerFunc(Event.OnClickService),
			OnClickLimiter:     telebot.HandlerFunc(Event.OnClickService),
			OnClickConnLimiter: telebot.HandlerFunc(Event.OnClickService),
			OnClickRateLimiter: telebot.HandlerFunc(Event.OnClickService),
			OnClickConfig:      telebot.HandlerFunc(Event.OnClickService),
			OnClickSaveConfig:  telebot.HandlerFunc(Event.OnClickSaveConfig),

			OnClickNode:    telebot.HandlerFunc(Event.OnClickNodes),
			OnClickAddNode: AddNodeConversation(OnClickAddNode, "/cancel"),
			// back services
			OnBackServices: telebot.HandlerFunc(Event.OnBackServices),
			OnBackHosts:    telebot.HandlerFunc(Event.OnClickNode),
			// TextCommand
			StartTextCommand:   telebot.HandlerFunc(Event.OnStartCommand),
			NodeTextCommand:    telebot.HandlerFunc(Event.OnClickNode),
			ParsingTextCommand: telebot.HandlerFunc(Event.OnParsingCommand),
			GostTextCommand:    telebot.HandlerFunc(Event.OnGostCommand),
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
